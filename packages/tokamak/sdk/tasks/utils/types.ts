export const L1Interface = [
  'function deposits(address _firstKey, address _secondKey) public view returns (uint256)',
  'event ETHWithdrawalFinalized(address indexed _from, address indexed _to,uint256 _amount,bytes _data)',
  'event ERC20WithdrawalFinalized(address indexed _l1Token, address indexed _l2Token, address indexed _from,address _to,uint256 _amount,bytes _data)',
]

export const L2Interface = [
  'event DepositFinalized(address indexed _l1Token, address indexed _l2Token, address indexed _from, address _to, uint256 _amount, bytes _data)',
  'event Mint(address indexed _account, uint256 _amount)',
  'event Burn(address indexed _account, uint256 _amount)',
  'event WithdrawalInitiated( address indexed _l1Token, address indexed _l2Token, address indexed _from, address _to, uint256 _amount, bytes _data)',
]

export const ERC20 = [
  'function balanceOf(address account) external view returns (uint256)',
  'function totalSupply() external view returns (uint256)',
  'event Transfer(address indexed _from, address indexed _to, uint256 _value)',
  'function name() public view returns (string)',
  'function symbol() public view returns (string)',
  'function decimals() public view returns (uint8)',
]

// v3
export const NonfungibleTokenPositionManager = [
  'event IncreaseLiquidity(uint256 indexed tokenId, uint128 liquidity, uint256 amount0, uint256 amount1)',
  'event DecreaseLiquidity(uint256 indexed tokenId, uint128 liquidity, uint256 amount0, uint256 amount1)',
  'event Collect(uint256 indexed tokenId, address recipient, uint256 amount0, uint256 amount1)',
  'function positions(uint256 tokenId) external view returns (uint96 nonce,address operator,address token0,address token1,uint24 fee,int24 tickLower,int24 tickUpper,uint128 liquidity,uint256 feeGrowthInside0LastX128,uint256 feeGrowthInside1LastX128,uint128 tokensOwed0,uint128 tokensOwed1)',
  'function decreaseLiquidity((uint256 tokenId, uint128 liquidity, uint256 amount0Min, uint256 amount1Min, uint256 deadline)) external returns (uint256 amount0, uint256 amount1)',
  'function collect((uint256 tokenId, address recipient, uint128 amount0Max, uint128 amount1Max)) external returns (uint256 amount0, uint256 amount1)',
  'function ownerOf(uint256 tokenId) public view returns (address)',
  'function multicall(bytes[] data) public returns (bytes[] memory results)',
  'function totalSupply() public view returns (uint256)',
]

export const Pool = [
  'function liquidity() external view returns (uint128)',
  'function token0() external view returns (address)',
  'function token1() external view returns (address)',
  'function slot0() external view returns (uint160 sqrtPriceX96,int24 tick,uint16 observationIndex,uint16 observationCardinality,uint16 observationCardinalityNext,uint8 feeProtocol,bool unlocked)',
  'function fee() external view returns (uint24)',
  'event Mint( address sender, address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1)',
  'event Burn( address indexed owner, int24 indexed tickLower, int24 indexed tickUpper, uint128 amount, uint256 amount0, uint256 amount1 )',
  'function factory() external view returns (address)',
]

export type Closed = {
  L1startBlock: any
  L1endBlock: any
  L2startBlock: any
  L2endBlock: any
  data: Info[]
}

export type Info = {
  l1Token: string
  l2Token: string
  tokenName: string
  data: User[]
}

export type User = {
  claimer: string
  amount: string
  type: number // 0 : eoa, 1 : contract(pool)
}

export type WithdrawClaimed = {
  txHash: string
  event: {
    l1Token: any
    l2Token: any
    from: any
    to: any
    amount: any
    data: any
  }
}

// ============================================
// ForceWithdraw Types (Tech Spec v2.2)
// ============================================

/**
 * generate-assets3.json 스키마 (Titan-style)
 * L1 강제 출금 스크립트의 입력 소스
 */
export interface ForceWithdrawAssetEntry {
  l1Token: string // L1 token address (0x0 for ETH)
  l2Token: string // L2 token address
  tokenName: string // Human-readable token name
  data: ForceWithdrawClaim[]
}

export interface ForceWithdrawClaim {
  claimer: string // L1 수신자 주소
  amount: string // wei 단위 string
  hash: string // keccak256(abi.encodePacked(l1Token, claimer, amount))
}

/**
 * Shutdown State (trh-sdk용 상태 추적)
 */
export interface ShutdownState {
  chainId: number
  l1DeploymentJsonPath: string
  lastGenAt?: string // ISO timestamp
  lastSendAt?: string // ISO timestamp
  lastSnapshotFile?: string
  lastDryRunGasEstimate?: string
}

/**
 * UpgradeL1Bridge ABI (Titan Legacy 참조)
 * forceWithdrawClaim 함수 호출용
 */
export const UpgradeL1BridgeABI = [
  'function forceWithdrawClaim(address _position, string calldata _hash, address _token, uint256 _amount, address _address) external',
  'function forceWithdrawClaimAll((address position, string hashed, address token, uint256 amount, address getAddress)[] calldata params) external',
  'function forceRegistry(address[] calldata _position) external',
  'function forceActive(bool _state) external',
  'function setCloser(address _closer) external',
  'function position(address) external view returns (bool)',
  'function claimState(bytes32) external view returns (bool)',
  'function closer() external view returns (address)',
  'event ForceWithdraw(bytes32 indexed _index, address indexed _token, uint256 amount, address _claimer, address _requester)',
]
