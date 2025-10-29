# VM 이미지 빌드 및 공유 가이드

> 로컬에서 빌드한 VM 이미지를 공유하는 방법

## 📚 목차

1. [개요](#개요)
2. [현재 구조 vs 이미지 기반 구조](#현재-구조-vs-이미지-기반-구조)
3. [이미지 빌드 방법](#이미지-빌드-방법)
4. [이미지 공유 방법](#이미지-공유-방법)
5. [배포 스크립트 수정](#배포-스크립트-수정)
6. [실제 사용 시나리오](#실제-사용-시나리오)

---

## 개요

매번 VM을 빌드하는 대신, 사전에 빌드된 이미지를 사용하여 배포 시간을 크게 단축할 수 있습니다.

### 주요 이점
- ✅ **배포 시간 단축**: 30분 → 2분
- ✅ **일관된 바이너리 보장**
- ✅ **버전 관리 용이**
- ✅ **롤백 가능**
- ✅ **팀 간 공유 간편**

---

## 현재 구조 vs 이미지 기반 구조

### 현재 구조 (매번 빌드)
```
1. 배포 시마다 VM 빌드 (20-30분)
2. 로컬 바이너리를 볼륨 마운트
3. 매번 동일한 빌드 반복
```

### 제안 구조 (이미지 사용)
```
1. VM 이미지 사전 빌드 (1회)
2. Registry에 푸시
3. 배포 시 이미지 pull (1-2분)
```

---

## 이미지 빌드 방법

### 1. 개별 VM 이미지 빌드

```bash
# Cannon VM 빌드
cd /Users/zena/tokamak-projects/tokamak-thanos
docker build -f ops/docker/op-stack-go/Dockerfile \
  --target cannon-builder \
  --platform linux/amd64 \
  -t tokamak/cannon:latest .

# Asterisc VM 빌드
docker build -f ops/docker/op-stack-go/Dockerfile \
  --target asterisc-builder \
  --platform linux/amd64 \
  -t tokamak/asterisc:latest .

# Challenger 전체 빌드 (모든 VM 포함)
docker build -f ops/docker/op-stack-go/Dockerfile \
  --target op-challenger-target \
  --platform linux/amd64 \
  -t tokamak/challenger:latest .
```

### 2. 통합 Dockerfile 작성 (선택사항)

`Dockerfile.vms` 파일 생성:

```dockerfile
# Cannon VM 이미지
FROM --platform=linux/amd64 golang:1.21.3-alpine3.18 as cannon-builder
WORKDIR /app
COPY . .
RUN cd cannon && make cannon

FROM scratch as cannon-image
COPY --from=cannon-builder /app/cannon/bin/cannon /cannon
COPY --from=cannon-builder /app/op-program/bin/prestate*.json /

# Asterisc VM 이미지
FROM --platform=linux/amd64 golang:1.21.3-alpine3.18 as asterisc-builder
WORKDIR /app
COPY . .
RUN cd asterisc && make asterisc

FROM scratch as asterisc-image
COPY --from=asterisc-builder /app/asterisc/bin/* /

# Kona 이미지
FROM --platform=linux/amd64 rust:1.82-bookworm as kona-builder
WORKDIR /app
COPY . .
RUN cargo build --release -p kona-client

FROM scratch as kona-image
COPY --from=kona-builder /app/target/release/kona-client /kona-client
```

---

## 이미지 공유 방법

### 방법 1: Docker Hub (가장 쉬움)

```bash
# 1. Docker Hub 로그인
docker login
# Username: your-dockerhub-username
# Password: your-password

# 2. 이미지 태그
docker tag tokamak/challenger:latest your-username/tokamak-challenger:latest

# 3. 이미지 푸시
docker push your-username/tokamak-challenger:latest

# 4. 다른 사람이 사용
docker pull your-username/tokamak-challenger:latest
```

### 방법 2: GitHub Container Registry (추천)

```bash
# 1. GitHub Personal Access Token으로 로그인
export CR_PAT=YOUR_TOKEN
echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin

# 2. 이미지 태그
docker tag tokamak/challenger:latest ghcr.io/your-github-username/tokamak-challenger:latest

# 3. 이미지 푸시
docker push ghcr.io/your-github-username/tokamak-challenger:latest

# 4. 다른 사람이 사용
docker pull ghcr.io/your-github-username/tokamak-challenger:latest
```

### 방법 3: 이미지 파일로 공유 (오프라인)

```bash
# 1. 이미지를 tar 파일로 저장
docker save -o tokamak-vms.tar \
  tokamak/cannon:latest \
  tokamak/asterisc:latest \
  tokamak/challenger:latest

# 2. 파일 공유 (Google Drive, S3 등)
# 파일 크기: 약 200-500MB

# 3. 다른 사람이 로드
docker load -i tokamak-vms.tar
```

---

## 배포 스크립트 수정

### 1. docker-compose.yml 수정

```yaml
# 기존 (로컬 빌드)
op-challenger:
  build:
    context: ${PROJECT_ROOT}
    dockerfile: ops/docker/op-stack-go/Dockerfile
    target: op-challenger-target
  image: tokamaknetwork/thanos-op-challenger:latest

# 변경 후 (이미지 사용)
op-challenger:
  # build 섹션 제거
  image: ghcr.io/your-username/tokamak-challenger:latest
```

### 2. vm-build.sh 수정 (선택사항)

```bash
# vm-build.sh에 추가
build_vms_from_images() {
    local image_tag="${VM_IMAGE_TAG:-latest}"
    local registry="${VM_REGISTRY:-ghcr.io/your-username}"

    log_info "Pulling pre-built VM images..."

    # Challenger 이미지에서 바이너리 추출
    docker create --name temp-challenger ${registry}/tokamak-challenger:${image_tag}

    # Cannon 바이너리 복사
    docker cp temp-challenger:/usr/local/bin/cannon ${CANNON_BIN}/cannon
    docker cp temp-challenger:/usr/local/bin/op-program ${OP_PROGRAM_BIN}/op-program

    # Asterisc 바이너리 복사 (마운트로 제공되므로 필요시)
    # docker cp temp-challenger:/usr/local/bin/asterisc ${ASTERISC_BIN}/asterisc

    docker rm temp-challenger

    log_success "VM binaries extracted from images"
}

# 환경변수로 제어
if [ "${USE_PREBUILT_VMS:-false}" = "true" ]; then
    build_vms_from_images
else
    # 기존 로컬 빌드 로직
    build_vms "$build_cannon" "$build_asterisc"
fi
```

---

## 실제 사용 시나리오

### 시나리오 1: 이미지 제공자 (당신)

```bash
# 1. 모든 변경사항 커밋
git add .
git commit -m "Update for GameType 3"

# 2. VM 이미지 빌드
docker build -f ops/docker/op-stack-go/Dockerfile \
  --target op-challenger-target \
  --platform linux/amd64 \
  -t ghcr.io/tokamak-network/challenger:v1.0.0 .

# 3. 이미지 푸시
docker push ghcr.io/tokamak-network/challenger:v1.0.0

# 4. 팀에 공유
echo "이미지 준비됨: ghcr.io/tokamak-network/challenger:v1.0.0"
```

### 시나리오 2: 이미지 사용자 (팀원)

```bash
# 1. docker-compose.yml 수정
# image: ghcr.io/tokamak-network/challenger:v1.0.0

# 2. 배포 실행
./op-challenger/scripts/cleanup.sh
./op-challenger/scripts/deploy-modular.sh --dg-type 3

# VM 빌드 없이 바로 배포 완료!
```

### 시나리오 3: 특정 버전 사용

```bash
# 환경변수로 버전 지정
export VM_IMAGE_TAG=v1.0.0
export USE_PREBUILT_VMS=true

# 배포
./op-challenger/scripts/deploy-modular.sh --dg-type 3
```

---

## CI/CD 자동화 (선택사항)

`.github/workflows/build-vms.yml`:

```yaml
name: Build and Push VM Images

on:
  push:
    branches: [feature/challenger-gametype3]
    paths:
      - 'cannon/**'
      - 'asterisc/**'
      - 'op-program/**'
      - 'op-challenger/**'

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push Challenger
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./ops/docker/op-stack-go/Dockerfile
          target: op-challenger-target
          platforms: linux/amd64
          push: true
          tags: |
            ghcr.io/${{ github.repository_owner }}/challenger:latest
            ghcr.io/${{ github.repository_owner }}/challenger:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

---

## 주의사항

### ⚠️ 중요 사항
- **플랫폼**: 반드시 `linux/amd64`로 빌드
- **이미지 크기**: 각 VM 약 50-200MB
- **프라이빗 이미지**: 사용자가 로그인 필요
- **버전 태그**: 명확한 버전 관리 필요
- **호환성**: 코드 변경 시 이미지도 업데이트 필요

### 🔒 보안 고려사항
- 프라이빗 레지스트리 사용 권장
- 민감한 정보는 이미지에 포함하지 않음
- 정기적인 보안 스캔 실행

---

## FAQ

### Q: 이미지 빌드가 실패합니다
A: `--platform linux/amd64` 옵션 확인, Docker 데몬 상태 확인

### Q: 이미지가 너무 큽니다
A: multi-stage 빌드 사용, 불필요한 파일 제거

### Q: 팀원이 이미지를 pull할 수 없습니다
A: 레지스트리 권한 확인, 로그인 상태 확인

### Q: 버전 관리는 어떻게 하나요?
A: Git 태그와 동일한 버전 사용 권장 (예: v1.0.0)

---

## 요약

VM 이미지를 사전 빌드하여 공유하면:
1. **배포 시간이 30분에서 2분으로 단축**
2. **팀 전체가 동일한 바이너리 사용**
3. **버전 관리와 롤백이 용이**

가장 간단한 방법은 Docker Hub 또는 GitHub Container Registry를 사용하는 것입니다.