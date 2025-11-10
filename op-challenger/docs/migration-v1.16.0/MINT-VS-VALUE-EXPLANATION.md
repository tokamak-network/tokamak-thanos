# _mint vs _value 상세 설명

## 🎯 핵심 요약 (30초 버전)

```
_mint  = L1에서 lock → L2 계정 잔액 증가량
_value = L2 트랜잭션 실행 시 msg.value
```

**관계**: `_value ≤ _mint` (항상!)

---

## 📚 자세한 설명

### 1. _mint (Mint Amount)
**"L1에서 lock하고 L2에서 mint할 토큰 총량"**

```solidity
// L1 OptimismPortal2에서:
if (_mint > 0) {
    IERC20(nativeToken).transferFrom(user, portal, _mint);
    // → Portal에 _mint만큼 lock
}

// L2에서:
recipient.balance += _mint;
// → L2 계정 잔액 증가
```

### 2. _value (Call Value)
**"L2 트랜잭션 실행 시 msg.value로 사용할 량"**

```solidity
// L2에서:
_to.call{value: _value}(_data);
// → 컨트랙트 호출 시 msg.value
```

---

## 💡 실제 시나리오

### 시나리오 1: 단순 토큰 전송
**"Alice가 Bob에게 100 TON을 L2로 전송"**

```solidity
portal.depositTransaction(
    _to: Bob,
    _mint: 100 TON,    // Bob L2 계정에 100 TON 생성
    _value: 0,         // 컨트랙트 호출 없음
    _gasLimit: 21000,
    _isCreation: false,
    _data: ""
);
```

**결과**:
- L1: Portal에 100 TON lock
- L2: Bob 계정에 100 TON mint
- 컨트랙트 호출 없음

---

### 시나리오 2: 토큰 전송 + 컨트랙트 호출
**"Alice가 100 TON을 L2로 보내면서, 30 TON으로 DEX swap 실행"**

```solidity
portal.depositTransaction(
    _to: DEXContract,
    _mint: 100 TON,    // DEX 계정에 100 TON 생성
    _value: 30 TON,    // swap 함수 호출 시 msg.value
    _gasLimit: 500000,
    _isCreation: false,
    _data: abi.encodeWithSignature("swap()")
);
```

**L2 실행 순서**:
1. `DEXContract.balance += 100 TON`
2. `DEXContract.swap{value: 30 TON}()`
3. swap 함수 내부에서 `msg.value == 30 TON` ✅
4. 나머지 70 TON은 DEXContract 잔액으로 유지

---

### 시나리오 3: 모든 토큰 즉시 사용
**"Alice가 NFT를 100 TON에 구매"**

```solidity
portal.depositTransaction(
    _to: NFTMarket,
    _mint: 100 TON,    // NFTMarket에 100 TON 생성
    _value: 100 TON,   // 전부 buy 함수에 사용
    _gasLimit: 300000,
    _isCreation: false,
    _data: abi.encodeWithSignature("buy(uint256)", tokenId)
);
```

**L2 실행**:
1. `NFTMarket.balance += 100 TON`
2. `NFTMarket.buy{value: 100 TON}(tokenId)`
3. buy 함수에서 100 TON 전부 사용
4. NFTMarket 잔액 = 0 TON

---

## ⚠️ 잘못된 사용 예시

### ❌ Case 1: _value > _mint
```solidity
portal.depositTransaction(
    _to: SomeContract,
    _mint: 50 TON,     // 50 TON만 생성
    _value: 100 TON,   // 100 TON 사용하려고 시도
    ...
);
// → L2에서 revert! (가진 것보다 많이 사용 불가)
```

### ❌ Case 2: 잔액 없이 value 지정
```solidity
portal.depositTransaction(
    _to: SomeContract,
    _mint: 0,          // 토큰 생성 안 함
    _value: 10 TON,    // 10 TON 사용하려고 시도
    ...
);
// → L2에서 revert! (잔액 0인데 value 지정)
```

---

## 🔄 Ethereum과의 비교

### Ethereum (ETH)
```solidity
// 하나의 값만 사용
to.call{value: msg.value}(data);

// msg.value가 곧 계정 잔액 증가량
```

### Tokamak (Native Token)
```solidity
// 두 개의 값 사용
depositTransaction(
    to,
    _mint,   // 계정 잔액 증가량
    _value,  // 트랜잭션 msg.value
    ...
);

// _mint ≥ _value 관계
```

---

## 🎓 왜 이런 구조인가?

### Ethereum의 제약:
- ETH는 네이티브 자산
- `msg.value`와 잔액 증가량이 항상 동일
- **분리 불가능**

### Tokamak의 유연성:
- Native Token은 ERC20
- 잔액 증가량(_mint)과 사용량(_value)을 **분리 가능**
- 더 유연한 사용 패턴:
  - 토큰만 전송 (value=0)
  - 토큰 전송 + 일부 사용 (value < mint)
  - 모두 즉시 사용 (value = mint)

---

## 📊 opaque Data 구조

```solidity
// Tokamak
bytes memory opaqueData = abi.encodePacked(
    _mint,      // uint256 - L2 잔액 증가량
    _value,     // uint256 - L2 msg.value
    _gasLimit,  // uint64
    _isCreation,// bool
    _data       // bytes
);

// Optimism v1.16.0 (비교)
bytes memory opaqueData = abi.encodePacked(
    msg.value,  // uint256 - L1 ETH → L2 잔액 = msg.value
    _value,     // uint256 - L2 msg.value (동일)
    _gasLimit,
    _isCreation,
    _data
);
```

**차이점**: Tokamak은 `_mint`와 `_value`를 분리하여 더 유연함!

---

## ✅ 체크리스트

마이그레이션 시 확인사항:
- [ ] `_mint` 파라미터 유지 (6-parameter 함수)
- [ ] L1에서 `_mint`만큼 token lock
- [ ] opaqueData에 `_mint`, `_value` 순서로 인코딩
- [ ] L2에서 `_mint` → 잔액 증가, `_value` → msg.value 처리
- [ ] `_value ≤ _mint` 검증 로직

---

**작성일**: 2025-11-10
**참고**: `OPTIMISMPORTAL2-MIGRATION-STRATEGY.md`

