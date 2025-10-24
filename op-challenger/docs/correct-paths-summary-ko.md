# Tokamak-Thanos 올바른 경로 정보

## 중요 정정

Tokamak-Thanos의 컨트랙트 위치는:

### ❌ 틀린 경로
```
/Users/zena/tokamak-projects/tokamak-thanos/packages/contracts-bedrock/
```

### ✅ 올바른 경로
```
/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/
```

**Tokamak 커스텀 버전**입니다!

## 현재 상태

### Optimism 경로
```
/Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/
├── src/vendor/asterisc/RISCV.sol          # ✅ 있음
├── interfaces/vendor/asterisc/IRISCV.sol  # ✅ 있음
└── scripts/deploy/DeployAsterisc.s.sol    # ✅ 있음
```

### Tokamak-Thanos 경로 (정정됨)
```
/Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/
├── src/vendor/
│   ├── AddressAliasHelper.sol             # ✅ 있음
│   ├── WNativeToken.sol                   # ✅ 있음
│   └── asterisc/                          # ❌ 없음 (생성 필요!)
├── interfaces/vendor/
│   └── asterisc/                          # ❌ 없음 (생성 필요!)
└── scripts/deploy/
    └── DeployAsterisc.s.sol               # ❌ 없음 (복사 필요!)
```

## 올바른 명령어

### 디렉토리 생성
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

mkdir -p src/vendor/asterisc
mkdir -p interfaces/vendor/asterisc
mkdir -p test/opcm
```

### 파일 복사
```bash
# RISCV.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/vendor/asterisc/RISCV.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/src/vendor/asterisc/

# IRISCV.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/vendor/asterisc/IRISCV.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/interfaces/vendor/asterisc/

# DeployAsterisc.s.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAsterisc.s.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/scripts/deploy/

# DeployAsterisc.t.sol 복사
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/test/opcm/DeployAsterisc.t.sol \
   /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock/test/opcm/
```

### 컴파일
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

forge build
```

### 테스트
```bash
forge test --match-contract DeployAsterisc_Test -vvv
```

## 업데이트 필요한 문서

다음 문서들의 경로를 정정해야 합니다:

1. **onchain-contracts-integration-ko.md**
   - 모든 Tokamak 경로: `packages/contracts-bedrock/` → `packages/tokamak/contracts-bedrock/`

2. **gametype2-integration-plan-ko.md**
   - Phase 0 명령어 경로 수정

3. **contract-files-checklist-ko.md**
   - 전체 체크리스트 경로 수정

## 빠른 시작 (정정된 버전)

```bash
# 1. Tokamak contracts-bedrock로 이동
cd /Users/zena/tokamak-projects/tokamak-thanos/packages/tokamak/contracts-bedrock

# 2. 디렉토리 생성
mkdir -p src/vendor/asterisc interfaces/vendor/asterisc test/opcm

# 3. 파일 복사 (한 번에)
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/src/vendor/asterisc/RISCV.sol src/vendor/asterisc/
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/interfaces/vendor/asterisc/IRISCV.sol interfaces/vendor/asterisc/
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/scripts/deploy/DeployAsterisc.s.sol scripts/deploy/
cp /Users/zena/tokamak-projects/optimism/packages/contracts-bedrock/test/opcm/DeployAsterisc.t.sol test/opcm/

# 4. 확인
ls -lh src/vendor/asterisc/RISCV.sol
ls -lh interfaces/vendor/asterisc/IRISCV.sol

# 5. 컴파일
forge build

# 6. 테스트
forge test --match-contract DeployAsterisc_Test -vvv
```

이제 올바른 경로로 작업하시면 됩니다!
