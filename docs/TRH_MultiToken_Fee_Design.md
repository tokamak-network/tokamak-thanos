# Multi-Token Fee 설계서 — Thanos L2

> **대상 레포**: `tokamak-thanos` (Gaming/Full Preset)
> **의존 문서**: TRH_AA_TON_Design_v2.md (EIP-7702 + ERC-4337 v0.8)
> **작성일**: 2026.03.16 | Tokamak Network | v1.1
> **v1.1 변경**: fee token 설정을 `trh-sdk deploy` 배포 플로우에 통합. 별도 CLI 커맨드 제거. Preset이 fee token을 자동 결정.

---

## 1. 개요

Thanos L2에서 유저가 TON(네이티브), ETH(L2 ERC-20), USDT, USDC 중 원하는 토큰으로 가스비를 지불할 수 있게 한다.

### 1.1 핵심 원칙

**프로토콜 레벨의 fee token은 TON으로 유지한다.** op-geth의 가스비 결제 로직을 건드리지 않는다. EVM 레이어 변경 없음. OP Stack 호환성 유지. Level 1 (컨트랙트) 변경만으로 구현한다.

유저가 USDC로 가스비를 내는 것처럼 보이지만, 프로토콜 레벨에서는 다음이 발생한다:

```
유저 → USDC 지불 → MultiTokenPaymaster → TON으로 가스비 대납 → EntryPoint → Bundler 보상

    유저가 보는 것: "USDC 0.03 차감"
    프로토콜이 하는 것: "TON 10 소비 (가스비) + USDC 0.03 수거 (환산 금액)"
```

### 1.2 지원 토큰

| 토큰 | L2 주소 | 유형 | 소수점 | 비고 |
|------|---------|------|--------|------|
| TON | 네이티브 (msg.value) | Native | 18 | 가스비 직접 지불. Paymaster 불필요. |
| ETH | 0x4200…0031 (L2 ETH ERC-20) | ERC-20 | 18 | L1 ETH를 브릿지하면 이 ERC-20으로 수령 |
| USDT | 배포 시 결정 | ERC-20 | 6 | USDC Bridge와 유사한 방식으로 L2에 존재 |
| USDC | 배포 시 결정 | ERC-20 | 6 | USDC Bridge predeploy로 L2에 존재 (DeFi/Full Preset) |

### 1.3 유저 경험

```
[게임 설정 화면]

  가스비 지불 수단:
  ┌─────────────────────────┐
  │  ◉ TON  (수수료 없음)    │  ← 네이티브. 직접 지불. 가장 저렴.
  │  ○ ETH  (마크업 5%)      │
  │  ○ USDC (마크업 3%)      │
  │  ○ USDT (마크업 3%)      │
  └─────────────────────────┘

  * 최초 1회 approve 필요 (선택한 토큰)
  * 이후 자동으로 선택한 토큰에서 가스비 차감
```

---

## 2. 아키텍처

### 2.1 v2.0 TonTokenPaymaster → v2.1 MultiTokenPaymaster

AA 설계 v2.0의 TonTokenPaymaster는 L2 ETH만 지원했다. 이를 일반화하여 여러 ERC-20 토큰을 지원하는 MultiTokenPaymaster로 확장한다.

```
v2.0: TonTokenPaymaster
  └── l2Eth (하드코딩) → TON 환산 → 가스비 대납

v2.1: MultiTokenPaymaster
  └── supportedTokens (맵)
       ├── ETH  → ETH/TON 오라클 → 가스비 대납
       ├── USDC → USDC/TON 오라클 → 가스비 대납
       └── USDT → USDT/TON 오라클 → 가스비 대납
```

### 2.2 컴포넌트 구조

```
┌───────────────────────────────────────────────────────────────┐
│  MultiTokenPaymaster                                           │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │  TokenConfig │  │  TokenConfig │  │  TokenConfig │            │
│  │  ETH         │  │  USDC        │  │  USDT        │            │
│  │  oracle: A   │  │  oracle: B   │  │  oracle: C   │            │
│  │  markup: 5%  │  │  markup: 3%  │  │  markup: 3%  │            │
│  │  decimals:18 │  │  decimals: 6 │  │  decimals: 6 │            │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘            │
│         │                │                │                     │
│         ▼                ▼                ▼                     │
│  ┌──────────────────────────────────────────┐                  │
│  │  PriceOracleRouter                        │                  │
│  │                                            │                  │
│  │  ETH/TON  → SimplePriceOracle             │                  │
│  │  USDC/TON → SimplePriceOracle             │                  │
│  │  USDT/TON → SimplePriceOracle             │                  │
│  │                                            │                  │
│  │  (Phase 2: UniswapV3TwapOracle)           │                  │
│  └──────────────────────────────────────────┘                  │
│                                                                 │
│  TON deposit (EntryPoint) ←── 가스비는 항상 TON으로 결제         │
└───────────────────────────────────────────────────────────────┘
```

---

## 3. MultiTokenPaymaster 컨트랙트

### 3.1 핵심 구조

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.23;

import "../core/BasePaymaster.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

