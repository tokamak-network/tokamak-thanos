# [가이드] GameType 2 DevNet 배포 및 E2E 테스트

**작성일**: 2025년 10월 23일
**목적**: GameType 2 (Asterisc) DevNet 배포 및 E2E 테스트 실행 가이드

---

## 📋 현재 상태 체크리스트

### ✅ 완료된 작업
- [x] Off-chain 구현 완료 (Phase 1-4)
- [x] 빌드 성공
- [x] Unit 테스트 100% 통과
- [x] Proposal 타입 통합
- [x] StateConverter 구현
- [x] **DevNet 배포 완료** (2025-10-23)
- [x] **GameType 2 온체인 게임 생성 성공** (2025-10-23)
- [x] **Challenger Bisection/Step 테스트 100% 통과** (2025-10-23)
- [x] **온체인 컨트랙트 배포 및 검증** (2025-10-23)

### ⏳ 진행 중
- [ ] op-challenger 실시간 게임 참여 테스트
- [ ] Asterisc E2E 자동화 테스트 작성

---

## 🚀 단계별 가이드

### Step 1: 사전 준비

#### 1.1 필수 도구 확인
```bash
# Python 3.9 이상
python3 --version

# Node.js (v16 이상)
node --version

# pnpm
pnpm --version

# Docker
docker --version

# foundry (forge, cast)
forge --version
cast --version
```

#### 1.2 의존성 빌드
```bash
cd /Users/zena/tokamak-projects/tokamak-thanos

# Go 바이너리 빌드
make build-go

# TypeScript 빌드 (contracts)
make build-ts
```

---

### Step 2: 온체인 컨트랙트 상태 확인

#### 2.1 RISCV.sol 컴파일 상태 확인
```bash
cd packages/tokamak/contracts-bedrock

# RISCV.sol 확인
ls -la src/dispute/RISCV.sol

# 컴파일 테스트
forge build --sizes

# RISCV 컨트랙트만 빌드
forge build --contracts src/dispute/RISCV.sol
```

**예상 결과**:
```
✅ RISCV.sol 컴파일 성공
✅ PreimageOracle.sol 컴파일 성공
✅ FaultDisputeGame.sol 컴파일 성공
```

#### 2.2 GameType 2 설정 확인
```bash
# deploy-config 확인
cat deploy-config/devnetL1-template.json | grep -A 5 "gameType"
```

**확인 사항**:
- `faultGameWithdrawalDelay`: 설정됨
- GameType 2 관련 설정 존재 여부

---

### Step 3: DevNet 실행

#### 3.1 DevNet 환경 정리 (선택)
```bash
# 기존 DevNet 정리
make devnet-clean

# .devnet 디렉토리 제거
rm -rf .devnet

# Docker 컨테이너 정리
docker-compose down -v
```

#### 3.2 Allocation 파일 생성
```bash
# DevNet allocation 생성
make devnet-allocs

# 생성된 파일 확인
ls -la .devnet/
```

**예상 출력**:
```
.devnet/
├── allocs-l1.json
├── allocs-l2.json
├── genesis-l1.json
└── genesis-l2.json
```

#### 3.3 DevNet 시작
```bash
# DevNet 실행
make devnet-up
```

**예상 로그**:
```
Starting L1 (anvil)...
Deploying L1 contracts...
✅ DisputeGameFactory deployed
✅ RISCV deployed
✅ PreimageOracle deployed
Starting L2...
Starting op-node...
Starting op-batcher...
Starting op-proposer...
```

#### 3.4 DevNet 상태 확인
```bash
# 다른 터미널에서
# L1 RPC 확인
cast block-number --rpc-url http://localhost:8545

# L2 RPC 확인
cast block-number --rpc-url http://localhost:9545

# 로그 확인
make devnet-logs
```

---

### Step 4: 온체인 컨트랙트 배포 확인

> **✅ 2025-10-23 배포 완료**: 모든 컨트랙트가 성공적으로 배포되었습니다.
>
> **배포된 주소**:
> ```
> DisputeGameFactory: 0x11c81c1A7979cdd309096D1ea53F887EA9f8D14d
> RISCV VM:          0xEad59ca6b38c83EcD7735006Db68a29c5e8A96A2
> GameType 2 Impl:   0xCe8950f4c5597E721b82f63185784a0971E69662
> ```

