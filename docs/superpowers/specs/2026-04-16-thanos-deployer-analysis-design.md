---
date: 2026-04-16
title: thanos-deployer 배포 로직 분석 설계서
author: Claude Code
status: approved
---

# thanos-deployer 배포 로직 분석 설계서

## 개요

Electron 앱에서 L1 배포 요청이 들어왔을 때, **어떤 경로로 어떤 deployer가 실행되고, 어떤 로직으로 L1 컨트랙트를 배포하며 L2 Genesis를 생성하는지**를 코드 수준에서 완벽하게 분석하고 문서화한다.

## 분석 목표

### 주 질문들
1. Electron UI에서 배포 버튼 클릭 → HTTP 요청 → 어느 백엔드 핸들러가 수신하는가?
2. 백엔드에서 어떤 SDK 함수를 호출하는가?
3. Foundry 스크립트는 어떻게 호출되고, 어떤 순서로 L1 컨트랙트를 배포하는가?
4. op-chain-ops deployer는 어떻게 L2 Genesis를 생성하는가?
5. 각 단계의 입력/출력 데이터 구조는 무엇인가?

## 분석 범위

### 시간적 범위
**시작점**: Electron 배포 UI → `POST /api/v1/stacks/thanos`
**종료점**: `build/genesis.json` + `build/rollup.json` 생성 완료

### 공간적 범위 (코드베이스)
- **trh-platform** — Electron 셸, AWS SSO 통합
- **trh-backend** — Gin 백엔드, 요청 핸들러, TaskManager
- **trh-sdk** — L1 배포 오케스트레이션, Foundry 셸아웃
- **tokamak-thanos** — Foundry 스크립트, op-chain-ops deployer

### 분석 깊이
**코드 수준**:
- 파일 경로 (상대 경로, 레포별 명시)
- 함수/메서드명 + 라인 번호
- 핵심 로직 (의사코드/텍스트 설명)
- 입/출력 데이터 구조 (함수 시그니처, JSON 스키마)
- 호출 관계 (누가 누를 부르는가)

## 산출물 구조

### 1. 주 분석 문서
**파일**: `docs/analysis/thanos-deployer-flow-analysis.md`

**내용**:
- Phase별 상세 분석 (6개 Phase)
  - Phase 1: Electron SSO 인증
  - Phase 2: Web UI 배포 요청
  - Phase 3: trh-backend 영속화 + 비동기 큐잉
  - Phase 4: L1 컨트랙트 배포 (Foundry)
  - Phase 5: L2 Genesis 생성 (op-chain-ops)
  - Phase 6: 결과 반영

- 각 Phase마다:
  - 개략적 설명
  - 핵심 파일 목록
  - 함수별 상세 분석 (파일:라인, 함수명, 입/출력)
  - 데이터 흐름 (JSON/구조체)
  - 호출 시퀀스

- 부록:
  - 통합 호출 그래프 (모든 함수의 부모-자식 관계)
  - 데이터 구조 스키마 (DTO, JSON 포맷)
  - 알려진 함정 및 개선 포인트

### 2. 다이어그램 세트
Mermaid로 작성, SVG 렌더링:

**2a. 전체 시스템 아키텍처** (`system-architecture.mmd`)
- Electron → trh-backend → trh-sdk → Foundry/Go deployer 계층 구조
- 모듈 간 통신 방식 (HTTP, 모듈 import, 셸아웃)

**2b. L1 배포 플로우** (`l1-deploy-flow.mmd`)
- Foundry `start-deploy.sh` 실행
- 컨트랙트 배포 순서 (OptimismPortal → SystemConfig → ... → AnchorStateRegistry)
- 산출물 생성 (deploy.json, rollup.json)

**2c. L2 Genesis 생성 플로우** (`l2-genesis-flow.mmd`)
- Genesis allocs 생성
- Predeploy 배치
- 상태 초기화
- Fault-proof 패칭

**2d. 모듈 호출 그래프** (`call-graph.mmd`)
- 트레이스: `Deploy` → `DeployContracts` → `DeployAWSInfrastructure` 등
- 각 화살표에 파일:라인 명시

**2e. 데이터 변환 흐름** (`data-flow.mmd`)
- 입력: DeployThanosRequest JSON
- 중간: .env 파일, 스크립트 패치 상태
- 출력: genesis.json, rollup.json, deploy.json

### 3. 코드 참조 테이블
**파일**: `docs/analysis/thanos-deployer-code-map.md` (또는 상기 문서의 부록)