/// @title MultiTokenPaymaster
/// @notice 유저가 ETH, USDC, USDT 등 ERC-20 토큰으로 가스비를 지불.
///         Paymaster는 TON(네이티브) deposit에서 가스비를 대납하고,
///         유저의 ERC-20 토큰을 환율에 맞게 수거.
contract MultiTokenPaymaster is BasePaymaster {
    using SafeERC20 for IERC20;

    struct TokenConfig {
        bool enabled;              // 이 토큰이 가스비 수단으로 활성화되었는가
        ITokenPriceOracle oracle;  // 토큰/TON 환율 오라클
        uint256 markupPercent;     // 마크업 (예: 5 = 5%)
        uint8 decimals;            // 토큰 소수점 (ETH=18, USDC=6, USDT=6)
    }

    /// @notice 지원 토큰 목록
    mapping(address => TokenConfig) public supportedTokens;

    /// @notice 토큰 주소 배열 (순회용)
    address[] public tokenList;

    /// @notice 수거된 토큰 잔액 (운영자가 주기적으로 TON으로 스왑하여 재deposit)
    mapping(address => uint256) public collectedFees;

    event TokenAdded(address indexed token, address oracle, uint256 markupPercent, uint8 decimals);
    event TokenRemoved(address indexed token);
    event TokenConfigUpdated(address indexed token, address oracle, uint256 markupPercent);
    event FeesCollected(address indexed token, address indexed sender, uint256 amount);
    event FeesWithdrawn(address indexed token, address indexed to, uint256 amount);

    constructor(IEntryPoint _entryPoint) BasePaymaster(_entryPoint) {}

    // ═══════════════════════════════════════════
    // 관리 함수 (운영자)
    // ═══════════════════════════════════════════

    /// @notice 새 토큰을 가스비 수단으로 추가
    function addToken(
        address token,
        ITokenPriceOracle oracle,
        uint256 markupPercent,
        uint8 decimals
    ) external onlyOwner {
        require(!supportedTokens[token].enabled, "already enabled");
        require(address(oracle) != address(0), "zero oracle");
        require(markupPercent <= 50, "markup too high"); // 최대 50%

        supportedTokens[token] = TokenConfig({
            enabled: true,
            oracle: oracle,
            markupPercent: markupPercent,
            decimals: decimals
        });
        tokenList.push(token);

        emit TokenAdded(token, address(oracle), markupPercent, decimals);
    }

    /// @notice 토큰을 가스비 수단에서 제거
    function removeToken(address token) external onlyOwner {
        require(supportedTokens[token].enabled, "not enabled");
        supportedTokens[token].enabled = false;
        emit TokenRemoved(token);
    }

    /// @notice 토큰 설정 변경 (오라클, 마크업)
    function updateTokenConfig(
        address token,
        ITokenPriceOracle oracle,
        uint256 markupPercent
    ) external onlyOwner {
        require(supportedTokens[token].enabled, "not enabled");
        require(markupPercent <= 50, "markup too high");

        supportedTokens[token].oracle = oracle;
        supportedTokens[token].markupPercent = markupPercent;

        emit TokenConfigUpdated(token, address(oracle), markupPercent);
    }

    /// @notice 수거된 토큰을 인출 (운영자가 스왑 후 재deposit에 사용)
    function withdrawCollectedFees(
        address token,
        address to,
        uint256 amount
    ) external onlyOwner {
        require(collectedFees[token] >= amount, "insufficient collected");
        collectedFees[token] -= amount;
        IERC20(token).safeTransfer(to, amount);
        emit FeesWithdrawn(token, to, amount);
    }

    // ═══════════════════════════════════════════
    // Paymaster 핵심 로직
    // ═══════════════════════════════════════════

    /// @notice UserOp 검증: 토큰 allowance 확인 + 사전 수거
    function _validatePaymasterUserOp(
        PackedUserOperation calldata userOp,
        bytes32 userOpHash,
        uint256 maxCost              // ← TON 단위 (네이티브)
    ) internal override returns (bytes memory context, uint256 validationData) {

        // 1. paymasterAndData에서 토큰 주소 추출
        //    paymasterAndData = [paymaster_address(20)] [token_address(20)] [signature(...)]
        address token = address(bytes20(userOp.paymasterAndData[20:40]));

        // 2. 토큰이 지원되는지 확인
        TokenConfig memory config = supportedTokens[token];
        require(config.enabled, "PM: token not supported");

        // 3. maxCost(TON) → 선택된 토큰으로 환산
        uint256 tokenCost = _tonToToken(maxCost, token, config);
        uint256 tokenCostWithMarkup = tokenCost * (100 + config.markupPercent) / 100;

        // 4. 유저의 토큰 allowance 확인
        address sender = userOp.getSender();
        uint256 allowance = IERC20(token).allowance(sender, address(this));
        require(allowance >= tokenCostWithMarkup, "PM: insufficient allowance");

        // 5. 토큰 사전 수거 (pre-charge)
        IERC20(token).safeTransferFrom(sender, address(this), tokenCostWithMarkup);

        // 6. context에 정산 정보 기록
        context = abi.encode(sender, token, tokenCostWithMarkup, maxCost);
        validationData = 0;
    }

    /// @notice 실행 후 정산: 실제 가스비만큼 수거하고 나머지 환불
    function _postOp(
        PostOpMode mode,
        bytes calldata context,
        uint256 actualGasCost,       // ← TON 단위
        uint256 actualUserOpFeePerGas
    ) internal override {
        (
            address sender,
            address token,
            uint256 preCharged,
            uint256 maxCost
        ) = abi.decode(context, (address, address, uint256, uint256));

        TokenConfig memory config = supportedTokens[token];

        // 1. 실제 가스비(TON) → 선택된 토큰으로 환산
        uint256 actualTokenCost = _tonToToken(actualGasCost, token, config);
        uint256 actualWithMarkup = actualTokenCost * (100 + config.markupPercent) / 100;

        // 2. 과다 수거분 환불
        if (preCharged > actualWithMarkup) {
            IERC20(token).safeTransfer(sender, preCharged - actualWithMarkup);
        }

        // 3. 실제 수거 금액 기록
        uint256 actualCollected = preCharged > actualWithMarkup
            ? actualWithMarkup
            : preCharged;
        collectedFees[token] += actualCollected;

        emit FeesCollected(token, sender, actualCollected);
    }

    // ═══════════════════════════════════════════
    // 환율 변환
    // ═══════════════════════════════════════════

    /// @notice TON → 토큰 환율 변환 (소수점 처리 포함)
    function _tonToToken(
        uint256 tonAmount,
        address token,
        TokenConfig memory config
    ) internal view returns (uint256) {
        // oracle.getPrice()는 1 TON당 해당 토큰의 가격 (18 decimals 고정)
        // 예: TON/USDC = 0.65e18 → 1 TON = 0.65 USDC
        uint256 price = config.oracle.getPrice();

        if (config.decimals == 18) {
            // ETH: 18 decimals → 직접 환산
            return tonAmount * price / 1e18;
        } else {
            // USDC/USDT: 6 decimals → 스케일링 필요
            // tonAmount (18 dec) * price (18 dec) / 1e18 → 18 dec 결과
            // 18 dec → 6 dec: / 1e12
            return tonAmount * price / 1e18 / (10 ** (18 - config.decimals));
        }
    }

    /// @notice 토큰 → TON 환율 변환 (가스 추정용 view 함수)
    function getTokenToTonRate(address token) external view returns (uint256) {
        TokenConfig memory config = supportedTokens[token];
        require(config.enabled, "token not supported");
        return config.oracle.getPrice();
    }

    /// @notice 예상 토큰 가스비 조회 (프론트엔드용)
    function estimateTokenCost(
        address token,
        uint256 estimatedTonGasCost
    ) external view returns (uint256 tokenCost, uint256 tokenCostWithMarkup) {
        TokenConfig memory config = supportedTokens[token];
        require(config.enabled, "token not supported");

        tokenCost = _tonToToken(estimatedTonGasCost, token, config);
        tokenCostWithMarkup = tokenCost * (100 + config.markupPercent) / 100;
    }
}
```

### 3.2 소수점 처리 상세

가장 까다로운 부분이 토큰별 소수점 차이다. TON/ETH는 18 decimals, USDC/USDT는 6 decimals.

```
예시: 가스비 10 TON, TON 가격 = 0.65 USDC

