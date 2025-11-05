# Traefik 네트워크 문제 자동 해결 가이드

## 문제 설명

Kurtosis devnet 배포 시 Traefik 역방향 프록시에서 다음과 같은 네트워크 문제가 발생할 수 있습니다:

```
Could not find network named "aec9d72af47f16cdd3285619cc874ec4f96a6cf7a1a4675c4d69325c7871eef7"
```

이는 Traefik이 하드코딩된 네트워크 ID를 참조하지만, 실제 Docker 네트워크 ID가 다를 때 발생합니다.

## 자동 해결 방법

### 1. 단독 실행

```bash
# 프로젝트 루트에서
just fix-traefik

# 또는 직접 스크립트 실행
cd kurtosis-devnet && ./fix-traefik.sh
```

### 2. 배포와 함께 실행

```bash
# Simple devnet with automatic fix
just simple-devnet-fixed

# 또는 kurtosis-devnet 디렉토리에서
just devnet-with-fix simple.yaml
```

### 3. AUTOFIX와 함께 사용

```bash
# 기존 방식
AUTOFIX=true just simple-devnet
# 배포 완료 후 (Traefik 에러 발생 시)
just fix-traefik
```

## 스크립트 동작 원리

`fix-traefik.sh` 스크립트는 다음을 수행합니다:

1. **컨테이너 감지**: 실행 중인 Traefik 컨테이너 자동 찾기
2. **설정 백업**: 현재 설정을 `/tmp/`에 백업
3. **동적 설정 적용**:
   - 하드코딩된 네트워크 ID 제거
   - `watch: true` - 실시간 Docker 변화 감지
   - `pollInterval: "5s"` - 5초마다 네트워크 상태 확인
4. **컨테이너 재시작**: 새 설정 적용
5. **상태 검증**: 성공 여부 확인

## 개선된 Traefik 설정

**변경 전 (문제 발생):**
```yaml
providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    network: "aec9d72af47f16cdd3285619cc874ec4f96a6cf7a1a4675c4d69325c7871eef7"
```

**변경 후 (문제 해결):**
```yaml
providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    useBindPortIP: false
    watch: true
    pollInterval: "5s"
    defaultRule: "Host(`{{ normalize .Name }}`)"
```

## 추가 정보

### 백업 위치
- 설정 백업: `/tmp/traefik-backup-YYYYMMDD-HHMMSS.yml`

### 수동 복구
문제 재발 시:
```bash
docker restart $(docker ps --filter "name=kurtosis-reverse-proxy" --format "{{.Names}}")
```

### 로그 확인
```bash
docker logs $(docker ps --filter "name=kurtosis-reverse-proxy" --format "{{.Names}}") --follow
```

## 영구적 해결책

현재 이 스크립트는 임시 해결책입니다. 영구적 해결을 위해서는:

1. **Kurtosis Engine 수정**: 상위 레포지토리의 Traefik 설정 템플릿 개선
2. **Optimism Package 수정**: `optimism-package`의 네트워크 설정 개선

하지만 개발 목적으로는 현재 자동화 스크립트로 충분합니다.

## 문제 해결

### 스크립트 실행 실패
```bash
# 권한 확인
chmod +x kurtosis-devnet/fix-traefik.sh

# Docker 데몬 확인
docker ps
```

### Traefik 컨테이너 없음
```bash
# Kurtosis 엔클레이브 확인
kurtosis enclave ls

# 엔클레이브 재시작
AUTOFIX=true just simple-devnet
```