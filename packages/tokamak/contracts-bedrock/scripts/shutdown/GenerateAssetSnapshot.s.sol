// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {stdJson} from 'forge-std/StdJson.sol';

import {IERC20} from './interfaces/IERC20.sol';
import {IL1StandardBridge} from './interfaces/IL1StandardBridge.sol';
import {IL2StandardToken} from './interfaces/IL2StandardToken.sol';
import {IOptimismPortal} from './interfaces/IOptimismPortal.sol';
import {IL1UsdcBridge} from './interfaces/IL1UsdcBridge.sol';
import {IL2UsdcBridge} from './interfaces/IL2UsdcBridge.sol';

interface IMulticall3 {
  struct Call3 {
    address target;
    bool allowFailure;
    bytes callData;
  }
  struct Result {
    bool success;
    bytes returnData;
  }
  function aggregate3(
    Call3[] calldata calls
  ) external payable returns (Result[] memory returnData);
  function getEthBalance(address addr) external view returns (uint256 balance);
}

/**
 * @title GenerateAssetSnapshot
 * @notice L2 asset snapshot generator
 * @dev Key principles:
 *      1. Data collection via Python scripts
 *      2. Solidity only handles validation, conversion, and JSON output
 *      3. Avoid complex event scanning
 */
contract GenerateAssetSnapshot is Script {
  using stdJson for string;

  // ========== Configuration ==========
  string l1RpcUrl;
  string l2RpcUrl;
  address l1Bridge;
  address optimismPortal;
  uint256 l2ChainId;
  uint256 finalizedNativeWithdrawals;

  // Custom bridges
  address l1UsdcBridge;
  address l2UsdcBridge;

  address constant NATIVE_TOKEN = address(0);
  address constant L1_NATIVE_TOKEN = 0xa30fe40285B8f5c0457DbC3B7C8A280373c40044;

  // Predeploy addresses - L2 Native tokens
  address constant PREDEPLOY_ETH = 0x4200000000000000000000000000000000000486;
  address constant PREDEPLOY_USDC = 0x4200000000000000000000000000000000000778;

  address constant MULTICALL3 = 0xcA11bde05977b3631167028862bE2a173976CA11;
  uint256 constant MULTICALL_BATCH_SIZE = 200;
  uint256 constant NATIVE_BALANCE_TOLERANCE = 1e16; // 0.01 with 18 decimals

  // ========== Data Structures ==========
  struct TokenInfo {
    address l1Token;
    address l2Token;
    string name;
    uint256 l1Deposit;
    uint256 unclaimedWithdrawals;
    uint256 l2TotalSupply;
    uint256 burnAdjustment;
  }

  struct UserAsset {
    address claimer;
    uint256 amount;
  }

  struct TokenAssets {
    address l1Token;
    address l2Token;
    string tokenName;
    UserAsset[] holders;
  }

  struct UnclaimedWithdrawal {
    address l1Token;
    address l2Token;
    address holder;
    bytes32 withdrawalHash;
    uint256 amount;
  }

  TokenInfo[] public tokens;
  TokenAssets[] public assets;
  UnclaimedWithdrawal[] public unclaimedWithdrawals;

  // ========== Phase 1: Setup ==========

  function setUp() public {
    console.log('=== Asset Snapshot Generator ===\n');

    // 1.1 Load environment variables and set L2 chain ID
    l1RpcUrl = vm.envString('L1_RPC_URL');
    l2RpcUrl = vm.envString('L2_RPC_URL');

    // Require .env.example.thanos-sepolia keys.
    l1Bridge = vm.envAddress('BRIDGE_PROXY');
    optimismPortal = vm.envAddress('OPTIMISM_PORTAL_PROXY');

    // Finalized native withdrawals total on L1 (sum of NativeTokenWithdrawalFinalized).
    // Used to reconcile L2 native balances with L1 deposit accounting.
    finalizedNativeWithdrawals = _computeFinalizedNativeWithdrawals();

    // Custom bridges (optional)
    l1UsdcBridge = vm.envOr('L1_USDC_BRIDGE_PROXY', address(0));
    l2UsdcBridge = vm.envOr('L2_USDC_BRIDGE_PROXY', address(0));

    vm.createSelectFork(l2RpcUrl);
    l2ChainId = block.chainid;

    console.log('L1 RPC:', l1RpcUrl);
    console.log('L2 RPC:', l2RpcUrl);
    console.log('L2 ChainId:', vm.toString(l2ChainId));
    console.log(
      'L1 Finalized Native Withdrawals:',
      vm.toString(finalizedNativeWithdrawals)
    );
    if (l1UsdcBridge != address(0)) {
      console.log('L1 USDC Bridge:', vm.toString(l1UsdcBridge));
    }
    if (l2UsdcBridge != address(0)) {
      console.log('L2 USDC Bridge:', vm.toString(l2UsdcBridge));
    }
    console.log();

    // 1.2 Fetch data via Python (pass chain ID) - skip when SKIP_FETCH=true
    bool skipFetch = vm.envOr('SKIP_FETCH', false);
    if (skipFetch) {
      console.log('[INFO] SKIP_FETCH=true, skipping Python data fetching...\n');
    } else {
      _fetchDataViaPython(l2ChainId);
    }

    // 1.3 Load token list
    _loadTokens();

    // 1.4 Load L2 burn adjustments
    _loadBurnAdjustments();

    // 1.5 Load unclaimed withdrawals
    _loadUnclaimedWithdrawals();
  }

  function _fetchDataViaPython(uint256 chainId) internal {
    console.log('[Fetch] Running Python script...');
    string[] memory fetchInputs = new string[](3);
    fetchInputs[0] = 'python3';
    fetchInputs[1] = 'scripts/shutdown/fetch_explorer_assets.py';
    fetchInputs[2] = vm.toString(chainId);

    require(vm.exists(fetchInputs[1]), 'fetch_explorer_assets.py not found');
    try vm.ffi(fetchInputs) {
      console.log('[OK] Explorer assets fetched');
    } catch {
      console.log('[ERROR] Explorer assets fetch failed (check ffi=true)');
      revert('Explorer assets fetch failed');
    }

    string[] memory burnInputs = new string[](4);
    burnInputs[0] = 'python3';
    burnInputs[1] = 'scripts/shutdown/compute_l2_burns.py';
    burnInputs[2] = l2RpcUrl;
    burnInputs[3] = vm.toString(chainId);

    require(vm.exists(burnInputs[1]), 'compute_l2_burns.py not found');
    try vm.ffi(burnInputs) {
      console.log('[OK] L2 burn adjustments fetched\n');
    } catch {
      console.log('[ERROR] L2 burn adjustment fetch failed (check ffi=true)\n');
      revert('L2 burn adjustment fetch failed');
    }
  }

  function _loadTokens() internal {
    // Try chain-specific file first, then fall back to generic file
    string memory path = string.concat(
      'data/l2-tokens-',
      vm.toString(l2ChainId),
      '.json'
    );
    uint256 skippedZeroCount = 0;
    uint256 skippedMappingCount = 0;
    address firstZeroToken = address(0);
    address firstMappingToken = address(0);

    if (!vm.exists(path)) {
      console.log(
        '[WARN] Chain-specific token list not found, trying generic...'
      );
      path = 'data/l2-tokens.json';
    }

    require(vm.exists(path), 'Token list file not found');

    string memory json = vm.readFile(path);
    address[] memory l2Tokens = json.readAddressArray('$');

    console.log(
      string.concat(
        '[Load] Tokens found in JSON: ',
        vm.toString(l2Tokens.length)
      )
    );

    // Add Native Token
    tokens.push(
      TokenInfo({
        l1Token: L1_NATIVE_TOKEN,
        l2Token: NATIVE_TOKEN,
        name: 'Tokamak Network',
        l1Deposit: 0,
        unclaimedWithdrawals: 0,
        l2TotalSupply: 0,
        burnAdjustment: 0
      })
    );

    // Add other tokens
    for (uint i = 0; i < l2Tokens.length; i++) {
      address l1Token = address(0);

      // Check if address is a contract
      if (l2Tokens[i].code.length == 0) {
        continue;
      }

      // Check if this is a custom bridge token
      if (l2Tokens[i] == PREDEPLOY_USDC) {
        if (l2UsdcBridge == address(0)) {
          if (skippedMappingCount == 0) {
            firstMappingToken = l2Tokens[i];
          }
          skippedMappingCount++;
          continue;
        }
        // USDC uses custom bridge
        try IL2UsdcBridge(l2UsdcBridge).l1Usdc() returns (address _l1Usdc) {
          l1Token = _l1Usdc;
          console.log('[Custom Bridge] USDC L1 Token:', vm.toString(l1Token));
        } catch {
          console.log(
            '[ERROR] Failed to get L1 USDC from custom bridge:',
            vm.toString(l2Tokens[i])
          );
          revert('Failed to get L1 USDC from custom bridge');
        }
      } else {
        // Standard OptimismMintableERC20
        try IL2StandardToken(l2Tokens[i]).l1Token() returns (address _l1) {
          l1Token = _l1;
        } catch {
          if (skippedMappingCount == 0) {
            firstMappingToken = l2Tokens[i];
          }
          skippedMappingCount++;
          continue;
        }
      }

      if (l1Token == address(0)) {
        if (skippedZeroCount == 0) {
          firstZeroToken = l2Tokens[i];
        }
        skippedZeroCount++;
        continue;
      }

      string memory name = _getTokenName(l2Tokens[i]);

      tokens.push(
        TokenInfo({
          l1Token: l1Token,
          l2Token: l2Tokens[i],
          name: name,
          l1Deposit: 0,
          unclaimedWithdrawals: 0,
          l2TotalSupply: 0,
          burnAdjustment: 0
        })
      );
    }

    if (skippedZeroCount > 0) {
      console.log(
        '[INFO] Skipped tokens with zero L1 mapping:',
        vm.toString(skippedZeroCount),
        'example:',
        vm.toString(firstZeroToken)
      );
    }
    if (skippedMappingCount > 0) {
      console.log(
        '[INFO] Skipped tokens without L1 mapping:',
        vm.toString(skippedMappingCount),
        'example:',
        vm.toString(firstMappingToken)
      );
    }
  }

  function _getTokenName(address token) internal view returns (string memory) {
    try IERC20(token).name() returns (string memory name) {
      return name;
    } catch {
      return 'Unknown';
    }
  }

  function _getTokenSymbol(
    address token
  ) internal view returns (string memory) {
    try IERC20(token).symbol() returns (string memory symbol) {
      return symbol;
    } catch {
      return '';
    }
  }

  function _loadUnclaimedWithdrawals() internal {
    string memory path = string.concat(
      'data/unclaimed-withdrawals-',
      vm.toString(l2ChainId),
      '.json'
    );

    require(vm.exists(path), 'Unclaimed withdrawals file not found');

    string memory json = vm.readFile(path);
    if (_isEmptyJsonArray(json)) {
      console.log('[INFO] No unclaimed withdrawals\n');
      return;
    }

    // Parse JSON array manually (simplified approach)
    console.log('[Load] Unclaimed withdrawals file found');

    uint256 count = _countUnclaimedEntries(json);
    if (count == 0) {
      console.log('[INFO] No unclaimed withdrawals\n');
      return;
    }

    console.log('[INFO] Processing unclaimed withdrawals:', count);

    for (uint i = 0; i < count; i++) {
      string memory base = string.concat('$[', vm.toString(i), '].');
      address l1Token = json.readAddress(string.concat(base, 'l1Token'));
      address l2Token = json.readAddress(string.concat(base, 'l2Token'));
      address holder = json.readAddress(string.concat(base, 'holder'));
      bytes32 withdrawalHash = json.readBytes32(
        string.concat(base, 'withdrawalHash')
      );
      uint256 amount = _parseUint(
        json.readString(string.concat(base, 'amount'))
      );
      unclaimedWithdrawals.push(
        UnclaimedWithdrawal({
          l1Token: l1Token,
          l2Token: l2Token,
          holder: holder,
          withdrawalHash: withdrawalHash,
          amount: amount
        })
      );

      // Update token info for Phase 2 validation
      bool found = false;
      for (uint j = 0; j < tokens.length; j++) {
        if (tokens[j].l1Token == l1Token && tokens[j].l2Token == l2Token) {
          tokens[j].unclaimedWithdrawals += amount;
          found = true;
          break;
        }
      }

      if (!found) {
        // Skip unclaimed entries for tokens not in the active token list.
        continue;
      }
    }
    console.log('[INFO] Unclaimed withdrawals aggregated\n');
  }

  function _loadBurnAdjustments() internal {
    string memory path = string.concat(
      'data/l2-burns-',
      vm.toString(l2ChainId),
      '.json'
    );

    require(vm.exists(path), 'L2 burn adjustment file not found');

    string memory json = vm.readFile(path);
    if (_isEmptyJsonArray(json)) {
      console.log('[INFO] No L2 burn adjustments\n');
      return;
    }

    uint256 count = _countEntries(json, bytes('"l2Token"'));
    if (count == 0) {
      console.log('[INFO] No L2 burn adjustments\n');
      return;
    }

    console.log('[INFO] Processing L2 burn adjustments:', count);

    uint256 missingTokenCount = 0;
    address firstMissingToken = address(0);
    for (uint i = 0; i < count; i++) {
      string memory base = string.concat('$[', vm.toString(i), '].');
      address l2Token = json.readAddress(string.concat(base, 'l2Token'));
      uint256 extraBurn = _parseUint(
        json.readString(string.concat(base, 'extraBurn'))
      );

      bool found = false;
      for (uint j = 0; j < tokens.length; j++) {
        if (tokens[j].l2Token == l2Token) {
          tokens[j].burnAdjustment = extraBurn;
          found = true;
          break;
        }
      }

      if (!found) {
        if (missingTokenCount == 0) {
          firstMissingToken = l2Token;
        }
        missingTokenCount++;
      }
    }
    if (missingTokenCount > 0) {
      console.log(
        '[INFO] Burn adjustment tokens not in list:',
        vm.toString(missingTokenCount),
        'example:',
        vm.toString(firstMissingToken)
      );
    }
    console.log('[INFO] L2 burn adjustments aggregated\n');
  }

  function _parseUint(string memory value) internal pure returns (uint256) {
    bytes memory data = bytes(value);
    require(data.length > 0, 'Invalid uint string');
    uint256 start = 0;
    uint256 end = data.length;

    // Trim leading/trailing whitespace from FFI output.
    while (start < end && _isWhitespace(data[start])) {
      start++;
    }
    while (end > start && _isWhitespace(data[end - 1])) {
      end--;
    }
    require(end > start, 'Invalid uint string');
    uint256 result = 0;

    // Accept hex-encoded output (e.g., "0x...") from FFI.
    if (
      end - start >= 2 &&
      data[start] == '0' &&
      (data[start + 1] == 'x' || data[start + 1] == 'X')
    ) {
      start += 2;
      for (uint i = start; i < end; i++) {
        uint8 c = uint8(data[i]);
        if (c >= 48 && c <= 57) {
          result = (result << 4) + (c - 48);
        } else if (c >= 65 && c <= 70) {
          result = (result << 4) + (c - 55);
        } else if (c >= 97 && c <= 102) {
          result = (result << 4) + (c - 87);
        } else {
          revert('Invalid uint string');
        }
      }
      return result;
    }

    // Extract first contiguous decimal number from output (e.g., "DEC:123").
    uint256 i = start;
    while (i < end && (uint8(data[i]) < 48 || uint8(data[i]) > 57)) {
      i++;
    }
    require(i < end, 'Invalid uint string');
    while (i < end && uint8(data[i]) >= 48 && uint8(data[i]) <= 57) {
      result = result * 10 + (uint8(data[i]) - 48);
      i++;
    }
    return result;
  }

  function _computeFinalizedNativeWithdrawals() internal returns (uint256) {
    string memory scriptPath = 'scripts/shutdown/compute_finalized_native_withdrawals.py';
    require(vm.exists(scriptPath), 'compute_finalized_native_withdrawals.py not found');
    string[] memory inputs = new string[](4);
    inputs[0] = 'python3';
    inputs[1] = scriptPath;
    inputs[2] = l1RpcUrl;
    inputs[3] = vm.toString(l1Bridge);

    try vm.ffi(inputs) returns (bytes memory stdout) {
      string memory output = string(stdout);
      return _parseUint(output);
    } catch {
      console.log('[ERROR] Failed to compute finalized native withdrawals');
      revert('Failed to compute finalized native withdrawals');
    }
  }

  function _countUnclaimedEntries(
    string memory json
  ) internal pure returns (uint256) {
    return _countEntries(json, bytes('"withdrawalHash"'));
  }

  function _countEntries(
    string memory json,
    bytes memory needle
  ) internal pure returns (uint256) {
    if (needle.length == 0) {
      return 0;
    }

    bytes memory data = bytes(json);
    if (data.length < needle.length) {
      return 0;
    }

    uint256 count = 0;
    for (uint256 i = 0; i + needle.length <= data.length; i++) {
      bool matchFound = true;
      for (uint256 j = 0; j < needle.length; j++) {
        if (data[i + j] != needle[j]) {
          matchFound = false;
          break;
        }
      }
      if (matchFound) {
        count++;
        i += needle.length - 1;
      }
    }
    return count;
  }

  function _isEmptyJsonArray(string memory json) internal pure returns (bool) {
    bytes memory data = bytes(json);
    uint256 i = 0;

    while (i < data.length && _isWhitespace(data[i])) {
      i++;
    }
    if (i >= data.length || data[i] != '[') {
      return false;
    }
    i++;
    while (i < data.length && _isWhitespace(data[i])) {
      i++;
    }
    if (i >= data.length || data[i] != ']') {
      return false;
    }
    i++;
    while (i < data.length) {
      if (!_isWhitespace(data[i])) {
        return false;
      }
      i++;
    }
    return true;
  }

  function _isWhitespace(bytes1 c) internal pure returns (bool) {
    return c == 0x20 || c == 0x0a || c == 0x0d || c == 0x09;
  }

  function _abs(uint256 a, uint256 b) internal pure returns (uint256) {
    return a > b ? a - b : b - a;
  }

  function _isL2Native(
    address l2Token,
    address l1Token,
    string memory tokenName,
    string memory tokenSymbol
  ) internal pure returns (bool) {
    // Check token name/symbol for "Bridged" or ".e" suffix
    // These indicate L1-bridged tokens even if l1Token is unknown
    bytes memory nameBytes = bytes(tokenName);
    bytes memory symbolBytes = bytes(tokenSymbol);

    // "Bridged" in name → NOT L2 native
    for (uint i = 0; i + 7 <= nameBytes.length; i++) {
      if (
        nameBytes[i] == 'B' &&
        nameBytes[i + 1] == 'r' &&
        nameBytes[i + 2] == 'i' &&
        nameBytes[i + 3] == 'd' &&
        nameBytes[i + 4] == 'g' &&
        nameBytes[i + 5] == 'e' &&
        nameBytes[i + 6] == 'd'
      ) {
        return false; // Bridged token
      }
    }

    // ".e" suffix in symbol → NOT L2 native (e.g., USDC.e, ETH.e)
    if (
      symbolBytes.length >= 2 &&
      symbolBytes[symbolBytes.length - 2] == '.' &&
      symbolBytes[symbolBytes.length - 1] == 'e'
    ) {
      return false; // Ethereum-bridged token
    }

    // L2 tokens with no L1 counterpart (l1Token = 0x000 and not ETH)
    if (
      l1Token == address(0) &&
      l2Token != PREDEPLOY_ETH &&
      l2Token != NATIVE_TOKEN
    ) {
      return true;
    }

    return false;
  }

  function _diagnoseMismatch(
    string memory tokenName,
    uint256 l1Locked,
    uint256 l2Supply,
    address l1Token,
    address l2Token,
    string memory tokenSymbol
  ) internal pure returns (string memory) {
    if (l1Locked == l2Supply) return '';

    if (l1Locked > l2Supply) {
      // L1 is larger - common case
      return
        '  Possible causes: In-flight withdrawals, Pending deposits, Failed L2 mints';
    } else {
      // L2 is larger
      if (_isL2Native(l2Token, l1Token, tokenName, tokenSymbol)) {
        return '  [INFO] L2 Native token - L1 deposit = 0 is expected';
      } else {
        // Bridged token with missing L1 deposits
        return
          '  [WARN] Bridged token with missing L1 deposits - Check custom bridge or L1 token mapping';
      }
    }
  }

  function _serializeHolders(
    UserAsset[] storage holders
  ) internal view returns (string memory) {
    if (holders.length == 0) return '[]';
    string memory json = '[\n';
    for (uint i = 0; i < holders.length; i++) {
      if (i > 0) json = string.concat(json, ',\n');
      json = string.concat(
        json,
        '      {"claimer": "',
        vm.toString(holders[i].claimer),
        '", "amount": "',
        vm.toString(holders[i].amount),
        '"}'
      );
    }
    return string.concat(json, '\n    ]');
  }

  function _escapeJson(
    string memory str
  ) internal pure returns (string memory) {
    bytes memory strBytes = bytes(str);
    uint escapeCount = 0;

    // Count characters that need escaping
    for (uint i = 0; i < strBytes.length; i++) {
      if (strBytes[i] == '"' || strBytes[i] == '\\') {
        escapeCount++;
      }
    }

    if (escapeCount == 0) return str;

    // Create new string with escaped characters
    bytes memory result = new bytes(strBytes.length + escapeCount);
    uint j = 0;
    for (uint i = 0; i < strBytes.length; i++) {
      if (strBytes[i] == '"' || strBytes[i] == '\\') {
        result[j++] = '\\';
      }
      result[j++] = strBytes[i];
    }

    return string(result);
  }

  function _makeHash(
    address l1Token,
    address claimer,
    uint256 amount
  ) internal pure returns (bytes32) {
    return keccak256(abi.encodePacked(l1Token, claimer, amount));
  }

  function _serializeWithHash(
    address l1Token,
    UserAsset[] storage holders
  ) internal view returns (string memory) {
    if (holders.length == 0) return '[]';
    string memory json = '[\n';
    for (uint i = 0; i < holders.length; i++) {
      if (i > 0) json = string.concat(json, ',\n');
      bytes32 hash = _makeHash(l1Token, holders[i].claimer, holders[i].amount);
      json = string.concat(json, _serializeOneWithHash(holders[i], hash));
    }
    return string.concat(json, '\n    ]');
  }

  function _serializeOneWithHash(
    UserAsset storage holder,
    bytes32 hash
  ) internal view returns (string memory) {
    return
      string.concat(
        '      {\n',
        '        "claimer": "',
        vm.toString(holder.claimer),
        '",\n',
        '        "amount": "',
        vm.toString(holder.amount),
        '",\n',
        '        "hash": "',
        vm.toString(hash),
        '"\n',
        '      }'
      );
  }

  // ========== Phase 2: Validate Balances ==========

  function phase2_validateBalances() public {
    console.log('=== Phase 2: Validate Balances ===');
    console.log('Total tokens:', vm.toString(tokens.length));
    console.log('-------------------------------------------\n');

    for (uint i = 0; i < tokens.length; i++) {
      TokenInfo storage t = tokens[i];

      // Skip L2 native tokens (l1Token = 0x00 and not ETH/Native) - no log output
      if (
        t.l1Token == address(0) &&
        t.l2Token != PREDEPLOY_ETH &&
        t.l2Token != NATIVE_TOKEN
      ) {
        continue;
      }

      // First, check L2 supply to detect broken tokens early
      vm.createSelectFork(l2RpcUrl);
      if (t.l2Token == NATIVE_TOKEN) {
        t.l2TotalSupply = vm.envOr('L2_NATIVE_TOTAL_SUPPLY', uint256(0));
      } else {
        try IERC20(t.l2Token).totalSupply() returns (uint256 supply) {
          t.l2TotalSupply = supply;
        } catch {
          console.log('  [WARN] Failed to get L2 total supply');
        }
      }

      // Skip tokens with no L2 supply AND no unclaimed withdrawals (broken/abandoned tokens)
      if (
        t.l2TotalSupply == 0 &&
        t.unclaimedWithdrawals == 0 &&
        t.l2Token != NATIVE_TOKEN
      ) {
        continue;
      }

      console.log(
        string.concat(
          '[',
          vm.toString(i + 1),
          '/',
          vm.toString(tokens.length),
          '] ',
          t.name
        )
      );
      console.log('  L2 Token:', vm.toString(t.l2Token));
      console.log('  L1 Token:', vm.toString(t.l1Token));

      // Get L1 deposit
      vm.createSelectFork(l1RpcUrl);
      if (t.l2Token == NATIVE_TOKEN) {
        if (t.l1Token == address(0)) {
          t.l1Deposit = optimismPortal.balance;
        } else {
          t.l1Deposit = IERC20(t.l1Token).balanceOf(optimismPortal);
        }
      } else if (t.l2Token == PREDEPLOY_ETH) {
        // ETH (predeploy) is bridged via L1StandardBridge deposits mapping.
        try IL1StandardBridge(l1Bridge).deposits(address(0), PREDEPLOY_ETH) returns (
          uint256 deposit
        ) {
          t.l1Deposit = deposit;
        } catch {
          console.log('  [WARN] Failed to get L1 ETH deposit from bridge');
        }
      } else if (t.l2Token == PREDEPLOY_USDC && l1UsdcBridge != address(0)) {
        // USDC uses custom bridge
        if (t.l1Token == address(0)) {
          console.log('  [WARN] L1 USDC token is not set');
        } else {
          try
            IL1UsdcBridge(l1UsdcBridge).deposits(t.l1Token, t.l2Token)
          returns (uint256 deposit) {
            t.l1Deposit = deposit;
            console.log(
              '  [Custom Bridge] USDC Deposits:',
              vm.toString(deposit)
            );
          } catch {
            console.log(
              '  [WARN] Failed to get L1 USDC deposit from custom bridge'
            );
          }
        }
      } else {
        // Standard bridge
        try IL1StandardBridge(l1Bridge).deposits(t.l1Token, t.l2Token) returns (
          uint256 deposit
        ) {
          t.l1Deposit = deposit;
        } catch {
          console.log('  [WARN] Failed to get L1 deposit');
        }
      }

      console.log('  L1 Deposit:            ', vm.toString(t.l1Deposit));
      console.log(
        '  Unclaimed Withdrawals: ',
        vm.toString(t.unclaimedWithdrawals)
      );
      console.log('  Burn Adjustment:       ', vm.toString(t.burnAdjustment));

      // Calculation logic:
      // - l1Deposit (Bridge.deposits[]): Increases on deposit, decreases only when withdrawal is finalized on L1.
      // - l2TotalSupply: Increases on deposit mint, decreases immediately on L2 burn (withdrawal start).
      // - unclaimedWithdrawals: Assets burned on L2 but not yet finalized on L1.
      // - burnAdjustment: Extra L2 burns not represented in withdrawal logs.
      // Correct equation: l1Deposit == l2TotalSupply + unclaimedWithdrawals + burnAdjustment

      uint256 expectedL1 = t.l2TotalSupply +
        t.unclaimedWithdrawals +
        t.burnAdjustment;

      if (t.l2Token == NATIVE_TOKEN) {
        console.log('  L2 Supply:              (Calculated in Phase 3)');
      } else {
        console.log('  L2 Supply:             ', vm.toString(t.l2TotalSupply));

        if (t.l1Deposit == expectedL1) {
          console.log('  Status: [OK] Match');
        } else {
          console.log('  Status: [WARN] Mismatch');
          uint256 diff = _abs(t.l1Deposit, expectedL1);
          console.log('  Difference:            ', vm.toString(diff));

          // Diagnose the mismatch on L2 fork
          vm.createSelectFork(l2RpcUrl);
          string memory tokenSymbol = _getTokenSymbol(t.l2Token);

          string memory diagnosis = _diagnoseMismatch(
            t.name,
            t.l1Deposit, // Pass actual deposit
            expectedL1, // Pass what we expected based on L2 state
            t.l1Token,
            t.l2Token,
            tokenSymbol
          );
          if (bytes(diagnosis).length > 0) {
            console.log(diagnosis);
          }
          // Switch back to L1 if needed (though the loop ends or restarts with a fork select)
          vm.createSelectFork(l1RpcUrl);
        }
      }
      console.log('');
    }

    console.log('-------------------------------------------');
    console.log('Phase 2 Complete\n');
  }

  // ========== Phase 3: Collect Assets ==========

  function phase3_collectAssets() public {
    console.log('=== Phase 3: Collect Assets ===');
    vm.createSelectFork(l2RpcUrl);

    // Load holders from JSON
    address[] memory holders = _loadHolders();
    console.log('Total holders:', vm.toString(holders.length));
    console.log('Total tokens:', vm.toString(tokens.length));
    console.log('-------------------------------------------\n');

    for (uint i = 0; i < tokens.length; i++) {
      TokenInfo storage t = tokens[i];

      // Skip L2 native tokens (l1Token = 0x00 and not ETH/Native) - no log output
      if (
        t.l1Token == address(0) &&
        t.l2Token != PREDEPLOY_ETH &&
        t.l2Token != NATIVE_TOKEN
      ) {
        continue;
      }

      // Skip tokens with no L2 supply AND no unclaimed withdrawals
      if (
        t.l2TotalSupply == 0 &&
        t.unclaimedWithdrawals == 0 &&
        t.l2Token != NATIVE_TOKEN
      ) {
        continue;
      }

      console.log(
        string.concat(
          '[',
          vm.toString(i + 1),
          '/',
          vm.toString(tokens.length),
          '] ',
          t.name
        )
      );
      console.log('  L2 Token:', vm.toString(t.l2Token));

      TokenAssets storage ta = assets.push();
      ta.l1Token = t.l1Token;
      ta.l2Token = t.l2Token;
      ta.tokenName = t.name;

      // Query balances in batches using Multicall3
      uint256 total = 0;
      uint256 len = holders.length;
      for (uint j = 0; j < len; ) {
        uint256 currentBatch = len - j;
        if (currentBatch > MULTICALL_BATCH_SIZE) {
          currentBatch = MULTICALL_BATCH_SIZE;
        }

        IMulticall3.Call3[] memory calls = new IMulticall3.Call3[](
          currentBatch
        );
        if (t.l2Token == NATIVE_TOKEN) {
          // Native token (ETH) - individual calls to avoid multicall self-call issues in fork
          for (uint k = 0; k < currentBatch; k++) {
            uint256 bal = holders[j + k].balance;
            if (bal > 0) {
              ta.holders.push(UserAsset(holders[j + k], bal));
              unchecked {
                total += bal;
              }
            }
          }
        } else {
          // ERC20 tokens - Batch using Multicall3
          for (uint k = 0; k < currentBatch; k++) {
            address holder = holders[j + k];
            calls[k] = IMulticall3.Call3({
              target: t.l2Token,
              allowFailure: true,
              callData: abi.encodeWithSelector(
                IERC20.balanceOf.selector,
                holder
              )
            });
          }

          // Execute batch
          IMulticall3.Result[] memory results = IMulticall3(MULTICALL3)
            .aggregate3(calls);

          // Process results
          for (uint k = 0; k < results.length; k++) {
            if (results[k].success && results[k].returnData.length >= 32) {
              uint256 bal = abi.decode(results[k].returnData, (uint256));
              if (bal > 0) {
                ta.holders.push(UserAsset(holders[j + k], bal));
                unchecked {
                  total += bal;
                }
              }
            }
          }
        }

        j += currentBatch;
      }

      console.log('  Holders with balance:', vm.toString(ta.holders.length));
      console.log('  Collected total:     ', vm.toString(total));

      uint256 expected = t.l2TotalSupply;
      string memory label = 'Expected (L2 Supply):';

      if (t.l2Token == NATIVE_TOKEN) {
        expected = t.l1Deposit + finalizedNativeWithdrawals;
        label = 'Expected (L1 Deposit + Finalized Withdrawals):';
      }

      console.log('  ', label, vm.toString(expected));

      if (
        total == expected ||
        (t.l2Token == NATIVE_TOKEN && _abs(total, expected) <= NATIVE_BALANCE_TOLERANCE)
      ) {
        console.log('  Status: [OK] Match');
      } else {
        console.log('  Status: [WARN] Mismatch');
        console.log(
          '  Difference:          ',
          vm.toString(_abs(total, expected))
        );
      }
      console.log('');
    }

    console.log('-------------------------------------------');
    console.log('Phase 3 Complete\n');
  }

  function _loadHolders() internal returns (address[] memory) {
    address[] memory holders = _loadAddressList(
      'data/l2-holders-',
      'data/l2-holders.json',
      'Holder list file not found'
    );

    address[] memory contracts = _loadOptionalAddressList(
      'data/l2-contracts-',
      'data/l2-contracts.json'
    );

    if (contracts.length > 0) {
      console.log('  [INFO] Merging contract addresses into holders');
      holders = _mergeUnique(holders, contracts);
    }

    if (unclaimedWithdrawals.length > 0) {
      console.log('  [INFO] Merging unclaimed withdrawal holders');
      address[] memory unclaimedHolders = new address[](
        unclaimedWithdrawals.length
      );
      for (uint i = 0; i < unclaimedWithdrawals.length; i++) {
        unclaimedHolders[i] = unclaimedWithdrawals[i].holder;
      }
      holders = _mergeUnique(holders, unclaimedHolders);
    }

    console.log('  [INFO] Total holders after merge:', holders.length);

    // Filter out systemic addresses (Forge default sender, msg.sender, address(0))
    address forgeDefaultSender = 0x1804c8AB1F12E6bbf3894d4083f33e07309d1f38;
    address[] memory filtered = new address[](holders.length);
    uint256 count = 0;

    for (uint i = 0; i < holders.length; i++) {
      address h = holders[i];
      if (h == address(0) || h == forgeDefaultSender || h == msg.sender)
        continue;
      filtered[count++] = h;
    }

    address[] memory result = new address[](count);
    for (uint i = 0; i < count; i++) {
      result[i] = filtered[i];
    }

    console.log('  [INFO] Total holders after filtering:', result.length);
    return result;
  }

  function _loadAddressList(
    string memory chainPrefix,
    string memory fallbackPath,
    string memory missingError
  ) internal returns (address[] memory) {
    string memory path = string.concat(
      chainPrefix,
      vm.toString(l2ChainId),
      '.json'
    );

    if (!vm.exists(path)) {
      console.log('[WARN] Chain-specific file not found, trying generic...');
      path = fallbackPath;
    }

    require(vm.exists(path), missingError);

    string memory json = vm.readFile(path);
    return json.readAddressArray('$');
  }

  function _loadOptionalAddressList(
    string memory chainPrefix,
    string memory fallbackPath
  ) internal returns (address[] memory) {
    string memory path = string.concat(
      chainPrefix,
      vm.toString(l2ChainId),
      '.json'
    );

    if (!vm.exists(path)) {
      path = fallbackPath;
    }

    if (!vm.exists(path)) {
      return new address[](0);
    }

    string memory json = vm.readFile(path);
    return json.readAddressArray('$');
  }

  function _mergeUnique(
    address[] memory a,
    address[] memory b
  ) internal pure returns (address[] memory) {
    if (b.length == 0) return a;

    address[] memory temp = new address[](a.length + b.length);
    uint256 count = 0;

    for (uint i = 0; i < a.length; i++) {
      temp[count++] = a[i];
    }

    for (uint j = 0; j < b.length; j++) {
      if (!_contains(temp, count, b[j])) {
        temp[count++] = b[j];
      }
    }

    address[] memory merged = new address[](count);
    for (uint k = 0; k < count; k++) {
      merged[k] = temp[k];
    }
    return merged;
  }

  function _contains(
    address[] memory arr,
    uint256 length,
    address target
  ) internal pure returns (bool) {
    for (uint i = 0; i < length; i++) {
      if (arr[i] == target) return true;
    }
    return false;
  }

  function _getBalance(
    address token,
    address user
  ) internal view returns (uint256) {
    if (token == NATIVE_TOKEN) {
      return user.balance;
    }
    try IERC20(token).balanceOf(user) returns (uint256 bal) {
      return bal;
    } catch {
      return 0;
    }
  }

  function _saveAssets(string memory filename) internal {
    string memory json = '[\n';

    for (uint i = 0; i < assets.length; i++) {
      if (i > 0) json = string.concat(json, ',\n');

      TokenAssets storage ta = assets[i];
      json = string.concat(
        json,
        '  {\n',
        '    "l1Token": "',
        vm.toString(ta.l1Token),
        '",\n',
        '    "l2Token": "',
        vm.toString(ta.l2Token),
        '",\n',
        '    "tokenName": "',
        _escapeJson(ta.tokenName),
        '",\n',
        '    "holders": ',
        _serializeHolders(ta.holders),
        '\n  }'
      );
    }

    json = string.concat(json, '\n]');
    vm.writeFile(filename, json);
    console.log('Saved:', filename, '\n');
  }

  // ========== Phase 4: Add Hash ==========

  function phase4_addHash() public {
    console.log('=== Phase 4: Add Hash ===');
    console.log('Generating final output with hash...');
    console.log('-------------------------------------------\n');

    string memory json = '[\n';
    uint256 totalHolders = 0;
    bool firstEntry = true;

    for (uint i = 0; i < assets.length; i++) {
      TokenAssets storage ta = assets[i];

      // Filter out L2 native tokens (l1Token = 0x00 and not ETH/Native) - no log output
      if (
        ta.l1Token == address(0) &&
        ta.l2Token != PREDEPLOY_ETH &&
        ta.l2Token != NATIVE_TOKEN
      ) {
        continue;
      }

      // Filter out tokens with no holders - no log output
      if (ta.holders.length == 0) {
        continue;
      }

      if (!firstEntry) json = string.concat(json, ',\n');
      firstEntry = false;

      console.log(
        string.concat(
          '[',
          vm.toString(i + 1),
          '/',
          vm.toString(assets.length),
          '] ',
          ta.tokenName
        )
      );
      console.log('  Holders with hash:', vm.toString(ta.holders.length));

      totalHolders += ta.holders.length;

      json = string.concat(
        json,
        '  {\n',
        '    "l1Token": "',
        vm.toString(ta.l1Token),
        '",\n',
        '    "l2Token": "',
        vm.toString(ta.l2Token),
        '",\n',
        '    "tokenName": "',
        _escapeJson(ta.tokenName),
        '",\n',
        '    "data": ',
        _serializeWithHash(ta.l1Token, ta.holders),
        '\n  }'
      );
    }

    json = string.concat(json, '\n]');

    string memory filename = string.concat(
      'data/generate-assets-',
      vm.toString(l2ChainId),
      '.json'
    );
    vm.writeFile(filename, json);

    console.log('');
    console.log('-------------------------------------------');
    console.log('Final output saved:', filename);
    console.log('Total tokens:', vm.toString(assets.length));
    console.log('Total holders:', vm.toString(totalHolders));

    console.log('Phase 4 Complete\n');
  }

  // ========== Run ==========

  function run() public {
    console.log('');
    console.log('===========================================');
    console.log('  Asset Snapshot Generation Started');
    console.log('===========================================');
    console.log('');

    phase2_validateBalances();
    phase3_collectAssets();
    phase4_addHash();

    console.log('===========================================');
    console.log('  All Phases Complete');
    console.log('===========================================');
    console.log('');
  }
}