잘못된 계산:
  10e18 (TON wei) × 0.65e18 (price) / 1e18 = 6.5e18
  → 이것은 6.5 USDC가 아니라 6,500,000,000,000 USDC (6.5e12 USDC)

올바른 계산:
  10e18 (TON wei) × 0.65e18 (price) / 1e18 / 1e12 (18→6 dec) = 6.5e6
  → 6.5e6 = 6.5 USDC ✓
```

| 토큰 | decimals | 스케일링 | 예시 (10 TON 가스비) |
|------|----------|---------|---------------------|
| ETH | 18 | / 1e18 | 10 TON × 0.0005 = 0.005 ETH |
| USDC | 6 | / 1e18 / 1e12 | 10 TON × 0.65 = 6.5 USDC |
| USDT | 6 | / 1e18 / 1e12 | 10 TON × 0.65 = 6.5 USDT |

---

## 4. 오라클 설계

### 4.1 토큰별 오라클

각 토큰에 독립적인 오라클을 배치한다. 인터페이스는 AA 설계 v2.0의 ITokenPriceOracle과 동일하다.

```solidity
interface ITokenPriceOracle {
    /// @notice 1 TON당 해당 토큰의 가격 (18 decimals 고정)
    /// @return price 예: USDC 오라클이면 0.65e18 = 1 TON = 0.65 USDC
    function getPrice() external view returns (uint256 price);

    /// @notice 가격 업데이트 시간
    function lastUpdated() external view returns (uint256);
}
```

| 오라클 | 토큰 | Phase 1 | Phase 2 |
|--------|------|---------|---------|
| ETH/TON | ETH | SimplePriceOracle (수동) | UniswapV3TwapOracle (TON/WTON 풀) |
| USDC/TON | USDC | SimplePriceOracle (수동) | UniswapV3TwapOracle (TON/USDC 풀) |
| USDT/TON | USDT | SimplePriceOracle (수동) | UniswapV3TwapOracle (TON/USDT 풀) |

### 4.2 Phase 1: SimplePriceOracle (수동 업데이트)

```solidity
contract SimplePriceOracle is ITokenPriceOracle {
    uint256 public override lastUpdated;
    uint256 private _price;
    address public owner;
    uint256 public constant MAX_STALENESS = 86400; // 24시간

    constructor(uint256 initialPrice) {
        owner = msg.sender;
        _price = initialPrice;
        lastUpdated = block.timestamp;
    }

    function getPrice() external view override returns (uint256) {
        require(block.timestamp - lastUpdated < MAX_STALENESS, "stale price");
        return _price;
    }

    function updatePrice(uint256 newPrice) external {
        require(msg.sender == owner, "only owner");
        _price = newPrice;
        lastUpdated = block.timestamp;
    }
}
```

운영자가 외부 소스(CoinMarketCap API, Binance API 등)에서 가격을 가져와 주기적으로 업데이트한다. SentinAI 또는 cron job으로 자동화 가능.

### 4.3 Phase 2: UniswapV3TwapOracle (자동)

DeFi/Full Preset에는 Uniswap V3가 genesis predeploy로 포함되어 있다. TON/USDC, TON/ETH 풀이 존재하면 TWAP(Time-Weighted Average Price)을 자동으로 읽을 수 있다.

```solidity
contract UniswapV3TwapOracle is ITokenPriceOracle {
    IUniswapV3Pool public pool;
    uint32 public twapInterval = 1800; // 30분 TWAP

    function getPrice() external view override returns (uint256) {
        (int24 arithmeticMeanTick, ) = OracleLibrary.consult(
            address(pool), twapInterval
        );
        return OracleLibrary.getQuoteAtTick(
            arithmeticMeanTick, 1e18, token0, token1
        );
    }
}
```

별도 외부 트리거 불필요. 온체인 자동화. 다만 풀의 유동성이 충분해야 TWAP이 조작에 강해진다.

### 4.4 오라클 안전장치

| 안전장치 | 설명 |
|---------|------|
| Stale price 체크 | 24시간 미갱신 시 revert. 유저의 토큰이 잘못된 환율로 수거되는 것을 방지. |
| 가격 변동폭 제한 | 이전 가격 대비 ±30% 초과 변동 시 revert. 오라클 조작 방지. |
| 마크업 버퍼 | 마크업(3~5%)이 환율 변동의 버퍼 역할. 약간의 환율 변동은 마크업으로 흡수. |
| Fallback | TWAP 오라클 실패 시 SimplePriceOracle로 자동 fallback. |

---

## 5. 배포 시 자동 설정 (trh-sdk deploy 통합)

### 5.1 핵심 원칙: 배포 = fee token 설정 완료

fee token 설정은 `trh-sdk deploy` 한 줄에 통합된다. 별도 커맨드 불필요. Preset이 fee token을 결정하고, SDK가 배포 과정에서 자동으로 오라클 배포, 토큰 등록, 가격 초기화, CronJob 등록까지 수행한다.

```bash
# 이 한 줄이 fee token 설정까지 포함
trh-sdk deploy --chain-name my-game --preset gaming --network testnet
```

운영자 입력: **0개.** Preset 선택이 곧 fee token 결정이다.

### 5.2 Preset → Fee Token 매핑

| Preset | 자동 등록 Fee Token | 이유 |
|--------|-------------------|------|
| General | TON만 | MultiTokenPaymaster 미배포. 가장 가벼운 구성. |
| DeFi | TON + ETH + USDC + USDT | USDC Bridge predeploy 존재. DeFi 유저는 스테이블코인 보유. 4개 모두 자동. |
| Gaming | TON + ETH | AA + Paymaster 포함. USDC Bridge 미포함이므로 스테이블코인은 미등록. |
| Full | TON + ETH + USDC + USDT | 모든 predeploy 포함. 4개 모두 자동. |

### 5.3 SDK 내부 실행 순서

`trh-sdk deploy --preset gaming` 실행 시 SDK가 자동으로 수행하는 작업:

```
Phase 1: L1 컨트랙트 배포 + Genesis 생성 (기존 플로우)
Phase 2: Helm install — L2 노드 구동 (기존 플로우)
Phase 3: 모듈 배포 (기존 플로우 — Explorer, Bridge 등)

