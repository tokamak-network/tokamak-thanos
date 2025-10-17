# DA Server 보안 취약점 분석

## 목차
1. [개요](#개요)
2. [주요 보안 취약점](#주요-보안-취약점)
3. [공격 시나리오](#공격-시나리오)
4. [보안 강화 방안](#보안-강화-방안)
5. [보안 체크리스트](#보안-체크리스트)
6. [참고 자료](#참고-자료)

---

## 개요

### DA Server란?

DA (Data Availability) Server는 Plasma 모드에서 L2 트랜잭션 데이터를 오프체인에 저장하는 핵심 컴포넌트입니다. L1에는 작은 크기의 commitment(해시)만 올리고 실제 데이터는 DA Server에 저장하여 가스 비용을 절감합니다.

### 작동 방식

```
┌─────────────────────────────────────────┐
│  op-batcher (L2 → L1 제출)              │
├─────────────────────────────────────────┤
│  1. L2 트랜잭션 데이터 생성              │
│  2. DA Server에 데이터 저장 (PUT)       │
│  3. Commitment(해시) 받기                │
│  4. L1에 commitment만 제출              │
└─────────────────────────────────────────┘
           ↓ PUT /put/{hash}
┌─────────────────────────────────────────┐
│  DA Server (오프체인 저장소)             │
├─────────────────────────────────────────┤
│  - 전체 트랜잭션 데이터 저장             │
│  - File Storage 또는 S3                 │
│  - Commitment 생성 및 반환               │
└─────────────────────────────────────────┘
           ↑ GET /get/{hash}
┌─────────────────────────────────────────┐
│  op-node (L2 실행)                      │
├─────────────────────────────────────────┤
│  1. L1에서 commitment 읽기               │
│  2. DA Server에서 실제 데이터 조회       │
│  3. 데이터 검증 및 L2 블록 생성          │
└─────────────────────────────────────────┘
```

### 보안 중요성

DA Server는 **L2의 데이터 가용성을 보장**하는 핵심 컴포넌트이므로, 보안 취약점은 다음과 같은 심각한 결과를 초래할 수 있습니다:

- 💰 **자금 손실**: 데이터 손실로 인한 인출 불가
- 🔒 **서비스 중단**: DA Server 다운 시 L2 전체 마비
- 🕵️ **정보 유출**: 모든 L2 트랜잭션 데이터 노출
- ⚠️ **무결성 침해**: 데이터 변조로 잘못된 상태 확정

---

## 주요 보안 취약점

### 1. 인증/인가 부재 (심각도: Critical)

#### 문제점

DA Server는 **어떠한 인증 메커니즘도 없이** HTTP 엔드포인트를 노출합니다.

**코드 분석** (`op-plasma/daserver.go`):
```go
func (d *DAServer) HandleGet(w http.ResponseWriter, r *http.Request) {
    // ❌ 인증 체크 없음
    // ❌ API 키 검증 없음
    // ❌ IP 화이트리스트 없음

    route := path.Dir(r.URL.Path)
    if route != "/get" {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // ... 데이터 조회
}

func (d *DAServer) HandlePut(w http.ResponseWriter, r *http.Request) {
    // ❌ 누구나 데이터 저장 가능
    input, err := io.ReadAll(r.Body)
    // ... 데이터 저장
}
```

#### 공격 예시

```bash
# 공격자가 임의의 데이터 저장
curl -X POST http://da-server:3100/put/ \
  -d "malicious_data"

# 공격자가 모든 데이터 조회
curl http://da-server:3100/get/0x1234567890abcdef...

# 스토리지 폭탄 공격
for i in {1..100000}; do
  dd if=/dev/urandom bs=1M count=1 | \
    curl -X POST http://da-server:3100/put/ --data-binary @-
done
```

#### 위험

- ✗ **데이터 오염**: 악의적인 데이터 주입으로 시스템 혼란
- ✗ **정보 유출**: 모든 L2 트랜잭션 데이터 노출 (개인정보, 거래 내역 등)
- ✗ **스토리지 고갈**: 무제한 데이터 업로드로 디스크/S3 비용 폭증
- ✗ **서비스 거부**: 악의적 요청으로 정상 서비스 방해

---

### 2. TLS/암호화 미적용 (심각도: High)

#### 문제점

Docker Compose 설정에서 **TLS가 활성화되지 않음**:

```yaml
# ops-bedrock/docker-compose.yml
da-server:
  command: >
    da-server
      --addr=0.0.0.0  # ❌ HTTP만 사용
      --port=3100
      --log.level=debug
```

#### 중간자 공격 시나리오

```
┌────────────┐   HTTP (평문)   ┌───────────┐
│ op-batcher ├──────────────────►│ DA Server │
└────────────┘                  └───────────┘
      │
      └─► 🕵️ 네트워크 도청
          ├─ 트랜잭션 데이터 탈취
          ├─ Commitment 변조
          └─ Replay 공격
```

#### 공격 예시

```bash
# Wireshark로 패킷 캡처
tcpdump -i any -w capture.pcap port 3100

# 평문 트랜잭션 데이터 추출
strings capture.pcap | grep -A 10 "POST /put"

# Replay 공격
curl -X POST http://da-server:3100/put/0xABCD... \
  --data "@captured_transaction.bin"
```

#### 위험

- ✗ **데이터 도청**: 모든 L2 트랜잭션이 평문으로 네트워크 전송
- ✗ **MITM 공격**: 중간자가 데이터를 변조하여 잘못된 commitment 생성
- ✗ **Replay 공격**: 이전 트랜잭션 재전송으로 이중 지불 시도
- ✗ **네트워크 스니핑**: 내부 네트워크 침투 시 모든 데이터 유출

---

### 3. Rate Limiting 부재 (심각도: High)

#### 문제점

```go
// op-plasma/daserver.go
func (d *DAServer) HandleGet(w http.ResponseWriter, r *http.Request) {
    // ❌ Rate limit 없음
    // ❌ 요청 빈도 제한 없음
    // ❌ 동시 연결 수 제한 없음
    // ❌ 리소스 사용량 제어 없음
}
```

#### DDoS 공격 시나리오

```bash
# HTTP Flood 공격
ab -n 1000000 -c 1000 http://da-server:3100/get/0x1234...

# Slowloris 공격
slowhttptest -c 1000 -H -g -o slowloris \
  -i 10 -r 200 -t GET -u http://da-server:3100/get/0x1234...

# Storage Exhaustion
while true; do
  dd if=/dev/urandom bs=10M count=1 | \
    curl -X POST http://da-server:3100/put/ --data-binary @- &
done
```

#### 비용 폭증 시나리오 (S3 사용 시)

```
공격자가 1시간 동안 공격:
- 초당 1000개 PUT 요청 × 3600초 = 3,600,000 요청
- 요청당 10MB 데이터 = 36TB 데이터
- AWS S3 비용:
  └─ PUT 요청: $0.005/1000 = $18
  └─ 저장 비용: 36TB × $0.023/GB = $828/월
  └─ 데이터 전송: 36TB × $0.09/GB = $3,240
  └─ 총 비용: ~$4,086 (1시간 공격으로!)
```

#### 위험

- ✗ **서비스 거부**: 정상 요청 처리 불가
- ✗ **리소스 고갈**: CPU/메모리/네트워크 대역폭 소진
- ✗ **비용 폭증**: S3 사용 시 API 호출 및 저장 비용 급증
- ✗ **가용성 저하**: DA Server 응답 지연으로 L2 블록 생성 지연

---

### 4. 입력 검증 부재 (심각도: Medium)

#### 문제점

```go
// op-plasma/daserver.go:136
func (d *DAServer) HandlePut(w http.ResponseWriter, r *http.Request) {
    input, err := io.ReadAll(r.Body)
    if err != nil {
        // ❌ 최대 크기 제한 없음
        // ❌ 데이터 형식 검증 없음
        // ❌ Content-Type 확인 없음
    }

    // ❌ MaxInputSize 체크 없음 (Plasma는 MaxInputSize 있지만 서버에선 미검증)
}
```

#### 공격 예시

```bash
# 메모리 폭탄 (10GB 업로드)
dd if=/dev/zero bs=1G count=10 | \
  curl -X POST http://da-server:3100/put/ --data-binary @-

# 압축 폭탄 (Zip Bomb)
curl -X POST http://da-server:3100/put/ \
  --data-binary @42.zip  # 압축 시 42KB, 해제 시 4.5PB

# Null byte injection
curl -X POST http://da-server:3100/put/ \
  -d "valid_data\x00../../etc/passwd"
```

#### 위험

- ✗ **메모리 초과**: OOM 킬러 발동으로 서비스 다운
- ✗ **디스크 고갈**: 스토리지 부족으로 정상 데이터 저장 불가
- ✗ **잘못된 데이터**: op-node가 데이터 파싱 실패로 블록 생성 중단
- ✗ **경로 조작**: 파일 스토리지 사용 시 디렉토리 탐색 공격

---

### 5. Challenge Window 공격 (심각도: Critical)

#### Plasma Challenge 메커니즘

Plasma 모드는 **Challenge-Response 시스템**으로 데이터 가용성을 보장합니다:

```solidity
// DataAvailabilityChallenge.sol
contract DataAvailabilityChallenge {
    uint256 public challengeWindow;  // 챌린지 가능 기간
    uint256 public resolveWindow;    // 해결 필요 기간

    // 챌린지 단계:
    // 1. Active: 챌린지 진행 중
    // 2. Resolved: 데이터 제공으로 해결
    // 3. Expired: 해결 실패로 만료
}
```

#### 공격 시나리오

```
시간 T0: Sequencer가 잘못된 commitment 제출
         └─ L1에 commitment: 0xBADDATA...

시간 T1: Challenger가 챌린지 시작
         └─ DA Server에서 데이터 조회 시도

시간 T2: 공격자가 DA Server DDoS 공격 🔥
         └─ GET /get/0xBADDATA... → Timeout

시간 T3: Challenge Window 만료 ⏰
         └─ 데이터를 제공하지 못해 챌린지 실패
         └─ 잘못된 commitment가 확정됨 💥

결과: 무효한 트랜잭션이 L2에 포함
      └─ 사용자 자금 손실 발생!
```

#### 코드 분석

```go
// op-plasma/damgr.go:195-208
switch ch.challengeStatus {
case ChallengeActive:
    if d.isExpired(ch.expiresAt) {
        // ❌ Challenge 만료 → 데이터 영구 손실
        return nil, ErrExpiredChallenge
    } else if notFound {
        // DA Server에서 데이터를 찾을 수 없으면
        // 챌린지 해결 불가
        return nil, ErrPendingChallenge
    }
}
```

#### 위험

- ✗ **데이터 가용성 실패**: Liveness 보장 불가
- ✗ **잘못된 상태 확정**: 무효한 트랜잭션이 최종 확정
- ✗ **자금 손실**: 사용자 인출 불가능
- ✗ **시스템 신뢰 상실**: L2 보안 모델 붕괴

---

### 6. 단일 장애점 (SPOF) (심각도: Critical)

#### 문제점

```yaml
# docker-compose.yml
services:
  da-server:
    image: tokamaknetwork/thanos-da-server
    # ❌ 단일 인스턴스만 실행
    # ❌ 백업 서버 없음
    # ❌ 자동 장애 복구 없음
    # ❌ 헬스체크 없음
```

#### 장애 전파 시나리오

```
┌─────────────────────────────────────────┐
│ DA Server 다운 🔥                        │
│ (하드웨어 장애, 네트워크 분리, 등)        │
└──────────────┬──────────────────────────┘
               ↓
┌──────────────────────────────────────────┐
│ op-batcher                               │
│ ❌ PUT 요청 실패 → 배치 제출 중단         │
└──────────────┬───────────────────────────┘
               ↓
┌──────────────────────────────────────────┐
│ op-node                                  │
│ ❌ GET 요청 실패 → 데이터 조회 불가       │
│ ❌ L2 블록 생성 중단                      │
└──────────────┬───────────────────────────┘
               ↓
┌──────────────────────────────────────────┐
│ L2 전체 중단 💥                           │
│ - 트랜잭션 처리 중단                      │
│ - 블록 생성 중단                          │
│ - 사용자 서비스 이용 불가                 │
└──────────────────────────────────────────┘
```

#### 데이터 손실 시나리오

```bash
# 볼륨 손상 시나리오
docker volume rm ops-bedrock_da_data  # 😱 모든 데이터 삭제

# 결과:
# - 모든 L2 트랜잭션 데이터 영구 손실
# - Commitment만 L1에 남음 (데이터는 복구 불가)
# - 사용자 인출 불가능
```

#### 위험

- ✗ **서비스 중단**: 단일 서버 장애로 L2 전체 마비
- ✗ **데이터 손실**: 백업 없이 영구적 손실
- ✗ **복구 불가능**: 데이터 없으면 재구축 불가
- ✗ **금융 손실**: 서비스 중단 시간에 비례한 손실

---

### 7. 스토리지 보안 (심각도: Medium)

#### File Storage 취약점

```go
// op-plasma/cmd/daserver/file.go
type FileStore struct {
    path string
}

func NewFileStore(path string) *FileStore {
    // ❌ 파일 권한 설정 없음
    // ❌ 디렉토리 접근 제어 없음
    // ❌ 파일 암호화 없음
    return &FileStore{path: path}
}
```

**컨테이너 탈출 시나리오**:
```bash
# 컨테이너 취약점 이용하여 호스트 파일 시스템 접근
docker exec -it da-server sh
cd /data
cat * > /tmp/leaked_data.tar.gz  # 모든 데이터 유출
```

#### S3 Storage 취약점

```bash
# docker-compose.yml 또는 .env에 평문 저장
AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE  # ❌ 평문 노출
AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY  # ❌
```

**S3 버킷 탈취 시나리오**:
```bash
# .env 파일 유출 시
export AWS_ACCESS_KEY_ID=<유출된 키>
export AWS_SECRET_ACCESS_KEY=<유출된 시크릿>

# S3 버킷 전체 다운로드
aws s3 sync s3://da-server-bucket ./stolen_data/

# 데이터 변조
echo "fake_data" > modified.bin
aws s3 cp modified.bin s3://da-server-bucket/0x1234.../
```

#### 위험

- ✗ **파일 시스템 노출**: 컨테이너 탈출 시 전체 데이터 유출
- ✗ **AWS 자격 증명 유출**: S3 버킷 완전 장악
- ✗ **데이터 무결성**: 외부에서 데이터 변조 가능
- ✗ **규정 준수 실패**: GDPR, 개인정보보호법 위반

---

### 8. 로깅 및 감사 추적 부족 (심각도: Medium)

#### 문제점

```go
// op-plasma/daserver.go:128
func (d *DAServer) HandlePut(w http.ResponseWriter, r *http.Request) {
    d.log.Info("PUT", "url", r.URL)
    // ❌ 요청자 IP 주소 로깅 없음
    // ❌ User-Agent 로깅 없음
    // ❌ 요청 크기 로깅 없음
    // ❌ 처리 시간 로깅 없음
    // ❌ 실패 원인 상세 로깅 없음
}
```

#### 부족한 감사 정보

```go
// 현재 로그 예시
INFO PUT url=/put/0x1234...

// 필요한 로그 정보 (없음):
// - 누가? (IP, User-Agent)
// - 언제? (정확한 타임스탬프)
// - 무엇을? (데이터 크기, 해시)
// - 결과는? (성공/실패, 응답 시간)
// - 왜 실패? (에러 상세)
```

#### 공격 탐지 불가 시나리오

```
공격 발생:
  ├─ DDoS 공격 진행 중 🔥
  ├─ 이상 트래픽 발생
  └─ 하지만 로그에 IP 정보 없음
      └─ 공격 출처 파악 불가
      └─ 차단 불가능

사후 분석:
  ├─ 침해 사고 발생
  ├─ 포렌식 조사 필요
  └─ 하지만 감사 추적 없음
      └─ 원인 분석 불가
      └─ 책임 소재 불명확
```

#### 위험

- ✗ **공격 탐지 불가**: 실시간 이상 행위 탐지 불가능
- ✗ **포렌식 불가능**: 사고 발생 시 원인 분석 불가
- ✗ **규정 준수 실패**: 감사 증적 부재로 규제 위반
- ✗ **책임 추적 불가**: 내부/외부 공격자 식별 불가

---

## 공격 시나리오

### 시나리오 1: 종합 공격 - L2 마비

```
공격자 목표: L2 서비스 완전 중단

1단계: 정찰 (Reconnaissance)
  └─ nmap으로 DA Server 포트 스캔
  └─ 인증 없음 확인
  └─ API 엔드포인트 매핑

2단계: DDoS 공격 (Availability Attack)
  └─ HTTP Flood: 초당 10만 GET 요청
  └─ Storage Exhaustion: 대용량 PUT 요청
  └─ DA Server CPU/메모리 고갈

3단계: Challenge Window 이용
  └─ Sequencer가 부정 commitment 제출
  └─ DA Server 다운으로 챌린지 불가
  └─ Challenge window 만료

4단계: 결과
  └─ 잘못된 상태 최종 확정
  └─ 사용자 자금 인출 불가
  └─ L2 신뢰 붕괴

예상 피해:
  - 서비스 중단 시간: 수 시간 ~ 수 일
  - 금전적 손실: 수억 ~ 수천억 원
  - 평판 손실: 회복 불가능
```

### 시나리오 2: 내부자 공격 - 데이터 유출

```
공격자: 내부 직원 또는 인프라 접근 권한 보유자

1단계: 접근
  └─ Docker 호스트 서버 로그인
  └─ da-server 컨테이너 접근

2단계: 데이터 탈취
  └─ 볼륨 마운트 경로 확인
  └─ 모든 트랜잭션 데이터 복사
  └─ S3 자격 증명 탈취

3단계: 데이터 판매
  └─ 다크웹에 데이터 판매
  └─ 경쟁사에 거래 패턴 유출
  └─ 사용자 프라이버시 침해

4단계: 증거 인멸
  └─ 로그 없어 추적 불가
  └─ 감사 추적 없음

예상 피해:
  - GDPR 위반: 최대 매출의 4%
  - 평판 손실
  - 사용자 이탈
  - 법적 소송
```

### 시나리오 3: 랜섬웨어 공격

```
공격자: 랜섬웨어 그룹

1단계: 침투
  └─ 피싱 이메일로 내부망 침투
  └─ DA Server 접근 권한 획득

2단계: 데이터 암호화
  └─ S3 버킷 전체 암호화
  └─ 또는 파일 스토리지 암호화
  └─ 백업도 함께 암호화 (백업이 있다면)

3단계: 랜섬 요구
  └─ 비트코인으로 100 BTC 요구
  └─ 48시간 내 지불하지 않으면 데이터 삭제

4단계: 결과
  └─ 지불해도 데이터 복구 보장 없음
  └─ 지불 거부 시 모든 데이터 손실
  └─ L2 서비스 영구 중단

예상 피해:
  - 랜섬 금액: 수십억 원
  - 복구 비용: 수억 원
  - 평판 손실: 측정 불가
```

---

## 보안 강화 방안

### 1. 인증/인가 구현

#### 1.1 API 키 기반 인증

```go
// 개선된 코드 예시
type DAServer struct {
    // ... 기존 필드
    apiKeys map[string]bool  // API 키 화이트리스트
}

func (d *DAServer) authenticate(r *http.Request) error {
    apiKey := r.Header.Get("X-API-Key")
    if apiKey == "" {
        return errors.New("missing API key")
    }

    if !d.apiKeys[apiKey] {
        return errors.New("invalid API key")
    }

    return nil
}

func (d *DAServer) HandleGet(w http.ResponseWriter, r *http.Request) {
    // 인증 체크
    if err := d.authenticate(r); err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    // ... 기존 로직
}
```

#### 1.2 Docker Compose 설정

```yaml
da-server:
  image: tokamaknetwork/thanos-da-server
  environment:
    # API 키 설정
    DA_SERVER_API_KEYS: "${DA_API_KEY_1},${DA_API_KEY_2}"
    # 또는 JWT 기반
    DA_SERVER_JWT_SECRET: "${JWT_SECRET}"
    DA_SERVER_JWT_ISSUER: "tokamak-thanos"
    DA_SERVER_JWT_AUDIENCE: "da-clients"
```

#### 1.3 클라이언트 설정

```yaml
op-batcher:
  environment:
    # API 키 주입
    OP_BATCHER_PLASMA_DA_API_KEY: "${DA_API_KEY}"

op-node:
  environment:
    OP_NODE_PLASMA_DA_API_KEY: "${DA_API_KEY}"
```

---

### 2. TLS/mTLS 구현

#### 2.1 TLS 인증서 생성

```bash
# 자체 서명 인증서 생성 (개발용)
mkdir -p ./certs

# CA 생성
openssl genrsa -out certs/ca-key.pem 4096
openssl req -new -x509 -days 365 -key certs/ca-key.pem \
  -out certs/ca.pem -subj "/CN=Tokamak-Thanos-CA"

# 서버 인증서 생성
openssl genrsa -out certs/server-key.pem 4096
openssl req -new -key certs/server-key.pem \
  -out certs/server.csr -subj "/CN=da-server"
openssl x509 -req -days 365 -in certs/server.csr \
  -CA certs/ca.pem -CAkey certs/ca-key.pem \
  -CAcreateserial -out certs/server.pem

# 클라이언트 인증서 생성 (mTLS용)
openssl genrsa -out certs/client-key.pem 4096
openssl req -new -key certs/client-key.pem \
  -out certs/client.csr -subj "/CN=da-client"
openssl x509 -req -days 365 -in certs/client.csr \
  -CA certs/ca.pem -CAkey certs/ca-key.pem \
  -CAcreateserial -out certs/client.pem
```

#### 2.2 Docker Compose 설정

```yaml
da-server:
  image: tokamaknetwork/thanos-da-server
  ports:
    - "3100:3100"
  command: >
    da-server
      --addr=0.0.0.0
      --port=3100
      --tls-cert=/certs/server.pem
      --tls-key=/certs/server-key.pem
      --client-ca=/certs/ca.pem  # mTLS 활성화
      --tls-verify-client=true
  volumes:
    - ./certs:/certs:ro
    - da_data:/data

op-batcher:
  environment:
    OP_BATCHER_PLASMA_DA_SERVER: "https://da-server:3100"  # HTTPS
    OP_BATCHER_PLASMA_DA_TLS_CERT: "/certs/client.pem"
    OP_BATCHER_PLASMA_DA_TLS_KEY: "/certs/client-key.pem"
    OP_BATCHER_PLASMA_DA_TLS_CA: "/certs/ca.pem"
  volumes:
    - ./certs:/certs:ro
```

---

### 3. Rate Limiting 및 보호 계층

#### 3.1 Nginx를 통한 Rate Limiting

```nginx
# nginx.conf
upstream da_servers {
    least_conn;
    server da-server-1:3100 max_fails=3 fail_timeout=30s;
    server da-server-2:3100 max_fails=3 fail_timeout=30s;
}

# Rate Limiting zones
limit_req_zone $binary_remote_addr zone=da_get:10m rate=100r/s;
limit_req_zone $binary_remote_addr zone=da_put:10m rate=10r/s;
limit_conn_zone $binary_remote_addr zone=da_conn:10m;

server {
    listen 443 ssl http2;
    server_name da-server.example.com;

    # TLS 설정
    ssl_certificate /etc/nginx/certs/server.pem;
    ssl_certificate_key /etc/nginx/certs/server-key.pem;
    ssl_client_certificate /etc/nginx/certs/ca.pem;
    ssl_verify_client on;
    ssl_protocols TLSv1.2 TLSv1.3;

    # 보안 헤더
    add_header Strict-Transport-Security "max-age=31536000" always;
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;

    # GET 엔드포인트
    location /get/ {
        limit_req zone=da_get burst=50 nodelay;
        limit_conn da_conn 10;
        proxy_pass http://da_servers;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # PUT 엔드포인트
    location /put/ {
        limit_req zone=da_put burst=10 nodelay;
        limit_conn da_conn 5;

        # 요청 크기 제한
        client_max_body_size 10M;

        proxy_pass http://da_servers;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

#### 3.2 Docker Compose 설정

```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    depends_on:
      - da-server-1
      - da-server-2

  da-server-1:
    image: tokamaknetwork/thanos-da-server
    # 외부 포트 노출하지 않음
    expose:
      - "3100"
    # ... 기타 설정

  da-server-2:
    image: tokamaknetwork/thanos-da-server
    expose:
      - "3100"
    # ... 기타 설정
```

---

### 4. 고가용성 (HA) 구성

#### 4.1 다중 DA Server 구성

```yaml
version: '3.4'

services:
  # DA Server 인스턴스 1
  da-server-1:
    image: tokamaknetwork/thanos-da-server
    hostname: da-server-1
    environment:
      DA_SERVER_ID: "server-1"
    volumes:
      - da_data_1:/data
    networks:
      - da-network

  # DA Server 인스턴스 2
  da-server-2:
    image: tokamaknetwork/thanos-da-server
    hostname: da-server-2
    environment:
      DA_SERVER_ID: "server-2"
    volumes:
      - da_data_2:/data
    networks:
      - da-network

  # DA Server 인스턴스 3 (최소 3개 권장)
  da-server-3:
    image: tokamaknetwork/thanos-da-server
    hostname: da-server-3
    environment:
      DA_SERVER_ID: "server-3"
    volumes:
      - da_data_3:/data
    networks:
      - da-network

  # Nginx Load Balancer
  nginx-lb:
    image: nginx:alpine
    ports:
      - "443:443"
    volumes:
      - ./nginx-lb.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    depends_on:
      - da-server-1
      - da-server-2
      - da-server-3
    networks:
      - da-network
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "-", "https://localhost:443/health"]
      interval: 10s
      timeout: 5s
      retries: 3

volumes:
  da_data_1:
  da_data_2:
  da_data_3:

networks:
  da-network:
```

#### 4.2 S3 복제 설정

```yaml
da-server:
  environment:
    # Primary S3 bucket
    AWS_S3_BUCKET: "da-primary-bucket"
    AWS_REGION: "us-east-1"

    # Replication 설정
    AWS_S3_REPLICATION_ENABLED: "true"
    AWS_S3_REPLICATION_BUCKET: "da-replica-bucket"
    AWS_S3_REPLICATION_REGION: "eu-west-1"
```

**AWS S3 복제 규칙 설정**:
```bash
# S3 버전 관리 활성화
aws s3api put-bucket-versioning \
  --bucket da-primary-bucket \
  --versioning-configuration Status=Enabled

# 복제 규칙 생성
aws s3api put-bucket-replication \
  --bucket da-primary-bucket \
  --replication-configuration file://replication.json
```

```json
// replication.json
{
  "Role": "arn:aws:iam::ACCOUNT:role/s3-replication-role",
  "Rules": [
    {
      "Status": "Enabled",
      "Priority": 1,
      "DeleteMarkerReplication": { "Status": "Enabled" },
      "Filter": {},
      "Destination": {
        "Bucket": "arn:aws:s3:::da-replica-bucket",
        "ReplicationTime": {
          "Status": "Enabled",
          "Time": { "Minutes": 15 }
        }
      }
    }
  ]
}
```

---

### 5. 백업 및 복구 전략

#### 5.1 자동 백업 스크립트

```bash
#!/bin/bash
# backup-da-server.sh

set -euo pipefail

BACKUP_DIR="/backups/da-server"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# 백업 디렉토리 생성
mkdir -p "$BACKUP_DIR"

# Docker 볼륨 백업
echo "Backing up DA server volumes..."
for volume in da_data_1 da_data_2 da_data_3; do
    docker run --rm \
        -v "${volume}:/data" \
        -v "${BACKUP_DIR}:/backup" \
        alpine \
        tar czf "/backup/${volume}_${TIMESTAMP}.tar.gz" -C /data .

    echo "Backed up $volume"
done

# S3 버킷 백업 (선택적)
if [ "${BACKUP_S3:-false}" = "true" ]; then
    echo "Backing up S3 bucket..."
    aws s3 sync \
        s3://da-primary-bucket \
        "${BACKUP_DIR}/s3_${TIMESTAMP}/" \
        --storage-class GLACIER
fi

# 오래된 백업 삭제
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: ${BACKUP_DIR}"
```

#### 5.2 Cron 작업 등록

```bash
# crontab -e
# 매일 새벽 2시에 백업
0 2 * * * /opt/scripts/backup-da-server.sh >> /var/log/da-backup.log 2>&1
```

#### 5.3 복구 스크립트

```bash
#!/bin/bash
# restore-da-server.sh

set -euo pipefail

if [ $# -lt 1 ]; then
    echo "Usage: $0 <backup_timestamp>"
    exit 1
fi

BACKUP_TIMESTAMP=$1
BACKUP_DIR="/backups/da-server"

# 시스템 중지
echo "Stopping DA server..."
cd /path/to/ops-bedrock
docker compose down

# 볼륨 복구
for volume in da_data_1 da_data_2 da_data_3; do
    BACKUP_FILE="${BACKUP_DIR}/${volume}_${BACKUP_TIMESTAMP}.tar.gz"

    if [ ! -f "$BACKUP_FILE" ]; then
        echo "Error: Backup file not found: $BACKUP_FILE"
        exit 1
    fi

    echo "Restoring $volume..."
    docker run --rm \
        -v "${volume}:/data" \
        -v "${BACKUP_DIR}:/backup" \
        alpine \
        sh -c "rm -rf /data/* && tar xzf /backup/${volume}_${BACKUP_TIMESTAMP}.tar.gz -C /data"
done

# 시스템 재시작
echo "Restarting DA server..."
docker compose up -d

echo "Restore completed. Timestamp: ${BACKUP_TIMESTAMP}"
```

---

### 6. 모니터링 및 알림

#### 6.1 Prometheus 메트릭 수집

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

# Alertmanager 설정
alerting:
  alertmanagers:
    - static_configs:
        - targets: ['alertmanager:9093']

# 메트릭 수집 대상
scrape_configs:
  - job_name: 'da-servers'
    static_configs:
      - targets:
        - 'da-server-1:9090'
        - 'da-server-2:9090'
        - 'da-server-3:9090'

  - job_name: 'nginx-lb'
    static_configs:
      - targets: ['nginx-lb:9113']  # nginx-prometheus-exporter
```

#### 6.2 알림 규칙

```yaml
# alert-rules.yml
groups:
  - name: da_server_alerts
    interval: 30s
    rules:
      # DA Server 다운 알림
      - alert: DAServerDown
        expr: up{job="da-servers"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "DA Server {{ $labels.instance }} is down"
          description: "DA Server has been down for more than 1 minute"

      # 높은 에러율
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate on {{ $labels.instance }}"

      # 디스크 사용량 높음
      - alert: HighDiskUsage
        expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) < 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Disk usage above 90% on {{ $labels.instance }}"

      # 메모리 사용량 높음
      - alert: HighMemoryUsage
        expr: (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) < 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Memory usage above 90% on {{ $labels.instance }}"

      # 높은 응답 시간
      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "95th percentile response time > 1s"
```

#### 6.3 Alertmanager 설정

```yaml
# alertmanager.yml
global:
  resolve_timeout: 5m
  slack_api_url: '${SLACK_WEBHOOK_URL}'

route:
  receiver: 'default'
  group_by: ['alertname', 'cluster']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  routes:
    - match:
        severity: critical
      receiver: 'critical'
      continue: true

receivers:
  - name: 'default'
    slack_configs:
      - channel: '#da-server-alerts'
        title: 'DA Server Alert'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}{{ end }}'

  - name: 'critical'
    slack_configs:
      - channel: '#da-server-critical'
        title: '🚨 CRITICAL: DA Server Alert'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}{{ end }}'
    pagerduty_configs:
      - service_key: '${PAGERDUTY_SERVICE_KEY}'
```

---

### 7. 보안 로깅 개선

#### 7.1 구조화된 로깅

```go
// 개선된 로깅 예시
type RequestLog struct {
    Timestamp   time.Time
    Method      string
    Path        string
    RemoteAddr  string
    UserAgent   string
    RequestSize int64
    ResponseCode int
    Duration    time.Duration
    Error       string
}

func (d *DAServer) logRequest(r *http.Request, code int, duration time.Duration, err error) {
    log := RequestLog{
        Timestamp:   time.Now(),
        Method:      r.Method,
        Path:        r.URL.Path,
        RemoteAddr:  r.RemoteAddr,
        UserAgent:   r.UserAgent(),
        RequestSize: r.ContentLength,
        ResponseCode: code,
        Duration:    duration,
    }

    if err != nil {
        log.Error = err.Error()
    }

    // JSON 형식으로 로깅
    d.logger.Info("request", "log", log)
}
```

#### 7.2 ELK Stack 연동

```yaml
# docker-compose.monitoring.yml
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  logstash:
    image: docker.elastic.co/logstash/logstash:8.11.0
    volumes:
      - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf:ro
    depends_on:
      - elasticsearch

  kibana:
    image: docker.elastic.co/kibana/kibana:8.11.0
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch

  da-server:
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "10"
        labels: "service=da-server"
```

---

### 8. 입력 검증 강화

```go
// 개선된 입력 검증
const (
    MaxInputSize = 10 * 1024 * 1024  // 10MB
    MaxRequestsPerIP = 100            // per minute
)

func (d *DAServer) validatePutRequest(r *http.Request) error {
    // Content-Length 확인
    if r.ContentLength > MaxInputSize {
        return fmt.Errorf("request size %d exceeds limit %d",
            r.ContentLength, MaxInputSize)
    }

    // Rate limit 확인
    if d.rateLimiter.Exceeded(r.RemoteAddr) {
        return fmt.Errorf("rate limit exceeded")
    }

    // Content-Type 확인
    contentType := r.Header.Get("Content-Type")
    if contentType != "application/octet-stream" {
        return fmt.Errorf("invalid content type: %s", contentType)
    }

    return nil
}

func (d *DAServer) HandlePut(w http.ResponseWriter, r *http.Request) {
    // 입력 검증
    if err := d.validatePutRequest(r); err != nil {
        d.log.Warn("invalid request", "error", err, "ip", r.RemoteAddr)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // 제한된 크기만 읽기
    limitedReader := io.LimitReader(r.Body, MaxInputSize)
    input, err := io.ReadAll(limitedReader)
    if err != nil {
        d.log.Error("failed to read body", "error", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // ... 기존 로직
}
```

---

## 보안 체크리스트

### 프로덕션 배포 전 필수 항목

#### ✅ Critical (심각도: 높음)

- [ ] **인증 구현**
  - [ ] API 키 또는 JWT 토큰 기반 인증
  - [ ] 강력한 API 키 생성 (최소 32바이트 랜덤)
  - [ ] 정기적인 키 로테이션 정책 수립

- [ ] **TLS/mTLS 적용**
  - [ ] TLS 1.2 이상 사용
  - [ ] 유효한 인증서 사용 (Let's Encrypt 등)
  - [ ] mTLS로 클라이언트 인증 강화

- [ ] **고가용성 구성**
  - [ ] 최소 3개 이상의 DA Server 인스턴스
  - [ ] Load Balancer 구성
  - [ ] Health check 설정
  - [ ] Auto-scaling 정책

- [ ] **백업 전략**
  - [ ] 정기 백업 스케줄 (최소 일 1회)
  - [ ] 지리적으로 분산된 백업 저장소
  - [ ] 정기적인 복구 테스트 (월 1회)
  - [ ] 백업 암호화

- [ ] **Challenge Window 보호**
  - [ ] 다중 DA provider 설정
  - [ ] Challenge 모니터링 시스템
  - [ ] 자동 알림 설정

#### ✅ High (심각도: 중상)

- [ ] **Rate Limiting**
  - [ ] GET 요청 제한 (권장: 100 req/s per IP)
  - [ ] PUT 요청 제한 (권장: 10 req/s per IP)
  - [ ] 동시 연결 수 제한
  - [ ] Burst 허용량 설정

- [ ] **DDoS 방어**
  - [ ] WAF (Web Application Firewall) 배치
  - [ ] CloudFlare 또는 AWS Shield 사용
  - [ ] IP 블랙리스트/화이트리스트 관리

- [ ] **입력 검증**
  - [ ] 최대 요청 크기 제한 (권장: 10MB)
  - [ ] Content-Type 검증
  - [ ] 악의적 페이로드 필터링

- [ ] **모니터링 및 알림**
  - [ ] Prometheus + Grafana 구성
  - [ ] Alertmanager로 실시간 알림
  - [ ] 24/7 on-call 체계 구축
  - [ ] 주요 메트릭 대시보드

#### ✅ Medium (심각도: 중)

- [ ] **스토리지 보안**
  - [ ] S3 버킷 암호화 (AES-256)
  - [ ] S3 버킷 정책 최소 권한 원칙
  - [ ] VPC Endpoint 사용 (AWS)
  - [ ] 파일 스토리지 권한 제한 (chmod 600)

- [ ] **로깅 및 감사**
  - [ ] 구조화된 로깅 (JSON)
  - [ ] 중앙 집중식 로그 관리 (ELK Stack)
  - [ ] 로그 보존 정책 (최소 90일)
  - [ ] 로그 무결성 보장

- [ ] **네트워크 보안**
  - [ ] 방화벽 규칙 설정
  - [ ] VPN 또는 Private Network 사용
  - [ ] 불필요한 포트 차단
  - [ ] 네트워크 세그먼테이션

- [ ] **컨테이너 보안**
  - [ ] 최소 권한 컨테이너 실행
  - [ ] Read-only 파일 시스템 (가능한 경우)
  - [ ] 보안 스캔 (Trivy, Clair 등)
  - [ ] 정기적인 이미지 업데이트

#### ✅ Low (권장 사항)

- [ ] **고급 보안**
  - [ ] Intrusion Detection System (IDS)
  - [ ] Security Information and Event Management (SIEM)
  - [ ] Honeypot 배치
  - [ ] 정기 침투 테스트

- [ ] **재해 복구**
  - [ ] 재해 복구 계획 (DRP) 수립
  - [ ] 지리적 복제 (Multi-region)
  - [ ] Failover 자동화
  - [ ] RTO/RPO 목표 설정

- [ ] **규정 준수**
  - [ ] GDPR 준수 (EU 사용자 대상)
  - [ ] 개인정보보호법 준수
  - [ ] 정기 보안 감사
  - [ ] 보안 정책 문서화

---

### 운영 체크리스트

#### 일일 점검

- [ ] DA Server 상태 확인
- [ ] 에러 로그 검토
- [ ] 디스크 사용량 확인
- [ ] 메모리 사용량 확인
- [ ] 응답 시간 확인

#### 주간 점검

- [ ] 보안 로그 분석
- [ ] 이상 트래픽 패턴 확인
- [ ] 백업 상태 검증
- [ ] 인증서 만료일 확인
- [ ] 성능 메트릭 분석

#### 월간 점검

- [ ] 복구 테스트 수행
- [ ] 보안 패치 적용
- [ ] API 키 로테이션
- [ ] 접근 권한 검토
- [ ] 보안 감사 보고서 작성

#### 분기별 점검

- [ ] 재해 복구 훈련
- [ ] 보안 정책 업데이트
- [ ] 침투 테스트 수행
- [ ] 규정 준수 검토
- [ ] 아키텍처 보안 검토

---

## 참고 자료

### 공식 문서

- [Optimism Plasma Mode](https://specs.optimism.io/experimental/plasma/)
- [Data Availability Challenge Spec](https://specs.optimism.io/experimental/fault-proof/)
- [AWS S3 Security Best Practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/security-best-practices.html)

### 관련 파일

- `op-plasma/daserver.go`: DA Server 구현
- `op-plasma/damgr.go`: Challenge 관리 로직
- `packages/contracts-bedrock/src/L1/DataAvailabilityChallenge.sol`: Challenge 컨트랙트
- `ops-bedrock/docker-compose.yml`: 서비스 구성

### 보안 도구

- **Trivy**: 컨테이너 이미지 보안 스캔
- **Falco**: 런타임 보안 모니터링
- **OWASP ZAP**: 웹 애플리케이션 보안 스캔
- **Nmap**: 네트워크 스캔
- **Wireshark**: 패킷 분석

### 표준 및 가이드라인

- OWASP Top 10
- NIST Cybersecurity Framework
- CIS Benchmarks
- ISO 27001

---

**대상 프로젝트**: tokamak-thanos (Optimism Fork)

---

## 면책 조항

본 문서는 DA Server의 보안 취약점을 분석하고 개선 방안을 제시하기 위한 목적으로 작성되었습니다. 실제 프로덕션 환경에서는 조직의 보안 정책과 규정 준수 요구사항에 따라 추가적인 보안 조치가 필요할 수 있습니다.

본 문서의 내용을 악의적인 목적으로 사용하는 것을 금지하며, 모든 보안 테스트는 적절한 권한과 승인을 받은 후에만 수행해야 합니다.