#### 4.1 DisputeGameFactory 확인
```bash
# DisputeGameFactory 주소
DGF="0x11c81c1A7979cdd309096D1ea53F887EA9f8D14d"

# 게임 개수 조회
cast call $DGF "gameCount()(uint256)" --rpc-url http://localhost:8545

# GameType 2가 등록되었는지 확인
cast call $DGF "gameImpls(uint32)(address)" 2 --rpc-url http://localhost:8545
# 출력: 0xCe8950f4c5597E721b82f63185784a0971E69662 (GameType 2 구현체)
```

#### 4.2 RISCV 컨트랙트 확인
```bash
# RISCV VM 주소
RISCV="0xEad59ca6b38c83EcD7735006Db68a29c5e8A96A2"

# 코드 배포 확인
cast code $RISCV --rpc-url http://localhost:8545 | wc -c
# 출력: 17806 bytes (RISCV VM 코드)
```

#### 4.3 GameType 2 구현체 확인
```bash
# GameType 2 구현체 주소
GT2_IMPL="0xCe8950f4c5597E721b82f63185784a0971E69662"

# VM 주소 확인
cast call $GT2_IMPL "vm()(address)" --rpc-url http://localhost:8545
# 출력: 0xEad59ca6b38c83EcD7735006Db68a29c5e8A96A2 (RISCV)

# GameType 확인
cast call $GT2_IMPL "gameType()(uint32)" --rpc-url http://localhost:8545
# 출력: 2

# MaxGameDepth 확인
cast call $GT2_IMPL "maxGameDepth()(uint256)" --rpc-url http://localhost:8545
# 출력: 50

# SplitDepth 확인
cast call $GT2_IMPL "splitDepth()(uint256)" --rpc-url http://localhost:8545
# 출력: 14
```

---

### Step 5: E2E 테스트 실행

#### 5.1 Alphabet 테스트 (간단한 테스트)
```bash
cd op-e2e

# Alphabet 게임 테스트 (GameType 0)
go test -v ./faultproofs -run TestOutputAlphabetGame_ChallengerWins
```

**예상 결과**:
```
=== RUN   TestOutputAlphabetGame_ChallengerWins
--- PASS: TestOutputAlphabetGame_ChallengerWins (30.25s)
PASS
```

#### 5.2 Cannon 테스트 (GameType 1)
```bash
# Cannon 게임 테스트
OP_E2E_CANNON_ENABLED=true go test -v ./faultproofs -run TestOutputCannonGame
```

**예상 결과**:
```
=== RUN   TestOutputCannonGame
--- PASS: TestOutputCannonGame (120.45s)
PASS
```

#### 5.3 전체 Fault Proof 테스트
```bash
# 모든 fault proof 테스트 실행
make test-fault-proofs
```

---

### Step 6: Asterisc (GameType 2) 통합 테스트

> **✅ 2025-10-23 완료**: GameType 2 통합 테스트가 성공적으로 완료되었습니다!
>
> **테스트 결과 요약**:
> - ✅ DevNet 배포 및 RISCV 컨트랙트 배포 성공
> - ✅ GameType 2 게임 온체인 생성 성공 (`0x328cfa286df1b8f099ac20a7921028b7e8ec5e0d`)
> - ✅ Bisection 알고리즘 테스트 100% 통과
> - ✅ Step 실행 테스트 100% 통과
> - ✅ 유닛 테스트 77개 모두 통과
>
> 상세한 테스트 결과는 [Challenger 테스트 리포트](./gametype2-challenger-test-report-ko.md)를 참조하세요.

#### 6.1 수동 통합 테스트

Asterisc E2E 테스트 파일이 아직 없으므로, 수동으로 통합을 확인합니다.