Phase 4: Fee Token 자동 설정 (신규)
  │
  ├── 4-1. MultiTokenPaymaster 배포
  │         → deployer 키로 L2에 배포
  │         → 배포 주소를 deployments.json에 기록
  │
  ├── 4-2. Preset에 정의된 토큰별 오라클 배포
  │         ├── ETH/TON SimplePriceOracle 배포
  │         ├── (DeFi/Full만) USDC/TON SimplePriceOracle 배포
  │         └── (DeFi/Full만) USDT/TON SimplePriceOracle 배포
  │
  ├── 4-3. 초기 가격 설정
  │         → SDK가 CoinGecko/Binance API에서 현재 환율 조회
  │         → 각 오라클에 updatePrice() 호출
  │         → 실패 시 하드코딩된 fallback 가격 사용:
  │            ETH/TON = 2600.0, USDC/TON = 0.65, USDT/TON = 0.65
  │
  ├── 4-4. 토큰 등록
  │         → paymaster.addToken(ETH, ethOracle, 5, 18)
  │         → paymaster.addToken(USDC, usdcOracle, 3, 6)  // DeFi/Full만
  │         → paymaster.addToken(USDT, usdtOracle, 3, 6)  // DeFi/Full만
  │
  ├── 4-5. Paymaster에 TON 초기 deposit
  │         → Testnet: Tokamak Faucet에서 자동 수령 → deposit
  │         → Mainnet: 펀딩 도우미에서 필요 금액 안내 → 유저 입금 후 자동 deposit
  │
  └── 4-6. 가격 업데이트 CronJob 등록
            → K8s CronJob: 1시간마다 CoinGecko API → oracle.updatePrice()
            → Helm values에 CronJob 스펙 포함 (tokamak-thanos-stack)
            → SentinAI 연동: 가격 업데이트 실패 시 알림
```

### 5.4 SDK 내부 구현 (Go)

```go
// pkg/modules/fee_token.go

type FeeTokenConfig struct {
    Token    common.Address
    Symbol   string
    Decimals uint8
    Markup   uint256
    Oracle   common.Address  // 배포 후 채워짐
}

// Preset별 fee token 정의
var PresetFeeTokens = map[string][]FeeTokenConfig{
    "general": {},  // MultiTokenPaymaster 미배포
    "defi": {
        {Symbol: "ETH",  Decimals: 18, Markup: 5},
        {Symbol: "USDC", Decimals: 6,  Markup: 3},
        {Symbol: "USDT", Decimals: 6,  Markup: 3},
    },
    "gaming": {
        {Symbol: "ETH",  Decimals: 18, Markup: 5},
    },
    "full": {
        {Symbol: "ETH",  Decimals: 18, Markup: 5},
        {Symbol: "USDC", Decimals: 6,  Markup: 3},
        {Symbol: "USDT", Decimals: 6,  Markup: 3},
    },
}

func SetupFeeTokens(ctx *DeployContext) error {
    presetID := ctx.Config.PresetID
    tokens := PresetFeeTokens[presetID]
    
    if len(tokens) == 0 {
        // General Preset: MultiTokenPaymaster 미배포
        return nil
    }
    
    // 1. MultiTokenPaymaster 배포
    paymaster, err := DeployMultiTokenPaymaster(ctx)
    if err != nil { return err }
    
    for i, token := range tokens {
        // 2. 토큰 주소 결정 (genesis predeploy 또는 deployments.json에서)
        tokenAddr := resolveTokenAddress(ctx, token.Symbol)
        
        // 3. 오라클 배포 + 초기 가격
        price := fetchPrice(token.Symbol)  // CoinGecko API
        oracle, err := DeploySimplePriceOracle(ctx, price)
        if err != nil { return err }
        tokens[i].Oracle = oracle
        
        // 4. 토큰 등록
        err = paymaster.AddToken(tokenAddr, oracle, token.Markup, token.Decimals)
        if err != nil { return err }
    }
    
    // 5. Paymaster TON deposit
    err = DepositToPaymaster(ctx, paymaster)
    if err != nil { return err }
    
    // 6. CronJob 등록 (Helm values에 추가)
    ctx.HelmValues.Set("feeToken.cronJob.enabled", true)
    ctx.HelmValues.Set("feeToken.cronJob.oracleAddresses", oracleAddresses)
    
    return nil
}
```

### 5.5 배포 후 토큰 추가/제거 (Platform UI)

배포 시 자동 설정된 fee token 외에 추가/제거가 필요하면 Platform UI에서 한다. CLI 커맨드는 없다.

```
[Chain Dashboard] → [Fee Tokens]

  ┌────────────────────────────────────────────────────────────────┐
  │  Fee Token 설정                                                 │
  │                                                                  │
  │  ✅ TON  (네이티브)        마크업: -    상태: 항상 활성           │
  │  ✅ ETH  0x4200…0031      마크업: 5%   상태: 활성  [편집] [해제] │
  │  ✅ USDC 0x1234…5678      마크업: 3%   상태: 활성  [편집] [해제] │
  │  ✅ USDT 0xabcd…ef01      마크업: 3%   상태: 활성  [편집] [해제] │
  │                                                                  │
  │  [+ 토큰 추가]                                                   │
  │                                                                  │
  │  ── Paymaster 상태 ──────────────────────────────────────────── │
  │  TON deposit 잔액: 1,234.5 TON                                  │
  │  수거된 수수료:  ETH 0.5 | USDC 320.4 | USDT 180.2              │
  │  가격 업데이트:  마지막 12분 전 (정상)                            │
  │                                                                  │
  │  [수거된 수수료 인출]  [TON 재충전]  [오라클 가격 갱신]           │
  └────────────────────────────────────────────────────────────────┘
