# State Dump 파일 생성 문제 해결 가이드

**작성일**: 2025-11-05
**문제**: `make devnet-allocs` 실행 시 state-dump 파일이 생성되지 않음
**원인**: `sortJsonByKeys` 함수에서 jq 명령 실패로 파일이 덮어쓰여짐
**중요성**: Challenger 테스트에 필수

## 문제 분석

### 근본 원인

`L2Genesis.s.sol`의 `sortJsonByKeys` 함수:

```solidity
// Line 556-562
function sortJsonByKeys(string memory _path) internal {
    string[] memory commands = new string[](3);
    commands[0] = "bash";
    commands[1] = "-c";
    commands[2] = string.concat("cat <<< $(jq -S '.' ", _path, ") > ", _path);
    Process.run(commands);
}
```

**문제점**:
1. `vm.dumpState(_path)`로 파일 생성 → 성공
2. `sortJsonByKeys(_path)`로 정렬 시도 → **jq 실패 시 파일 삭제**
3. 결과: 빈 파일 또는 파일 없음

## 해결 방안

### Option 1: sortJsonByKeys 함수 수정 (권장)

**파일**: `packages/tokamak/contracts-bedrock/scripts/L2Genesis.s.sol`

```solidity
function sortJsonByKeys(string memory _path) internal {
    string[] memory commands = new string[](3);
    commands[0] = "bash";
    commands[1] = "-c";
    // 안전한 버전: 먼저 임시 파일로 저장 후 성공 시에만 덮어쓰기
    commands[2] = string.concat(
        "jq -S '.' ", _path, " > ", _path, ".tmp && mv ", _path, ".tmp ", _path,
        " || echo 'Warning: Failed to sort JSON, keeping original file'"
    );
    Process.run(commands);
}
```

또는 더 간단하게:

```solidity
function sortJsonByKeys(string memory _path) internal {
    // 정렬을 건너뛰기 (정렬은 필수가 아님)
    console.log("Skipping JSON sorting to preserve file integrity");
}
```

### Option 2: 수동으로 state dump 생성

```bash
# 1. Forge 스크립트 직접 실행 (정렬 없이)
cd packages/tokamak/contracts-bedrock

# 2. L2Genesis.s.sol 임시 수정
# sortJsonByKeys 함수를 주석 처리

# 3. 스크립트 실행
forge script scripts/L2Genesis.s.sol:L2Genesis \
  --sig "runWithStateDump()" \
  --chain-id 31337 \
  --fork-url http://localhost:8545

# 4. 파일 확인
ls -la state-dump*.json
```

### Option 3: Python 스크립트로 우회

**파일**: `generate_state_dump.py`

```python
#!/usr/bin/env python3
import subprocess
import json
import os

def generate_state_dump():
    """Generate state dump without sorting"""

    # Run forge script without sorting
    cmd = [
        "forge", "script",
        "scripts/L2Genesis.s.sol:L2Genesis",
        "--sig", "runWithStateDump()",
        "--chain-id", "31337"
    ]

    result = subprocess.run(
        cmd,
        cwd="packages/tokamak/contracts-bedrock",
        capture_output=True,
        text=True
    )

    # Check for state dump files
    files = [
        "state-dump-901.json",
        "state-dump-901-delta.json",
        "state-dump-901-ecotone.json"
    ]

    for filename in files:
        filepath = f"packages/tokamak/contracts-bedrock/{filename}"
        if os.path.exists(filepath):
            print(f"✅ Generated: {filepath}")
            # Optional: sort JSON safely
            try:
                with open(filepath, 'r') as f:
                    data = json.load(f)
                with open(filepath, 'w') as f:
                    json.dump(data, f, indent=2, sort_keys=True)
                print(f"   Sorted: {filepath}")
            except Exception as e:
                print(f"   Warning: Could not sort {filepath}: {e}")

if __name__ == "__main__":
    generate_state_dump()
```

## 즉시 적용 가능한 해결책

### Quick Fix (임시 해결)