**테스트 스크립트 작성**: `test-asterisc-integration.sh`
```bash
#!/bin/bash
set -e

echo "🔍 GameType 2 (Asterisc) 통합 테스트"

# 1. op-challenger 빌드 확인
echo "1. Building op-challenger..."
cd /Users/zena/tokamak-projects/tokamak-thanos
go build ./op-challenger/...
echo "✅ op-challenger 빌드 성공"

# 2. Asterisc 패키지 테스트
echo "2. Running Asterisc unit tests..."
go test ./op-challenger/game/fault/trace/asterisc/... -v
echo "✅ Asterisc unit 테스트 통과"

# 3. Outputs 패키지 테스트
echo "3. Running Outputs tests..."
go test ./op-challenger/game/fault/trace/outputs/... -v
echo "✅ Outputs 테스트 통과"

# 4. 전체 op-challenger 빌드
echo "4. Building full op-challenger binary..."
go build -o op-challenger ./op-challenger/cmd
echo "✅ op-challenger 바이너리 생성 성공"

# 5. 버전 확인
echo "5. Checking op-challenger version..."
./op-challenger --version
echo "✅ op-challenger 실행 가능"

echo ""
echo "🎉 GameType 2 통합 테스트 완료!"
```

**실행**:
```bash
chmod +x test-asterisc-integration.sh
./test-asterisc-integration.sh
```

#### 6.2 DevNet에서 수동 테스트

**1. Dispute Game 생성 (수동)**

DevNet이 실행 중인 상태에서:

```bash
# L2 output proposal 확인
cast call <L2_OUTPUT_ORACLE_ADDRESS> "latestOutputIndex()(uint256)" --rpc-url http://localhost:8545

# 잘못된 output으로 dispute game 생성 (테스트용)
# DisputeGameFactory.create() 호출
cast send <DISPUTE_GAME_FACTORY_ADDRESS> \
  "create(uint32,bytes32,bytes)(address)" \
  2 \  # GameType 2 (Asterisc)
  0x0000000000000000000000000000000000000000000000000000000000000000 \  # 잘못된 root claim
  0x \  # 빈 extra data
  --rpc-url http://localhost:8545 \
  --private-key <PRIVATE_KEY>
```

**2. op-challenger 실행**

```bash
# op-challenger 설정 파일 작성
cat > challenger-config.toml <<EOF
[challenger]
l1-eth-rpc = "http://localhost:8545"
l1-beacon = "http://localhost:9596"
rollup-rpc = "http://localhost:9545"
game-factory-address = "<DISPUTE_GAME_FACTORY_ADDRESS>"
trace-type = "asterisc"
asterisc-bin = "./op-program/bin/op-program-client-riscv64"
asterisc-server = "./op-program/bin/op-program"
asterisc-prestate = "./op-program/bin/prestate.json"
datadir = "./challenger-data"
EOF

# op-challenger 실행
./op-challenger --config challenger-config.toml
```

**3. 로그 확인**

op-challenger 로그에서 다음을 확인:
```
✅ "Loaded GameType 2 configuration"
✅ "Asterisc trace provider initialized"
✅ "Monitoring dispute games"
```

---

### Step 7: 문제 해결

#### 7.1 DevNet이 시작되지 않는 경우

**증상**:
```
Error: Failed to start L1 node
```

**해결**:
```bash
# Docker 재시작
docker-compose down
docker system prune -f

# 포트 확인
lsof -i :8545
lsof -i :9545

# DevNet 정리 후 재시작
make devnet-clean
make devnet-up
```

#### 7.2 RISCV 컨트랙트가 배포되지 않은 경우

**증상**:
```
Error: RISCV contract not found
```

**해결**:
```bash
# deployment 스크립트 확인
cd packages/tokamak/contracts-bedrock
cat script/Deploy.s.sol | grep RISCV

# 수동 배포 (필요시)
forge script script/Deploy.s.sol:Deploy \
  --rpc-url http://localhost:8545 \
  --broadcast \
  --private-key <DEPLOYER_PRIVATE_KEY>
```

#### 7.3 GameType 2가 등록되지 않은 경우

**증상**:
```
Error: GameType 2 not registered in DisputeGameFactory
```

**해결**:
```bash
# DisputeGameFactory에 GameType 2 등록
cast send <DISPUTE_GAME_FACTORY_ADDRESS> \
  "setImplementation(uint32,address)" \
  2 \  # GameType 2
  <RISCV_GAME_ADDRESS> \
  --rpc-url http://localhost:8545 \
  --private-key <ADMIN_PRIVATE_KEY>

# 등록 확인
cast call <DISPUTE_GAME_FACTORY_ADDRESS> \
  "gameImpls(uint32)(address)" \
  2 \
  --rpc-url http://localhost:8545
```