```

Platform UI의 [+ 토큰 추가]가 내부적으로 수행하는 작업:
1. 운영자가 토큰 주소, 마크업 입력
2. Backend가 SimplePriceOracle 배포 + 초기 가격 설정
3. `paymaster.addToken()` 호출
4. CronJob에 새 오라클 주소 추가

### 5.6 MCP AI에서 fee token 추가

MCP Server의 기존 체인 관리 Tool을 통해서도 가능하다:

```
유저: "내 체인에 DAI를 가스비 수단으로 추가해줘"
AI: DAI 토큰 주소를 확인했습니다. 마크업은 기본 3%로 설정합니다.
    DAI를 fee token으로 추가할까요?
유저: "응"
AI: → Backend API: POST /api/chains/{id}/fee-tokens
    → SimplePriceOracle 배포 + addToken() 호출
    DAI가 가스비 수단으로 추가되었습니다.
```

---

## 6. 정산 및 수익 구조

### 6.1 수수료 흐름

```
유저                    MultiTokenPaymaster              운영자
  │                           │                            │
  ├── USDC 0.68 지불 ──────→ │                            │
  │   (가스비 0.65 + 마크업 0.03)                          │
  │                           ├── TON 10 대납 ──→ EntryPoint
  │                           │                            │
  │                           ├── USDC 0.68 수거 ──────→  │
  │                           │   (collectedFees에 축적)   │
  │                           │                            │
  │                           │         주기적으로          │
  │                           │   ←── withdrawCollectedFees │
  │                           │                            │
  │                           │         운영자가            │
  │                           │   USDC → TON 스왑 (Uniswap) │
  │                           │                            │
  │                           │   ←── TON re-deposit ──────│
  │                           │       (paymaster.deposit)   │
```

### 6.2 운영자 수익

| 항목 | 설명 |
|------|------|
| 마크업 수수료 | 가스비의 3~5%. 운영자의 직접 수익. |
| 환율 차익 | 오라클 가격과 실제 스왑 가격의 차이. 마크업이 이를 커버. |
| TON deposit 이자 | EntryPoint에 deposit된 TON은 인출 가능. 별도 이자 없음. |

### 6.3 자동 재충전

MultiTokenPaymaster의 TON deposit이 소진되면 가스비 대납이 불가능해진다. 자동 재충전 메커니즘:

```
SentinAI 모니터링:
  - MultiTokenPaymaster의 TON deposit 잔액 모니터링
  - 임계값(예: 100 TON) 이하 시 알림
  - 자동 실행: collectedFees에서 토큰 인출 → Uniswap 스왑 → TON deposit

또는 K8s CronJob:
  - 6시간마다 수거된 USDC/USDT/ETH를 TON으로 스왑
  - 스왑된 TON을 Paymaster에 재deposit
```

---

## 7. SDK / 프론트엔드 연동

### 7.1 가스비 토큰 선택

```typescript
// 지원 토큰 목록 조회
const tokens = await paymaster.getEnabledTokens();
// → [{ address: "0x4200...0031", symbol: "ETH", markup: 5 },
//    { address: "0x1234...5678", symbol: "USDC", markup: 3 }, ...]

// 예상 가스비 조회 (토큰별)
const tonGasEstimate = await entryPoint.estimateGas(userOp);

const estimates = await Promise.all(tokens.map(async (token) => {
  const { tokenCostWithMarkup } = await paymaster.estimateTokenCost(
    token.address,
    tonGasEstimate
  );
  return { token, cost: tokenCostWithMarkup };
}));

// → [{ token: "ETH",  cost: "0.005 ETH"  },
//    { token: "USDC", cost: "6.50 USDC"  },
//    { token: "USDT", cost: "6.50 USDT"  }]
```

### 7.2 UserOp 생성

```typescript
// 유저가 USDC를 선택한 경우
const selectedToken = USDC_ADDRESS;

// 1. 최초 1회: approve (무한 approve 또는 세션별)
if (!hasApproval) {
  await usdc.approve(PAYMASTER_ADDRESS, ethers.MaxUint256);
}

// 2. paymasterAndData 구성
// [paymaster_address(20bytes)] [token_address(20bytes)] [signature(65bytes)]
const paymasterAndData = ethers.concat([
  PAYMASTER_ADDRESS,
  selectedToken,
  await paymasterSigner.signPaymasterData(userOp, selectedToken),
]);

// 3. UserOp에 paymasterAndData 설정
userOp.paymasterAndData = paymasterAndData;