예시:
| Phase | 파일 | 함수 | 라인 | 역할 | 입력 | 출력 |
|-------|------|------|------|------|------|------|
| 1 | trh-platform/src/main/aws-auth.ts | startSsoLoginDirect | 335 | SSO 인증 | 없음 | ssoAccessToken |
| 2 | trh-backend/pkg/api/handlers/thanos/deployment.go | Deploy | 32 | 배포 요청 수신 | JSON body | {stackId} |
| 3 | trh-backend/pkg/services/thanos/stack_lifecycle.go | CreateThanosStack | 20 | 스택 생성 + 큐잉 | DeployThanosRequest | stackId |
| 4a | trh-sdk/pkg/stacks/thanos/deploy_contracts.go | DeployContracts | 33 | L1 배포 오케스트레이션 | input struct | {addresses} |
| 4b | tokamak-thanos/start-deploy.sh | (bash) | - | Foundry 호출 | .env | deploy.json |
| 5 | op-chain-ops/deployer/deployer.go | (various) | - | L2 Genesis 생성 | Deploy struct | genesis.json |

## 분석 방식

### 코드 추적 프로세스
1. **진입점 확인** — `POST /api/v1/stacks/thanos` 핸들러
2. **호출 체인 추적** — 각 함수의 의존성을 재귀적으로 따라감
3. **데이터 흐름 매핑** — 입력이 어떻게 변환되어 출력이 되는가
4. **셸아웃 지점 식별** — Foundry, terraform 등 외부 프로세스 호출
5. **다이어그램 작성** — 추적 결과를 시각화
6. **검증** — 현재 코드와 비교 확인

### 도구 사용
- `Grep` / `Read` — 코드 조회 및 읽기
- `LSP` — 함수 정의/참조 위치 탐색 (goToDefinition, findReferences)
- `Bash` — git log, git blame으로 변경 히스토리 확인
- `mcp__mermaid__generate` — Mermaid 다이어그램 생성
- `Write` — 마크다운 문서 작성

### 검증 체크리스트
- [ ] 모든 파일 경로/라인 번호가 현재 코드베이스와 일치
- [ ] 호출 순서가 코드 실행 순서와 일치
- [ ] 다이어그램과 텍스트 설명이 모순되지 않음
- [ ] wiki [[ec2-deploy]]와의 연결성 확인 (기존 문서와 상호참조)
- [ ] 대소문자, 함수명 스펠링 정확함

## 스코프 관리

### 포함 범위 ✓
- Electron → trh-backend → trh-sdk → Foundry/Go 전체 경로
- L1 컨트랙트 배포 상세 (contracts 목록, 배포 순서)
- L2 Genesis 생성 상세 (predeploy 배치, 상태 초기화)
- 에러 핸들링 및 재시도 로직

### 제외 범위 ✗
- Terraform/Helm 상세 (별도 문서 [[tokamak-thanos-stack]] 참조)
- 컨트랙트 Solidity 구현 (배포 메커니즘이 아님)
- 로컬 Docker 배포 (AWS 경로만 상세, 로컬은 비교 관점)
- CrossTrade/DRB 모듈 배포 (Preset 레이어, 별도 분석 가능)

## 성공 기준

완성된 분석이 다음을 만족해야 함:

1. **완전성**: Electron 클릭 → genesis.json 생성까지 모든 함수/파일이 추적됨
2. **정확성**: 코드베이스와 100% 일치 (라인 번호, 함수명, 로직)
3. **명확성**: 비기술자도 흐름을 이해 가능 (플로우 차트 + 텍스트)
4. **실용성**: 향후 버그 분석/최적화에 즉시 활용 가능
5. **유지보수성**: wiki에 저장되어 커밋/PR 업데이트와 함께 유지

## 일정 및 체크포인트

- **Phase 4 (L1 배포)**: 먼저 분석 (가장 복잡)
- **Phase 5 (L2 Genesis)**: 다음 (op-chain-ops deployer)
- **Phase 1-3**: 기존 wiki [[ec2-deploy]]와 비교, 코드 수준 추가
- **다이어그램**: 각 Phase 분석 후 즉시 작성
- **통합 검증**: 모든 Phase 완료 후

## 다음 단계

이 설계서 승인 후:
1. writing-plans 스킬 호출 → 상세 구현 계획 작성
2. 코드 분석 착수 (Phase별 추적)
3. 문서 + 다이어그램 병렬 작성
4. wiki 업데이트 및 commit