1. **sortJsonByKeys 함수 비활성화**:

```bash
# L2Genesis.s.sol 수정
sed -i.bak 's/sortJsonByKeys(_path);/\/\/ sortJsonByKeys(_path);/' \
  packages/tokamak/contracts-bedrock/scripts/L2Genesis.s.sol

# make devnet-allocs 재실행
make devnet-allocs

# 파일 확인
ls -la packages/tokamak/contracts-bedrock/state-dump*.json
```

2. **원복** (필요시):

```bash
mv packages/tokamak/contracts-bedrock/scripts/L2Genesis.s.sol.bak \
   packages/tokamak/contracts-bedrock/scripts/L2Genesis.s.sol
```

## 영구적 해결책

### L2Genesis.s.sol 수정

```diff
--- a/packages/tokamak/contracts-bedrock/scripts/L2Genesis.s.sol
+++ b/packages/tokamak/contracts-bedrock/scripts/L2Genesis.s.sol
@@ -549,7 +549,7 @@ contract L2Genesis is Deployer {

         console.log("Writing state dump to: %s", _path);
         vm.dumpState(_path);
-        sortJsonByKeys(_path);
+        // sortJsonByKeys(_path); // Disabled: causes file deletion on jq failure
     }

     /// @notice Sorts the allocs by address
```

또는 안전한 버전으로 교체:

```diff
@@ -556,7 +556,14 @@ contract L2Genesis is Deployer {
     function sortJsonByKeys(string memory _path) internal {
         string[] memory commands = new string[](3);
         commands[0] = "bash";
         commands[1] = "-c";
-        commands[2] = string.concat("cat <<< $(jq -S '.' ", _path, ") > ", _path);
+        // Safe version: only overwrite if jq succeeds
+        commands[2] = string.concat(
+            "if jq -S '.' ", _path, " > ", _path, ".tmp 2>/dev/null; then ",
+            "mv ", _path, ".tmp ", _path, "; ",
+            "else ",
+            "echo 'Warning: JSON sorting skipped'; ",
+            "fi"
+        );
         Process.run(commands);
     }
```

## 검증

수정 후 검증:

```bash
# 1. 기존 파일 삭제
rm -f packages/tokamak/contracts-bedrock/state-dump*.json

# 2. 재생성
make devnet-allocs

# 3. 파일 확인
ls -lh packages/tokamak/contracts-bedrock/state-dump*.json

# 4. 내용 확인
head -10 packages/tokamak/contracts-bedrock/state-dump-901.json
```

예상 결과:
```
-rw-r--r--  1 user  staff  1.2M Nov  5 11:00 state-dump-901-delta.json
-rw-r--r--  1 user  staff  1.2M Nov  5 11:00 state-dump-901-ecotone.json
-rw-r--r--  1 user  staff  1.2M Nov  5 11:00 state-dump-901.json
```

## Challenger 테스트를 위한 다음 단계

State dump 파일이 생성되면:

1. **원본 계획 (Option C) 구현**
   - State dump 로딩 방식 사용
   - 완전한 state 일치 보장

2. **E2E 테스트 실행**:
```bash
export USE_TOKAMAK_STATE_DUMP=true
go test -v ./op-e2e/faultproofs -timeout 30m
```

3. **Challenger 테스트**:
```bash
go test -v -run TestChallenger ./op-e2e/faultproofs
```

## 트러블슈팅

### jq가 설치되지 않은 경우

```bash
# macOS
brew install jq

# Linux
sudo apt-get install jq
```

### 파일 권한 문제

```bash
chmod +w packages/tokamak/contracts-bedrock/
```

### Forge 버전 문제

```bash
foundryup  # Foundry 업데이트
forge --version  # 버전 확인
```

## 결론

**즉시 적용 가능한 해결책**: `sortJsonByKeys` 주석 처리

**장기적 해결책**: 안전한 정렬 함수로 교체

이 수정으로 Challenger 테스트에 필요한 state dump 파일이 정상 생성됩니다.