// 4. Bundler에 제출
await bundler.sendUserOperation(userOp, ENTRY_POINT);
```

### 7.3 TON으로 직접 지불하는 경우

TON으로 가스비를 내는 경우 Paymaster가 필요 없다. 유저의 EOA(7702 delegation 활성)에서 직접 가스비를 TON으로 지불한다.

```typescript
// paymasterAndData를 비움 → Paymaster 미사용 → 유저의 TON deposit에서 가스비 차감
userOp.paymasterAndData = "0x";
```

---

## 8. 테스트 시나리오

### T1: ETH로 가스비 지불
1. 유저 EOA에 7702 delegation 설정
2. L2 ETH를 MultiTokenPaymaster에 approve
3. UserOp 전송 (paymasterAndData에 ETH 토큰 주소)
4. **검증**: 유저의 L2 ETH가 가스비+마크업(5%)만큼 차감
5. **검증**: Paymaster의 TON deposit이 가스비만큼 차감
6. **검증**: `collectedFees[ETH]`가 수거액만큼 증가

### T2: USDC로 가스비 지불 (6 decimals)
1. USDC를 MultiTokenPaymaster에 approve
2. UserOp 전송 (paymasterAndData에 USDC 주소)
3. **검증**: 소수점 변환이 정확한지 확인 (18→6 decimals)
4. **검증**: 유저의 USDC가 올바른 금액만큼 차감

### T3: USDT로 가스비 지불
1. T2와 동일한 플로우를 USDT로 실행
2. **검증**: USDT 특유의 approve(0) 요구사항 처리 (SafeERC20)

### T4: 지원되지 않는 토큰으로 시도
1. 등록되지 않은 토큰 주소로 UserOp 전송
2. **검증**: "PM: token not supported" revert

### T5: 오라클 장애
1. oracle.updatePrice를 24시간 이상 안 함
2. UserOp 제출
3. **검증**: "stale price" revert
4. **검증**: 유저의 토큰이 차감되지 않음

### T6: TON 직접 지불 (Paymaster 미사용)
1. paymasterAndData를 비운 UserOp 전송
2. **검증**: 유저의 TON deposit에서 가스비 차감
3. **검증**: ERC-20 토큰 변화 없음

### T7: 토큰 추가/제거
1. 운영자가 `addToken(DAI, oracle, 4, 18)` 호출
2. DAI로 가스비 지불 → 성공
3. 운영자가 `removeToken(DAI)` 호출
4. DAI로 가스비 지불 → "PM: token not supported" revert

### T8: 수거된 수수료 인출
1. 여러 UserOp 실행 → USDC 수거됨
2. 운영자가 `withdrawCollectedFees(USDC, operatorAddress, amount)` 호출
3. **검증**: 운영자 주소에 USDC 수신

### T9: 마크업 정확성
1. 마크업 3%로 설정된 USDC
2. 가스비 10 TON, TON/USDC = 0.65
3. 예상: 6.5 USDC × 1.03 = 6.695 USDC
4. **검증**: 사전 수거 금액이 6.695 USDC (6 decimals: 6695000)
5. **검증**: 실제 사용량이 적으면 차액 환불

### T10: Preset별 deploy 시 fee token 자동 설정
1. `trh-sdk deploy --preset gaming --network testnet` 실행
2. **검증**: MultiTokenPaymaster가 L2에 배포됨
3. **검증**: ETH/TON SimplePriceOracle이 배포되고 초기 가격이 설정됨
4. **검증**: `paymaster.supportedTokens(ETH)` → enabled=true, markup=5
5. **검증**: USDC/USDT는 미등록 (Gaming Preset)
6. **검증**: Paymaster의 TON deposit > 0 (Faucet 자동 충전)
7. **검증**: 가격 업데이트 CronJob이 K8s에 등록됨
8. DeFi Preset으로 같은 테스트 반복
9. **검증**: ETH + USDC + USDT 3개 모두 자동 등록됨

### T11: 배포 후 Platform UI에서 토큰 추가
1. Gaming Preset으로 배포된 체인
2. Platform UI에서 [+ 토큰 추가] → USDC 주소, 마크업 3% 입력
3. **검증**: USDC/TON 오라클 배포됨
4. **검증**: `paymaster.supportedTokens(USDC)` → enabled=true
5. USDC로 가스비 지불 UserOp 전송 → 성공

---

## 9. 보안 고려사항

### 9.1 감사 필요 포인트

| 포인트 | 심각도 | 상세 |
|--------|--------|------|
| 소수점 변환 정밀도 | **Critical** | 18→6 decimals 변환 시 rounding error. 특히 매우 작은 가스비에서 0으로 내림되면 무료 TX 가능성. 최소 수거 금액(minCharge) 설정 필요. |
| Pre-charge 동결 위험 | **High** | transferFrom 후 postOp에서 환불 실패 시 유저 자금 동결. PostOpMode.opReverted 케이스에서도 환불이 보장되어야 함. |
| 오라클 조작 | **High** | 환율 조작 시 arbitrage 가능. stale price 체크 + 가격 변동폭 제한. Phase 2에서 TWAP으로 교체. |
| USDT approve 특이사항 | **Medium** | USDT는 approve(0)을 먼저 호출해야 새 값으로 approve 가능. SafeERC20의 forceApprove 사용. |
| Reentrancy | **Medium** | ERC-20 transferFrom → postOp 사이 reentrancy. L2 predeploy ERC-20은 hook 없으므로 위험 낮지만, 외부 토큰 등록 시 주의. check-effects-interaction 패턴 준수. |
| 무한 approve 위험 | **Low** | 유저가 MaxUint256로 approve하면 Paymaster 컨트랙트가 탈취당할 경우 전체 잔액 위험. 세션별 approve 또는 금액 제한 approve 권장. |

### 9.2 최소 수거 금액 (rounding 방어)

```solidity
// 소수점 변환 후 0이 되는 것을 방지
uint256 public minChargePerToken; // 토큰별 최소 수거 금액

