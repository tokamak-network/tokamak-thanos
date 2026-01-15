// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Script} from 'forge-std/Script.sol';
import {console} from 'forge-std/console.sol';
import {ERC20} from '@openzeppelin/contracts/token/ERC20/ERC20.sol';

// Minimal Interfaces for Uniswap V3
interface IUniswapV3Factory {
  function createPool(
    address tokenA,
    address tokenB,
    uint24 fee
  ) external returns (address pool);
}

interface INonfungiblePositionManager {
  struct MintParams {
    address token0;
    address token1;
    uint24 fee;
    int24 tickLower;
    int24 tickUpper;
    uint256 amount0Desired;
    uint256 amount1Desired;
    uint256 amount0Min;
    uint256 amount1Min;
    address recipient;
    uint256 deadline;
  }
  function mint(
    MintParams calldata params
  )
    external
    payable
    returns (
      uint256 tokenId,
      uint128 liquidity,
      uint256 amount0,
      uint256 amount1
    );
  function factory() external view returns (address);
}

contract MockToken is ERC20 {
  constructor(string memory name, string memory symbol) ERC20(name, symbol) {
    _mint(msg.sender, 1000000 * 10 ** 18);
  }
}

interface IUniswapV3Pool {
  function initialize(uint160 sqrtPriceX96) external;
}

contract CreateDummyUniswapV3 is Script {
  function run() public {
    uint256 deployerPrivateKey = vm.envUint('PRIVATE_KEY');
    vm.startBroadcast(deployerPrivateKey);

    console.log('--- Deploying Mock Uniswap V3 Environment ---');

    // 1. Deploy Mock Tokens
    MockToken token0 = new MockToken('MockA', 'MKA');
    MockToken token1 = new MockToken('MockB', 'MKB');
    console.log('Token0:', address(token0));
    console.log('Token1:', address(token1));

    // 2. Address settings (If NPM is already deployed on devnet/testnet, use that address or find the newly deployed address)
    // Here, for verification purposes, if environment variables are not present, it can be replaced with logic to create a new address,
    // but for actual testing, it is better to receive the NPM address as an argument or read it from environment variables.
    address npmAddr = vm.envOr(
      'NONFUNGIBLE_POSITION_MANAGER',
      0x4200000000000000000000000000000000000504
    );
    INonfungiblePositionManager npm = INonfungiblePositionManager(npmAddr);
    address factoryAddr = npm.factory();

    address pool;
    try
      IUniswapV3Factory(factoryAddr).createPool(
        address(token0),
        address(token1),
        3000
      )
    returns (address p) {
      pool = p;
      console.log('Pool Created at:', pool);
      // Initialize pool with 1:1 price
      IUniswapV3Pool(pool).initialize(79228162514264337593543950336);
      console.log('Pool Initialized.');
    } catch {
      console.log('Pool already exists or failed.');
    }
    // 4. Mint Position
    token0.approve(address(npm), type(uint256).max);
    token1.approve(address(npm), type(uint256).max);

    INonfungiblePositionManager.MintParams
      memory params = INonfungiblePositionManager.MintParams({
        token0: address(token0),
        token1: address(token1),
        fee: 3000,
        tickLower: -60,
        tickUpper: 60,
        amount0Desired: 10 * 10 ** 18,
        amount1Desired: 10 * 10 ** 18,
        amount0Min: 0,
        amount1Min: 0,
        recipient: msg.sender,
        deadline: block.timestamp + 15 minutes
      });

    try npm.mint(params) returns (
      uint256 tokenId,
      uint128 liquidity,
      uint256,
      uint256
    ) {
      console.log('NFT Position Minted! ID:', tokenId);
      console.log('Liquidity:', uint256(liquidity));
    } catch (bytes memory reason) {
      console.log('Mint failed. Reason length:', reason.length);
    }
    vm.stopBroadcast();
  }
}