---

## 📊 테스트 체크리스트

### DevNet 배포
- [ ] DevNet 시작 성공 (`make devnet-up`)
- [ ] L1 노드 실행 중 (포트 8545)
- [ ] L2 노드 실행 중 (포트 9545)
- [ ] 컨트랙트 배포 완료

### 온체인 컨트랙트
- [ ] RISCV.sol 배포 확인
- [ ] PreimageOracle 배포 확인
- [ ] DisputeGameFactory에 GameType 2 등록
- [ ] FaultDisputeGame 배포 확인

### Off-chain 테스트
- [ ] op-challenger 빌드 성공
- [ ] Asterisc unit 테스트 통과
- [ ] Outputs 테스트 통과
- [ ] op-challenger 바이너리 실행 가능

### 통합 테스트
- [ ] Alphabet 테스트 통과 (GameType 0)
- [ ] Cannon 테스트 통과 (GameType 1) - 선택
- [ ] op-challenger가 GameType 2 인식
- [ ] 수동 dispute game 생성 가능
- [ ] op-challenger가 dispute 처리

---

## 🎯 다음 단계

### 1. Asterisc E2E 테스트 작성 (선택)

Cannon 테스트를 참고하여 Asterisc E2E 테스트 작성:

**op-e2e/faultproofs/output_asterisc_test.go** (신규)
```go
package faultproofs

import (
	"context"
	"testing"

	op_e2e "github.com/tokamak-network/tokamak-thanos/op-e2e"
)

func TestOutputAsteriscGame(t *testing.T) {
	op_e2e.InitParallel(t)
	ctx := context.Background()
	sys, l1Client := StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)

	// Asterisc dispute game 테스트 로직
	// ...
}
```

### 2. 성능 벤치마크

```bash
# Asterisc 실행 시간 측정
time ./op-program/bin/op-program-client-riscv64 run ...

# Cannon과 비교
time ./cannon/bin/cannon run ...
```

### 3. Testnet 배포 준비

DevNet 테스트 완료 후:
1. Sepolia testnet에 배포
2. 실제 dispute game 시뮬레이션
3. 모니터링 설정

---

## 📝 로그 및 디버깅

### DevNet 로그 확인
```bash
# 전체 로그
make devnet-logs

# 특정 컴포넌트 로그
docker logs op-node
docker logs op-batcher
docker logs op-proposer
```

### op-challenger 로그 레벨
```bash
# 디버그 모드로 실행
./op-challenger --log.level debug --config challenger-config.toml
```

### 주요 로그 패턴

**정상 작동**:
```
INFO  Loaded GameType 2 configuration
INFO  Asterisc trace provider initialized
INFO  Monitoring 0 active games
```

**문제 발생**:
```
ERROR Failed to load prestate file
ERROR RISCV contract not found
ERROR Invalid GameType configuration
```

---

## 🔗 참고 문서

- [통합 계획 문서](./gametype2-integration-plan-ko.md)
- [구현 작업 일지](./WORK-gametype2-implementation-2025-10-23-ko.md)
- [**Challenger 테스트 리포트**](./gametype2-challenger-test-report-ko.md) ⭐ NEW
- [Optimism Fault Proof 문서](https://docs.optimism.io/stack/fault-proofs/overview)

---

## ⚠️ 주의사항

1. **DevNet은 로컬 테스트용**
   - 실제 자산 사용 금지
   - 테스트 후 정리 (`make devnet-clean`)

2. **Asterisc 바이너리 준비**
   - `op-program-client-riscv64` 빌드 필요
   - prestate.json 파일 필요

3. **포트 충돌 주의**
   - L1: 8545
   - L2: 9545
   - Beacon: 9596

4. **리소스 사용량**
   - DevNet 실행 시 Docker 메모리 4GB 이상 권장
   - Asterisc 실행 시 CPU 사용량 높음

---

**작성 완료일**: 2025년 10월 23일
**다음 업데이트**: DevNet 테스트 완료 후