// _tonToToken 결과가 0이면 최소 수거 금액 적용
function _tonToTokenSafe(...) internal view returns (uint256) {
    uint256 result = _tonToToken(tonAmount, token, config);
    return result > 0 ? result : minChargePerToken;
}
```

---

## 10. Preset별 기본 설정 및 배포 자동화

### 10.1 Preset → Fee Token 매트릭스

| Preset | TON | ETH | USDC | USDT | MultiTokenPaymaster | 오라클 수 | CronJob |
|--------|-----|-----|------|------|---------------------|----------|---------|
| General | ✅ (항상) | ❌ | ❌ | ❌ | 미배포 | 0 | 미등록 |
| DeFi | ✅ | ✅ | ✅ | ✅ | 배포 | 3 | 등록 |
| Gaming | ✅ | ✅ | ❌ | ❌ | 배포 | 1 | 등록 |
| Full | ✅ | ✅ | ✅ | ✅ | 배포 | 3 | 등록 |

Gaming Preset에서 USDC/USDT가 미등록인 이유: Gaming Preset genesis에는 USDC Bridge predeploy가 포함되지 않으므로, L2에 USDC/USDT가 존재하지 않을 수 있다. 운영자가 이후 Platform UI에서 토큰을 추가할 수 있다.

### 10.2 Preset별 Helm values (자동 생성)

SDK가 `--preset` 플래그에 따라 아래 Helm values를 자동 생성한다. 운영자 입력 없음.

```yaml
# values-gaming.yaml (SDK가 자동 생성)
feeToken:
  enabled: true
  paymaster:
    address: ""           # 배포 후 자동 채움
    initialDeposit: "100" # TON (Testnet 기본)
  tokens:
    - symbol: ETH
      address: "0x4200000000000000000000000000000000000031"
      decimals: 18
      markup: 5
      oracle:
        type: simple      # Phase 1
        address: ""        # 배포 후 자동 채움
  cronJob:
    enabled: true
    schedule: "0 * * * *"  # 1시간마다
    priceSource: "coingecko"
    # failover: "binance"

# values-defi.yaml (SDK가 자동 생성)
feeToken:
  enabled: true
  paymaster:
    address: ""
    initialDeposit: "500"
  tokens:
    - symbol: ETH
      address: "0x4200000000000000000000000000000000000031"
      decimals: 18
      markup: 5
      oracle: { type: simple, address: "" }
    - symbol: USDC
      address: ""          # deployments.json에서 자동 추출 (USDC Bridge predeploy)
      decimals: 6
      markup: 3
      oracle: { type: simple, address: "" }
    - symbol: USDT
      address: ""          # deployments.json에서 자동 추출
      decimals: 6
      markup: 3
      oracle: { type: simple, address: "" }
  cronJob:
    enabled: true
    schedule: "0 * * * *"

# values-general.yaml
feeToken:
  enabled: false           # General Preset: Paymaster 미배포
```

### 10.3 배포 결과 출력

```
$ trh-sdk deploy --chain-name my-game --preset gaming --network testnet

[Phase 1] L1 컨트랙트 배포 ........................ ✅ (3분)
[Phase 2] Genesis 생성 ........................... ✅ (1분)
[Phase 3] Helm install (L2 노드) ................. ✅ (5분)
[Phase 4] Explorer + Bridge 배포 ................. ✅ (3분)
[Phase 5] Monitoring 스택 배포 ................... ✅ (3분)
[Phase 6] Fee Token 설정 ........................ ✅ (2분)
          ├── MultiTokenPaymaster 배포완료
          │   └── 주소: 0x7a3b...4c2d
          ├── ETH/TON 오라클 배포완료
          │   └── 초기 가격: 1 TON = 0.000385 ETH
          ├── ETH fee token 등록완료 (마크업 5%)
          ├── Paymaster TON deposit: 100 TON (Faucet)
          └── 가격 업데이트 CronJob 등록완료 (1시간 간격)
[Phase 7] DRB 노드 연동 ......................... ✅ (3분)

═══════════════════════════════════════════════════════════
  체인 배포 완료: my-game (Gaming Preset, Testnet)
  총 소요 시간: 20분

  엔드포인트:
    L2 RPC:     https://rpc.my-game.tokamak.network
    Explorer:   https://explorer.my-game.tokamak.network
    Bridge:     https://bridge.my-game.tokamak.network

  Fee Token:
    ✅ TON (네이티브, 직접 지불)
    ✅ ETH (Paymaster, 마크업 5%)

  Paymaster:
    주소: 0x7a3b...4c2d
    TON 잔액: 100 TON
    가격 업데이트: CronJob 활성 (1시간 간격)
═══════════════════════════════════════════════════════════
```

### 10.4 Mainnet 배포 시 차이

Testnet에서는 Faucet으로 Paymaster의 TON 초기 deposit을 자동 충전한다. Mainnet에서는 펀딩 도우미가 안내한다:

```
[Phase 6] Fee Token 설정
          ├── MultiTokenPaymaster 배포완료
          ├── 오라클 3개 배포완료 (ETH, USDC, USDT)
          │
          ├── ⚠️  Paymaster TON deposit 필요
          │   필요 금액: 최소 500 TON (권장 2,000 TON)
          │   입금 주소: 0x7a3b...4c2d
          │   
          │   QR 코드: [████████]
          │   
          │   입금 확인 중... (10초마다 폴링)
          │   ✅ 2,000 TON 입금 확인!
          │
          └── 가격 업데이트 CronJob 등록완료
```

이것은 기존 TRH 설계의 "펀딩 도우미 (Plan B)" 패턴과 동일하다. Operator EOA 펀딩과 Paymaster 펀딩이 같은 플로우에서 처리된다.

---

## 11. 구현 일정

### 11.1 SDK 통합 중심 일정

| 단계 | 기간 | 작업 | 전제조건 |
|------|------|------|----------|
| 11-1 | 1주 | MultiTokenPaymaster 컨트랙트 개발. 토큰 맵, 소수점 변환, addToken/removeToken. | AA 설계 v2.0 Phase 13-2 완료 |
| 11-2 | 1주 | `pkg/modules/fee_token.go` 구현. Preset→fee token 매핑, 오라클 배포, 토큰 등록, CronJob Helm values 생성. 이 모듈이 deploy 플로우의 Phase 6으로 자동 실행. | 11-1 완료 |
| 11-3 | 0.5주 | `tokamak-thanos-stack`에 가격 업데이트 CronJob chart 추가. Helm template에 `feeToken.cronJob` 조건부 렌더링. | 11-2 완료 |
| 11-4 | 0.5주 | trh-backend에 `POST /api/chains/{id}/fee-tokens` API 추가. Platform UI에서 토큰 추가/제거/마크업 변경. | 11-2 완료 |
| 11-5 | 1주 | 프론트엔드 토큰 선택 UI. estimateTokenCost 연동. approve 플로우. | 11-4 완료 |
| 11-6 | 1주 | E2E 테스트 T1~T10 전수 실행. 소수점 edge case. Preset별 deploy→fee token 자동 설정 검증. | 11-5 완료 |
| 11-7 | 병렬 | MultiTokenPaymaster 외부 감사 | 11-1 완료 시 시작 |

**총 예상: 5주 (감사 병렬)**

AA 설계 v2.0의 Phase 13-3 (TonTokenPaymaster 2주)을 MultiTokenPaymaster + SDK 통합으로 확장하므로, **AA v2.0 일정에 +3주**로 계산한다.

### 11.2 변경되는 레포별 작업

| 레포 | 작업 | 담당 |
|------|------|------|
| tokamak-thanos | MultiTokenPaymaster.sol, 오라클 컨트랙트 | Solidity 개발 |
| tokamak-thanos-stack | CronJob chart, Helm values에 feeToken 섹션 추가 | Helm/K8s 개발 |
| trh-sdk | `pkg/modules/fee_token.go`, deploy 플로우에 Phase 6 추가, Preset 정의에 fee token 매핑 추가 | Go 개발 |
| trh-backend | `POST /api/chains/{id}/fee-tokens` API, fee token 상태 조회 API | Go 개발 |
| trh-platform-ui | Fee Token 설정 패널, 토큰 추가/제거 UI | React/TS 개발 |

---

## 12. 변경 파일 총정리

| 파일 | 레포 | 변경 수준 | 내용 | 감사 |
|------|------|----------|------|------|
| **MultiTokenPaymaster.sol** | tokamak-thanos | **신규** | 다중 토큰 가스비 정산. 핵심. | **필수** |
| ITokenPriceOracle.sol | tokamak-thanos | 동일 | 인터페이스 변경 없음 | 필수 |
| SimplePriceOracle.sol | tokamak-thanos | 동일 | 토큰별 1개씩 배포 | 필수 |
| UniswapV3TwapOracle.sol | tokamak-thanos | 신규 (Phase 2) | Uniswap V3 TWAP 기반 자동 오라클 | 필수 |
| **pkg/modules/fee_token.go** | trh-sdk | **신규** | Preset→fee token 매핑, 오라클 배포, 토큰 등록 자동화 | 불필요 |
| **pkg/preset/presets.go** | trh-sdk | **수정** | 각 Preset에 FeeTokens 필드 추가 | 불필요 |
| commands/deploy.go | trh-sdk | 수정 | deploy 플로우에 Phase 6 (fee token 설정) 추가 | 불필요 |
| **charts/fee-token-updater/** | tokamak-thanos-stack | **신규** | 가격 업데이트 CronJob chart | 불필요 |
| charts/thanos-stack/presets/ | tokamak-thanos-stack | 수정 | values-*.yaml에 feeToken 섹션 추가 | 불필요 |
| internal/handler/fee_token_handler.go | trh-backend | **신규** | POST/GET/DELETE /api/chains/{id}/fee-tokens | 불필요 |
| src/pages/chains/[id]/fee-tokens.tsx | trh-platform-ui | **신규** | Fee Token 설정 패널 UI | 불필요 |
| Bundler config | — | 동일 | 변경 없음 (Paymaster가 처리) | 불필요 |

---

## 부록 A: paymasterAndData 인코딩

```
paymasterAndData 구조:

Offset  Length  Field
0       20      Paymaster 주소 (MultiTokenPaymaster)
20      20      Fee Token 주소 (선택된 ERC-20)
40      6       validUntil (서명 유효기간)
46      6       validAfter (서명 시작시간)
52      65      Paymaster 서명 (ECDSA)

총 117 bytes

UserOp에 포함되어 Bundler를 통해 EntryPoint에 전달됨.
EntryPoint가 _validatePaymasterUserOp 호출 시 이 데이터를 전달.
MultiTokenPaymaster가 offset 20~40에서 토큰 주소를 추출.
```

## 부록 B: 경쟁사 비교

| 플랫폼 | Multi-Token Fee | 방식 |
|--------|----------------|------|
| Thanos (본 설계) | TON, ETH, USDC, USDT | Paymaster 레벨. 프로토콜 변경 없음. |
| zkSync Era | ETH + ERC-20 (USDC 등) | 프로토콜 레벨 paymaster. 네이티브 AA. |
| StarkNet | ETH + STRK | 프로토콜 레벨. 두 네이티브 토큰 지원. |
| Arbitrum | ETH만 | Multi-token 미지원. |
| Base/OP Mainnet | ETH만 | Multi-token 미지원. Paymaster로 가능하지만 표준화 안 됨. |
| Caldera/Conduit (RaaS) | 체인별 설정 | 커스텀 fee token 가능하지만 Paymaster 인프라 미제공. |

Thanos의 차별화: **Preset 선택 = fee token 자동 설정.** `trh-sdk deploy --preset defi` 한 줄로 USDC/USDT가 가스비 수단에 자동 등록되고, 오라클 배포, 초기 가격 설정, CronJob 등록까지 완료된다. 운영자 입력 0개. Caldera/Conduit에서는 운영자가 직접 Paymaster를 개발하고 오라클을 배포해야 하며, zkSync Era는 프로토콜 레벨 변경이 필요하다. Thanos는 프로토콜을 건드리지 않으면서도 배포 시점에 Multi-Token Fee가 즉시 작동하는 유일한 OP Stack 기반 L2이다.

---

*End of Document*
