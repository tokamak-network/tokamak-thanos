# Phase 3 분석: 백엔드 배포 큐잉 및 오케스트레이션

## 목차
1. [개요](#개요)
2. [요청 수신 및 검증](#요청-수신-및-검증)
3. [스택 라이프사이클 관리](#스택-라이프사이클-관리)
4. [배포 오케스트레이션](#배포-오케스트레이션)
5. [작업 큐 시스템](#작업-큐-시스템)
6. [데이터베이스 스키마](#데이터베이스-스키마)
7. [상태 머신 및 에러 처리](#상태-머신-및-에러-처리)
8. [호출 시퀀스](#호출-시퀀스)

---

## 개요

Phase 3는 클라이언트의 배포 요청이 백엔드에서 어떻게 수신, 검증, 데이터베이스에 저장되고, 비동기 작업 큐를 통해 오케스트레이션되는지를 다룬다. 핵심 특징은:

- **Preset 기반 추상화**: 사용자는 미리 정의된 배포 템플릿을 선택하고 매개변수를 재정의
- **원자적 다중 엔티티 생성**: 스택, 배포, 통합 엔티티를 트랜잭션으로 원자적으로 생성
- **메모리 기반 작업 큐**: Redis/RabbitMQ 없이 Go 채널 기반 워커 풀
- **단계별 배포 오케스트레이션**: L1 컨트랙트 배포 → AWS 인프라 배포 순차 실행
- **컨텍스트 기반 취소**: 장시간 실행 작업의 우아한 중단 지원
- **실시간 로그 수집**: 배포 진행 중 로그 스트리밍

---

## 요청 수신 및 검증

### HTTP 핸들러: PresetDeploy

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/api/handlers/thanos/presets.go`
**함수**: `PresetDeploy` (라인 78-104)

```go
func (h *PresetHandler) PresetDeploy(c *gin.Context) {
	var req dtos.PresetDeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 요청 검증
	if req.PresetID == "" || req.StackName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "preset_id and stack_name required"})
		return
	}

	// 2. 서비스 호출: Preset → DeployThanosRequest 변환
	resp, err := h.service.CreateThanosStackFromPreset(c.Request.Context(), req)
	if err != nil {
		// 버그: 에러 시 resp가 nil일 수 있음
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
```

**요청 구조** (`dtos.PresetDeployRequest`):
```
{
  preset_id: string        // 배포 템플릿 ID
  stack_name: string       // 생성할 스택 이름
  chain_id: integer        // 체인 ID (선택적, preset 기본값 사용)
  seed_phrase: string      // 계정 유도용 시드 구문
  provider: string         // AWS 공급자 설정
  overrides: {             // Preset 기본값 재정의
    l1_rpc: string
    l2_config: object
    ...
  }
}
```

**검증 단계**:
1. JSON 파싱 검증
2. 필수 필드(`preset_id`, `stack_name`) 확인
3. Preset 존재 여부 확인 (서비스 레이어에서)
4. Seed phrase 형식 검증 (12/24 단어)
5. Provider 설정 유효성 확인

**응답 구조** (`entities.Response`):
```go
type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	StackID    string      `json:"stack_id"`
	DeploymentID string    `json:"deployment_id"`
}
```

---

## 스택 라이프사이클 관리

### 1단계: Preset → 배포 요청 변환

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/preset_deploy.go`
**함수**: `CreateThanosStackFromPreset` (라인 18-167)

이 함수는 사용자가 선택한 preset을 로드하고 사용자 입력과 병합하여 최종 `DeployThanosRequest`를 구성한다.

**주요 단계**:

```go
func (s *ThanosStackDeploymentService) CreateThanosStackFromPreset(
	ctx context.Context,
	req dtos.PresetDeployRequest,
) (*entities.Response, error) {
	
	// 1. Preset 템플릿 로드
	preset, err := s.loadPreset(ctx, req.PresetID)
	// 오류: preset이 없거나 비활성화됨
	
	// 2. 요청에서 시드 구문 추출
	seedPhrase := req.SeedPhrase // 공백 분리된 12/24 단어
	
	// 3. 시드 구문에서 계정 유도
	// Seed → HD Wallet → Account derivation path
	//   - m/44'/60'/0'/0/0 → Admin account
	//   - m/44'/60'/0'/0/1 → Sequencer account
	//   - m/44'/60'/0'/0/2 → Batcher account
	//   - m/44'/60'/0'/0/3 → Proposer account
	//   - m/44'/60'/0'/0/4 → Challenger account
	seedDerivedAccounts := deriveSeedAccounts(seedPhrase)
	
	// 4. Preset 기본값과 사용자 재정의 병합
	deployRequest := DeployThanosRequest{
		StackName:    req.StackName,
		ChainID:      coalesce(req.ChainID, preset.DefaultChainID),
		Network:      determineNetwork(req.Provider),
		L1RPC:        coalesce(req.Overrides.L1RPC, preset.L1RPC),
		L2Config:     mergeL2Config(preset.L2Config, req.Overrides),
		// ... 추가 필드
	}
	
	// 5. 시드 구문 안전 삭제 (메모리에서)
	clearSeedPhrase(seedPhrase)
	
	// 6. 서비스 호출: 스택 생성
	response, err := s.CreateThanosStack(ctx, deployRequest)
	
	return response, err
}
```

**계정 유도 프로세스**:

| 경로 | 계정 | 용도 |
|------|------|------|
| m/44'/60'/0'/0/0 | Admin | 스택 관리, 컨트랙트 배포 |
| m/44'/60'/0'/0/1 | Sequencer | 트랜잭션 집계 |
| m/44'/60'/0'/0/2 | Batcher | 배치 제출 |
| m/44'/60'/0'/0/3 | Proposer | 상태 제안 |
| m/44'/60'/0'/0/4 | Challenger | 부정행위 발견 |

**보안 특징**:
- 시드 구문은 메모리에 최소한으로만 유지
- 함수 반환 전 명시적으로 쓰레기 수집
- 유도된 계정만 저장, 시드 구문은 저장하지 않음

---

### 2단계: 스택 및 배포 엔티티 생성

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/stack_lifecycle.go`
**함수**: `CreateThanosStack` (라인 101-273)

이 함수는 데이터베이스 트랜잭션 내에서 스택, 배포, 통합 엔티티를 원자적으로 생성한다.

```go
func (s *ThanosStackDeploymentService) CreateThanosStack(
	ctx context.Context,
	request DeployThanosRequest,
) (*entities.Response, error) {
	
	// 1. 로컬 배포 사전 검사 (포트 충돌 방지)
	if request.Network == entities.LocalDevnet {
		if err := s.checkNoActiveLocalStack(); err != nil {
			return nil, fmt.Errorf("local stack already active: %w", err)
		}
	}
	
	// 2. 트랜잭션 시작 (원자성 보장)
	tx := s.deploymentRepo.BeginTx(ctx)
	
	// 3. StackEntity 생성
	stackEntity := &entities.StackEntity{
		ID:           uuid.New(),
		Name:         request.StackName,
		Type:         "Thanos",
		Network:      request.Network,
		Status:       entities.StackStatusPending,  // 초기 상태
		Config:       marshalConfig(request),
		Metadata:     initializeMetadata(request),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	err := s.stackRepo.CreateStackByTx(tx, stackEntity)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create stack: %w", err)
	}
	
	// 4. DeploymentEntity 생성 (단계별)
	// 4.1 L1 컨트랙트 배포 단계
	l1Deployment := &entities.DeploymentEntity{
		ID:       uuid.New(),
		StackID:  &stackEntity.ID,
		Step:     "deploy_l1_contracts",
		Status:   entities.DeploymentStatusPending,
		Config:   marshalL1Config(request),
		CreatedAt: time.Now().UTC(),
	}
	err = s.deploymentRepo.CreateDeployment(tx, l1Deployment)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create L1 deployment: %w", err)
	}
	
	// 4.2 AWS 인프라 배포 단계
	awsDeployment := &entities.DeploymentEntity{
		ID:       uuid.New(),
		StackID:  &stackEntity.ID,
		Step:     "deploy_aws_infrastructure",
		Status:   entities.DeploymentStatusPending,
		Config:   marshalAWSConfig(request),
		CreatedAt: time.Now().UTC(),
	}
	err = s.deploymentRepo.CreateDeployment(tx, awsDeployment)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create AWS deployment: %w", err)
	}
	
	// 5. IntegrationEntity 생성 (활성화된 모듈)
	if request.EnableBridge {
		bridgeIntegration := &entities.IntegrationEntity{
			ID:        uuid.New(),
			StackID:   stackEntity.ID,
			Type:      "Bridge",
			Status:    entities.DeploymentStatusPending,
			CreatedAt: time.Now().UTC(),
		}
		err = s.integrationRepo.CreateIntegration(tx, bridgeIntegration)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create bridge integration: %w", err)
		}
	}
	
	// 유사하게 BlockExplorer, Monitoring 등 생성...
	
	// 6. 트랜잭션 커밋
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	// 7. 배포 작업 큐잉
	taskID := fmt.Sprintf("deploy_%s", stackEntity.ID.String())
	s.taskManager.AddTask(taskID, func(ctx context.Context) {
		s.deploy(ctx, stackEntity.ID)
	})
	
	// 8. 응답 반환
	return &entities.Response{
		Success:      true,
		Message:      "Stack deployment queued",
		StackID:      stackEntity.ID.String(),
		DeploymentID: l1Deployment.ID.String(),
	}, nil
}
```

**상태 초기화**:

| 엔티티 | 초기 상태 | 의미 |
|--------|----------|------|
| StackEntity | Pending | 배포 대기 중 |
| DeploymentEntity (L1) | Pending | L1 배포 단계 대기 |
| DeploymentEntity (AWS) | Pending | AWS 배포 단계 대기 |
| IntegrationEntity | Pending | 통합 모듈 설치 대기 |

**트랜잭션 보장**:
- 모든 엔티티 생성이 성공하거나 모두 실패
- 부분적 생성 상태 방지
- Rollback 시 정확한 에러 추적

---

## 배포 오케스트레이션

### 배포 실행: deploy() 함수

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go`
**함수**: `deploy` (라인 31-437)

작업 워커에 의해 호출되는 진입점. 스택 상태를 Deploying으로 전환하고 배포 단계를 순차 실행한다.

```go
func (s *ThanosStackDeploymentService) deploy(
	ctx context.Context,
	stackID uuid.UUID,
) {
	logger.Info("Starting deployment", zap.String("stack_id", stackID.String()))
	
	// 1. 스택 상태 → Deploying
	stack, err := s.stackRepo.GetStackByID(stackID.String())
	if err != nil || stack == nil {
		logger.Error("Stack not found", zap.String("stack_id", stackID.String()))
		s.updateStackStatus(stackID, entities.StackStatusFailedToDeploy)
		return
	}
	
	err = s.stackRepo.UpdateStatus(stackID, entities.StackStatusDeploying)
	if err != nil {
		logger.Error("Failed to update status", zap.Error(err))
		return
	}
	
	// 2. 배포 단계 실행
	err = s.executeDeployments(ctx, stackID)
	if err != nil {
		logger.Error("Deployment failed", zap.String("stack_id", stackID.String()), zap.Error(err))
		s.updateStackStatus(stackID, entities.StackStatusFailedToDeploy)
		return
	}
	
	// 3. 스택 상태 → Deployed (성공)
	err = s.stackRepo.UpdateStatus(stackID, entities.StackStatusDeployed)
	if err != nil {
		logger.Error("Failed to mark as deployed", zap.Error(err))
		return
	}
	
	logger.Info("Deployment completed successfully", zap.String("stack_id", stackID.String()))
}
```

### 배포 단계 오케스트레이션: executeDeployments()

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/services/thanos/deployment.go`
**함수**: `executeDeployments` (라인 439-650+)

L1 컨트랙트 배포와 AWS 인프라 배포를 순차적으로 실행한다. 각 단계는 로그 수집, 컨텍스트 취소, 메타데이터 업데이트를 포함한다.

```go
func (s *ThanosStackDeploymentService) executeDeployments(
	ctx context.Context,
	stackID uuid.UUID,
) error {
	
	// 1. 배포 단계 목록 조회 (Pending 상태)
	deployments, err := s.deploymentRepo.GetDeploymentsByStackIDAndStatus(
		stackID,
		entities.DeploymentStatusPending,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch deployments: %w", err)
	}
	
	// 2. 단계 정렬 (L1 먼저, AWS 나중)
	sort.Slice(deployments, func(i, j int) bool {
		stepOrder := map[string]int{
			"deploy_l1_contracts":          1,
			"deploy_aws_infrastructure":    2,
		}
		return stepOrder[deployments[i].Step] < stepOrder[deployments[j].Step]
	})
	
	// 3. 각 배포 단계 순차 실행
	for _, deployment := range deployments {
		logger.Info("Executing deployment step",
			zap.String("step", deployment.Step),
			zap.String("deployment_id", deployment.ID.String()),
		)
		
		// 3.1 배포 상태 → InProgress
		err = s.deploymentRepo.UpdateDeploymentStatus(
			deployment.ID,
			entities.DeploymentStatusInProgress,
		)
		if err != nil {
			return fmt.Errorf("failed to update deployment status: %w", err)
		}
		
		// 3.2 로그 파일 생성
		logFile, err := os.CreateTemp("", fmt.Sprintf("deploy_%s_*.log", deployment.Step))
		if err != nil {
			return fmt.Errorf("failed to create log file: %w", err)
		}
		deployment.LogPath = logFile.Name()
		logFile.Close()
		
		// 3.3 단계별 실행 (비동기 상태 채널 사용)
		statusChan := make(chan DeploymentStatus, 10)
		go func() {
			switch deployment.Step {
			case "deploy_l1_contracts":
				s.executeL1Deployment(ctx, stackID, deployment, logFile.Name(), statusChan)
			case "deploy_aws_infrastructure":
				s.executeAWSDeployment(ctx, stackID, deployment, logFile.Name(), statusChan)
			default:
				statusChan <- DeploymentStatus{Error: fmt.Errorf("unknown step: %s", deployment.Step)}
			}
			close(statusChan)
		}()
		
		// 3.4 상태 채널에서 진행 상황 모니터링
		var stepError error
		for status := range statusChan {
			if status.Error != nil {
				stepError = status.Error
				logger.Error("Step error", zap.Error(status.Error))
			} else {
				logger.Info("Step progress",
					zap.String("message", status.Message),
					zap.Float64("progress", status.Progress),
				)
			}
		}
		
		// 3.5 컨텍스트 취소 확인
		select {
		case <-ctx.Done():
			logger.Info("Deployment cancelled", zap.String("stack_id", stackID.String()))
			s.deploymentRepo.UpdateDeploymentStatus(deployment.ID, entities.DeploymentStatusCancelled)
			return fmt.Errorf("deployment cancelled")
		default:
		}
		
		// 3.6 배포 상태 업데이트
		if stepError != nil {
			s.deploymentRepo.UpdateDeploymentStatus(deployment.ID, entities.DeploymentStatusFailed)
			return fmt.Errorf("step failed: %w", stepError)
		}
		
		s.deploymentRepo.UpdateDeploymentStatus(deployment.ID, entities.DeploymentStatusSuccess)
		
		// 3.7 스택 메타데이터 업데이트 (L1 배포 완료 후)
		if deployment.Step == "deploy_l1_contracts" {
			metadata := extractL1Metadata(logFile.Name()) // 로그에서 파싱
			s.stackRepo.UpdateMetadata(stackID, metadata)
		}
	}
	
	// 4. 통합 모듈 자동 설치 (preset 활성화된 모듈)
	s.autoInstallIntegrations(ctx, stackID)
	
	return nil
}

type DeploymentStatus struct {
	Message  string
	Progress float64
	Error    error
}
```

**L1 배포 단계** (`executeL1Deployment`):
- 계정 자금 조달 확인
- Admin 계정에 ETH 전송
- OP 컨트랙트 배포 (AddressManager, L1StandardBridge 등)
- L1 배포 config 저장
- Rollup config 생성

**AWS 배포 단계** (`executeAWSDeployment`):
- Kubernetes 클러스터 생성
- Persistent Volume 생성
- ConfigMap에서 L1 메타데이터 로드
- Docker 이미지 풀
- Pod 스펙 생성 (op-node, op-geth, op-batcher, op-proposer)
- Pod 실행 및 헬스 체크

---

## 작업 큐 시스템

### TaskManager: 메모리 기반 작업 큐

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/taskmanager/task_manager.go`

Thanos는 Redis/RabbitMQ 없이 Go의 채널과 워커 풀을 사용하는 간단한 메모리 기반 작업 큐를 구현한다.

```go
type TaskManager struct {
	taskChan    chan Task                    // 작업 채널 (버퍼 크기: 100)
	workers     int                          // 워커 수 (기본: 4)
	ctx         context.Context              // 취소용 컨텍스트
	cancel      context.CancelFunc
	activeTasks map[string]*TaskContext      // 실행 중인 작업 추적
	progress    map[string]*TaskProgress     // 진행 상황 추적
	results     map[string]interface{}       // 작업 결과 저장
	mu          sync.RWMutex                 // 맵 보호
}

type Task struct {
	ID       string
	Func     func(ctx context.Context)
	WithProgress bool
	ProgressCallback func(string, float64)
}

type TaskContext struct {
	ID     string
	Ctx    context.Context
	Cancel context.CancelFunc
	Started time.Time
}

type TaskProgress struct {
	TaskID   string
	Progress float64
	Message  string
	Status   string
	StartedAt time.Time
	UpdatedAt time.Time
}
```

### 워커 풀 패턴

```go
func (tm *TaskManager) Start() {
	for i := 0; i < tm.workers; i++ {
		go tm.worker(i)
	}
}

func (tm *TaskManager) worker(id int) {
	for {
		select {
		case <-tm.ctx.Done():
			logger.Info("Worker shutting down", zap.Int("worker_id", id))
			return
		case task := <-tm.taskChan:
			tm.executeTask(task, id)
		}
	}
}

func (tm *TaskManager) executeTask(task Task, workerID int) {
	tm.mu.Lock()
	taskCtx, cancel := context.WithCancel(tm.ctx)
	activeTask := &TaskContext{
		ID:      task.ID,
		Ctx:     taskCtx,
		Cancel:  cancel,
		Started: time.Now(),
	}
	tm.activeTasks[task.ID] = activeTask
	tm.mu.Unlock()
	
	// 진행 콜백 래퍼
	var progressFunc func(string, float64)
	if task.WithProgress {
		progressFunc = func(msg string, progress float64) {
			tm.updateProgress(task.ID, msg, progress)
		}
	}
	
	// Panic 복구 및 실행
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Task panicked",
				zap.String("task_id", task.ID),
				zap.Any("panic", r),
			)
			tm.setTaskStatus(task.ID, "failed")
		}
		tm.mu.Lock()
		delete(tm.activeTasks, task.ID)
		tm.mu.Unlock()
	}()
	
	logger.Info("Task started",
		zap.String("task_id", task.ID),
		zap.Int("worker_id", workerID),
	)
	
	task.Func(taskCtx)
	
	logger.Info("Task completed",
		zap.String("task_id", task.ID),
		zap.Duration("duration", time.Since(activeTask.Started)),
	)
}
```

### 작업 큐잉 및 취소

```go
// 기본 작업 큐잉
func (tm *TaskManager) AddTask(id string, task func(ctx context.Context)) {
	tm.taskChan <- Task{
		ID:       id,
		Func:     task,
		WithProgress: false,
	}
}

// 진행 콜백과 함께 큐잉
func (tm *TaskManager) AddTaskWithProgress(
	id string,
	task func(ctx context.Context, updateProgress func(string, float64)),
) {
	tm.taskChan <- Task{
		ID:           id,
		WithProgress: true,
		Func: func(ctx context.Context) {
			progressFunc := func(msg string, progress float64) {
				tm.updateProgress(id, msg, progress)
			}
			task(ctx, progressFunc)
		},
	}
}

// 실행 중인 작업 취소
func (tm *TaskManager) StopTask(id string) {
	tm.mu.RLock()
	taskCtx, exists := tm.activeTasks[id]
	tm.mu.RUnlock()
	
	if exists {
		logger.Info("Stopping task", zap.String("task_id", id))
		taskCtx.Cancel() // 컨텍스트 취소 신호
	}
}

// 작업 실행 여부 확인
func (tm *TaskManager) IsTaskRunning(id string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	_, exists := tm.activeTasks[id]
	return exists
}
```

### 진행 상황 추적

```go
func (tm *TaskManager) updateProgress(id string, msg string, progress float64) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if _, exists := tm.progress[id]; !exists {
		tm.progress[id] = &TaskProgress{
			TaskID:    id,
			StartedAt: time.Now(),
		}
	}
	
	tm.progress[id].Message = msg
	tm.progress[id].Progress = math.Min(progress, 100.0)
	tm.progress[id].UpdatedAt = time.Now()
}

func (tm *TaskManager) GetProgress(id string) *TaskProgress {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.progress[id]
}
```

**특징**:
- 메모리 효율적 (프로세스 종료 시 큐 손실 가능 → persistent storage 불필요)
- 단일 머신 배포에 최적화 (다중 서버 배포 시 대안 필요)
- 워커 수 설정 가능 (기본: 4개)
- 작업당 컨텍스트 기반 취소 지원
- Panic 복구로 워커 크래시 방지

---

## 데이터베이스 스키마

### StackEntity

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/domain/entities/stack.go`

```go
type StackEntity struct {
	ID                    uuid.UUID         // 주키
	Name                  string            // 스택 이름
	Type                  string            // 배포 타입 (e.g., "Thanos")
	Network               DeploymentNetwork // Mainnet, Testnet, LocalDevnet
	Status                StackStatus       // 현재 상태 (아래 참조)
	Config                datatypes.JSON    // L2 배포 설정 JSON
	                                         // {
	                                         //   "chain_id": 901,
	                                         //   "l1_rpc": "http://...",
	                                         //   "l2_config": {...},
	                                         //   "accounts": {
	                                         //     "admin": "0x...",
	                                         //     "sequencer": "0x...",
	                                         //     ...
	                                         //   }
	                                         // }
	Metadata              *StackMetadata    // L1 배포 결과 저장
	                                         // {
	                                         //   "rollup_config_url": "/path/to/rollup.json",
	                                         //   "l1_deployed_contracts": {
	                                         //     "address_manager": "0x...",
	                                         //     "bridge": "0x...",
	                                         //     ...
	                                         //   },
	                                         //   "genesis_config_url": "/path/to/genesis.json"
	                                         // }
	DeploymentPath        string            // 배포 결과물 저장 경로 (/data/stacks/{id}/)
	ImmutableConfig       datatypes.JSON    // 변경 불가능한 설정
	MainnetConfirmation   bool              // Mainnet 배포 승인 플래그
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt    // 소프트 삭제
}
```

**StackStatus enum**:

| 상태 | 의미 | 전환 |
|------|------|------|
| `Pending` | 생성되었으나 배포 대기 | → Deploying |
| `Deploying` | 배포 진행 중 | → Deployed / FailedToDeploy |
| `Deployed` | 배포 완료 | → Updating / Terminating |
| `Updating` | 업데이트 진행 중 | → Deployed / FailedToUpdate |
| `Terminating` | 종료 진행 중 | → Terminated / FailedToTerminate |
| `Terminated` | 정상 종료 | (최종 상태) |
| `Stopped` | 수동 중지 | → Deploying (재배포 시) |
| `FailedToDeploy` | 배포 실패 | → Deploying (재시도) |
| `FailedToUpdate` | 업데이트 실패 | → Updating (재시도) |
| `FailedToTerminate` | 종료 실패 | → Terminating (재시도) |

### DeploymentEntity

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/domain/entities/deployment.go`

```go
type DeploymentEntity struct {
	ID                  uuid.UUID           // 주키
	StackID             *uuid.UUID          // 외래키 (StackEntity.ID)
	Step                string              // 배포 단계명
	                                         // "deploy_l1_contracts"
	                                         // "deploy_aws_infrastructure"
	Status              DeploymentRunStatus // 실행 상태 (아래 참조)
	LogPath             string              // 로그 파일 경로 (/tmp/deploy_*.log)
	Config              datatypes.JSON      // 단계별 설정
	                                         // L1: {
	                                         //   "l1_rpc": "http://...",
	                                         //   "accounts": {...}
	                                         // }
	                                         // AWS: {
	                                         //   "region": "us-east-1",
	                                         //   "cluster_name": "thanos-xxx",
	                                         //   "namespace": "default"
	                                         // }
	CreatedAt           time.Time
	UpdatedAt           time.Time
	StartedAt           *time.Time          // 실행 시작 시각
	CompletedAt         *time.Time          // 실행 완료 시각
	PresetID            *uuid.UUID          // 사용된 preset ID (선택적)
	SeedDerivedAccounts datatypes.JSON      // 시드에서 유도한 계정 주소
	                                         // {
	                                         //   "admin": "0x...",
	                                         //   "sequencer": "0x...",
	                                         //   ...
	                                         // }
	FundingStatus       string              // 계정 자금 상태
	                                         // "unfunded", "funding", "funded"
	ModuleConfigs       datatypes.JSON      // 통합 모듈 설정
}
```

**DeploymentRunStatus enum**:

| 상태 | 의미 |
|------|------|
| `Pending` | 대기 중 |
| `InProgress` | 실행 중 |
| `Success` | 성공 |
| `Failed` | 실패 |
| `Stopped` | 중단됨 |
| `Cancelled` | 취소됨 |

### IntegrationEntity

**파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/domain/entities/integration.go`

```go
type IntegrationEntity struct {
	ID        uuid.UUID         // 주키
	StackID   uuid.UUID         // 외래키 (StackEntity.ID)
	Type      string            // 통합 타입
	                              // "Bridge", "BlockExplorer", "Monitoring",
	                              // "CrossChainBridge", "UptimeService", "DRB"
	Status    DeploymentStatus  // 설치 상태
	Config    datatypes.JSON    // 타입별 설정
	           // Bridge: {
	           //   "enabled": true,
	           //   "contracts": {...}
	           // }
	           // BlockExplorer: {
	           //   "enabled": true,
	           //   "api_url": "http://...",
	           //   "db_url": "..."
	           // }
	           // Monitoring: {
	           //   "prometheus_url": "http://...",
	           //   "grafana_url": "http://...",
	           //   "alert_webhook": "..."
	           // }
	CreatedAt time.Time
	UpdatedAt time.Time
	StartedAt *time.Time        // 설치 시작 시각
	CompletedAt *time.Time      // 설치 완료 시각
	LogPath   string            // 로그 파일 경로
}
```

**DeploymentStatus enum**:

| 상태 | 의미 |
|------|------|
| `Pending` | 대기 중 |
| `InProgress` | 설치/실행 중 |
| `Completed` | 완료 |
| `Failed` | 실패 |
| `Stopped` | 중단됨 |
| `Terminating` | 제거 중 |
| `Terminated` | 제거됨 |
| `Cancelled` | 취소됨 |
| `AwaitingConfig` | 설정 대기 중 |

### 데이터베이스 관계도

```
┌─────────────────┐
│  StackEntity    │
├─────────────────┤
│ ID (PK)         │ 1
│ Name            │
│ Status          │
│ Config          │
│ Metadata        │
│ CreatedAt       │
└─────────────────┘
        │
        │ 1:N
        ├─────────────────────────────┬──────────────────────────────┐
        │                             │                              │
┌───────────────────┐    ┌──────────────────────┐    ┌──────────────────────┐
│ DeploymentEntity  │    │ IntegrationEntity    │    │ LogEntity (Optional) │
├───────────────────┤    ├──────────────────────┤    ├──────────────────────┤
│ ID (PK)           │    │ ID (PK)              │    │ ID (PK)              │
│ StackID (FK)      │ N  │ StackID (FK)         │ N  │ DeploymentID (FK)    │
│ Step              │    │ Type                 │    │ LogPath              │
│ Status            │    │ Status               │    │ CreatedAt            │
│ LogPath           │    │ Config               │    └──────────────────────┘
│ Config            │    │ CreatedAt            │
│ StartedAt         │    │ CompletedAt          │
│ CompletedAt       │    └──────────────────────┘
└───────────────────┘
```

---

## 상태 머신 및 에러 처리

### 전체 상태 흐름

```
┌─────────────┐
│   Pending   │ ← 스택 생성 직후
└──────┬──────┘
       │
       └──→ [작업 큐잉]
            TaskManager.AddTask(deploy_task)
       │
       ▼
┌─────────────────────────────────────────┐
│  Worker executes deploy() task          │
│  ├─ deploy_l1_contracts (step)          │
│  │  ├─ Fund admin account               │
│  │  ├─ Deploy contracts on L1           │
│  │  └─ Save deployment config           │
│  │                                       │
│  └─ deploy_aws_infrastructure (step)    │
│     ├─ Create K8s cluster               │
│     ├─ Deploy pods (op-node, op-geth)   │
│     └─ Health check                     │
└──────────────────────────────────────────┘
       │
       ├─ 성공  ▼
       │    ┌─────────────┐
       │    │  Deployed   │ ← 모든 단계 완료
       │    └─────────────┘
       │
       └─ 실패  ▼
            ┌──────────────────┐
            │ FailedToDeploy   │ ← 재시도 가능
            └──────────────────┘
            
추가 전환:
Deployed → Updating → Deployed / FailedToUpdate
Deployed → Terminating → Terminated / FailedToTerminate
[Any State] → StopTask() → Stopped
[Any State] → CancelIntegration() → Cancelled (for IntegrationEntity)
```

### 배포 단계 상태 전환

**L1 배포 단계**:
```
┌─────────┐
│Pending  │ (DeploymentEntity 생성)
└────┬────┘
     │ executeDeployments() 시작
     ▼
┌──────────────┐
│ InProgress   │
└────┬─────────┘
     │ L1 배포 실행
     │
     ├─ 성공  ▼
     │    ┌─────────┐
     │    │ Success │ → AWS 단계로 진행
     │    └─────────┘
     │
     └─ 실패  ▼
          ┌────────┐
          │ Failed │ → 스택 상태 = FailedToDeploy
          └────────┘
```

**AWS 배포 단계**:
```
┌─────────┐
│Pending  │ (L1 완료 후 시작)
└────┬────┘
     │
     ▼
┌──────────────┐
│ InProgress   │
└────┬─────────┘
     │
     ├─ 성공  ▼
     │    ┌─────────┐
     │    │ Success │ → 스택 상태 = Deployed
     │    └─────────┘
     │
     └─ 실패  ▼
          ┌────────┐
          │ Failed │ → 스택 상태 = FailedToDeploy
          └────────┘
```

---

## 에러 처리 및 시나리오

### 시나리오 1: Seed Phrase 검증 실패

**상황**: 사용자가 유효하지 않은 시드 구문을 제공 (11 단어, 부정확한 체크섬)

**코드 경로**: `PresetDeploy` → `CreateThanosStackFromPreset` → `deriveSeedAccounts`

```go
// preset_deploy.go의 유효성 검사
seedPhrase := req.SeedPhrase
words := strings.Fields(seedPhrase)

// 검증 1: 단어 개수 확인
if len(words) != 12 && len(words) != 24 {
	return nil, fmt.Errorf("seed phrase must have 12 or 24 words, got %d", len(words))
}

// 검증 2: BIP39 체크섬 검증
_, err := bip39.NewMnemonic(seedPhrase)
if err != nil {
	return nil, fmt.Errorf("invalid seed phrase: %w", err)
}
```

**응답**:
```json
{
  "success": false,
  "message": "invalid seed phrase: checksum mismatch",
  "error": "seed phrase checksum validation failed"
}
```

**상태 변화**: 요청이 검증 단계에서 거부되어 스택 엔티티 생성 전에 반환됨. 데이터베이스 변경 없음.

---

### 시나리오 2: 트랜잭션 커밋 실패

**상황**: `CreateThanosStack` 중 데이터베이스 연결 끊김

**코드 경로**: `CreateThanosStack` → 트랜잭션 커밋

```go
// stack_lifecycle.go line ~262
if err := tx.Commit().Error; err != nil {
	logger.Error("Failed to commit transaction",
		zap.String("stack_id", stackEntity.ID.String()),
		zap.Error(err),
	)
	return nil, fmt.Errorf("failed to commit transaction: %w", err)
}
```

**응답**:
```json
{
  "success": false,
  "error": "failed to commit transaction: connection refused"
}
```

**상태 변화**: 
- 부분적으로 생성된 엔티티는 없음 (트랜잭션 롤백으로 보호)
- 클라이언트는 재시도 안전 (멱등성)

---

### 시나리오 3: L1 배포 컨트랙트 배포 실패

**상황**: Ethers.js가 컨트랙트 배포 트랜잭션 생성 실패 (가스 부족)

**코드 경로**: `executeDeployments` → `executeL1Deployment` → contract deployment

```go
// deployment.go의 L1 배포 실행
func (s *ThanosStackDeploymentService) executeL1Deployment(
	ctx context.Context,
	stackID uuid.UUID,
	deployment *entities.DeploymentEntity,
	logPath string,
	statusChan chan DeploymentStatus,
) {
	// ... 설정 파싱 ...
	
	// 계정 자금 확인
	balance, err := s.ethClient.GetBalance(adminAccount)
	if err != nil || balance < requiredGas {
		statusChan <- DeploymentStatus{
			Error: fmt.Errorf("insufficient gas: need %v, have %v", requiredGas, balance),
		}
		return
	}
	
	// 컨트랙트 배포
	_, err = s.deployContracts(ctx, deployment.Config)
	if err != nil {
		statusChan <- DeploymentStatus{
			Error: fmt.Errorf("contract deployment failed: %w", err),
		}
		return
	}
}
```

**응답**:
```json
{
  "success": false,
  "message": "Deployment failed",
  "error": "contract deployment failed: insufficient gas"
}
```

**상태 변화**:
- `StackEntity.Status` = `FailedToDeploy`
- `DeploymentEntity[0].Status` = `Failed` (L1 단계)
- `DeploymentEntity[1].Status` = `Pending` (AWS 단계는 실행 안 됨)
- 로그 파일은 유지되어 디버깅 가능

---

### 시나리오 4: 작업 컨텍스트 취소

**상황**: 사용자가 배포 중 `StopTask()` 호출

**코드 경로**: `TaskManager.StopTask()` → `taskCtx.Cancel()` → context 감지

```go
// deployment.go의 executeDeployments에서
for _, deployment := range deployments {
	// ... 배포 실행 ...
	
	// 3.5 컨텍스트 취소 확인
	select {
	case <-ctx.Done():
		logger.Info("Deployment cancelled by user",
			zap.String("stack_id", stackID.String()),
		)
		s.deploymentRepo.UpdateDeploymentStatus(
			deployment.ID,
			entities.DeploymentStatusCancelled,
		)
		return fmt.Errorf("deployment cancelled by user")
	default:
	}
}
```

**응답**: (비동기) 작업 취소 완료
```json
{
  "success": true,
  "message": "Task cancelled"
}
```

**상태 변화**:
- `StackEntity.Status` = `Stopped` (또는 현재 단계 상태 유지)
- 현재 실행 중이던 `DeploymentEntity` = `Cancelled`
- 미실행 단계 = `Pending` (유지)

---

### 시나리오 5: AWS 클러스터 생성 실패

**상황**: Kubernetes 클러스터 생성 타임아웃 (20분 이상)

**코드 경로**: `executeAWSDeployment` → EKS cluster creation

```go
// deployment.go
func (s *ThanosStackDeploymentService) executeAWSDeployment(
	ctx context.Context,
	stackID uuid.UUID,
	deployment *entities.DeploymentEntity,
	logPath string,
	statusChan chan DeploymentStatus,
) {
	statusChan <- DeploymentStatus{
		Message: "Creating EKS cluster",
		Progress: 10,
	}
	
	// 클러스터 생성 (타임아웃 설정)
	createCtx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()
	
	cluster, err := s.awsClient.CreateCluster(createCtx, &eks.CreateClusterInput{
		Name: aws.String(clusterName),
		// ... 설정 ...
	})
	
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			statusChan <- DeploymentStatus{
				Error: fmt.Errorf("cluster creation timeout: %w", err),
			}
		} else {
			statusChan <- DeploymentStatus{
				Error: fmt.Errorf("cluster creation failed: %w", err),
			}
		}
		return
	}
}
```

**로그 출력** (logPath에 기록):
```
2024-04-16T12:30:15Z INFO  Starting AWS infrastructure deployment
2024-04-16T12:30:16Z INFO  Creating EKS cluster (progress=10%)
2024-04-16T12:50:15Z ERROR Cluster creation timeout: context deadline exceeded
2024-04-16T12:50:15Z INFO  Deployment step cancelled
```

**응답**:
```json
{
  "success": false,
  "message": "Deployment failed",
  "error": "cluster creation timeout: context deadline exceeded"
}
```

**상태 변화**:
- `StackEntity.Status` = `FailedToDeploy`
- `DeploymentEntity[0].Status` = `Success` (L1은 완료)
- `DeploymentEntity[1].Status` = `Failed` (AWS 실패)

---

### 시나리오 6: 로컬 스택 포트 충돌

**상황**: 로컬 배포 중 이미 실행 중인 스택이 있음

**코드 경로**: `CreateThanosStack` → `checkNoActiveLocalStack()`

```go
// stack_lifecycle.go line ~105
if request.Network == entities.LocalDevnet {
	if err := s.checkNoActiveLocalStack(); err != nil {
		return nil, fmt.Errorf("local stack already active: %w", err)
	}
}

func (s *ThanosStackDeploymentService) checkNoActiveLocalStack() error {
	// 로컬 포트 확인 (예: 8545, 8546, 9545 등)
	activeStacks, err := s.stackRepo.GetStacksByNetworkAndStatus(
		entities.LocalDevnet,
		entities.StackStatusDeploying,
		entities.StackStatusDeployed,
	)
	if err != nil {
		return err
	}
	
	if len(activeStacks) > 0 {
		return fmt.Errorf("local stack already running: %s", activeStacks[0].Name)
	}
	
	return nil
}
```

**응답**:
```json
{
  "success": false,
  "error": "local stack already active: dev-stack-1"
}
```

**상태 변화**: 스택 엔티티 생성 전에 거부됨. 데이터베이스 변경 없음.

---

### 시나리오 7: 통합 모듈 자동 설치 부분 실패

**상황**: BlockExplorer 설치 중 데이터베이스 연결 실패 (Bridge는 이미 설치됨)

**코드 경로**: `executeDeployments` → `autoInstallIntegrations`

```go
// deployment.go
func (s *ThanosStackDeploymentService) autoInstallIntegrations(
	ctx context.Context,
	stackID uuid.UUID,
) error {
	stack, _ := s.stackRepo.GetStackByID(stackID.String())
	
	integrations := s.integrationRepo.GetIntegrationsByStackID(stackID)
	for _, integration := range integrations {
		if integration.Status != entities.DeploymentStatusPending {
			continue
		}
		
		switch integration.Type {
		case "Bridge":
			s.integrationMgr.InstallBridge(ctx, stackID.String())
		case "BlockExplorer":
			err := s.integrationMgr.InstallBlockExplorer(ctx, stackID, blockExplorerReq)
			if err != nil {
				// 부분 실패: Bridge는 설치됨, BlockExplorer는 실패
				logger.Error("Failed to install BlockExplorer",
					zap.String("stack_id", stackID.String()),
					zap.Error(err),
				)
				// 더 진행하거나 롤백할지 결정
				return fmt.Errorf("integration installation failed: %w", err)
			}
		}
	}
	return nil
}
```

**상태 변화**:
- `StackEntity.Status` = `Deployed` (기본 배포는 완료)
- `IntegrationEntity[Bridge].Status` = `Completed`
- `IntegrationEntity[BlockExplorer].Status` = `Failed`
- `IntegrationEntity[Monitoring].Status` = `Pending` (미실행)

**사용자 조치**: 실패한 BlockExplorer를 다시 설치하기 위해 `RetryIntegration()` API 호출

---

## 호출 시퀀스

### 전체 흐름

```
┌──────────────────────────────────────────────────────────────────┐
│ 1. HTTP 요청                                                      │
│    POST /api/v1/stacks/deploy                                    │
│    Body: {preset_id, stack_name, seed_phrase, ...}              │
└──────┬───────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────┐
│ 2. PresetHandler.PresetDeploy() [presets.go:78-104]            │
│    - JSON 파싱                                                    │
│    - 필드 검증                                                    │
│    - 서비스 호출 (에러 시 응답 반환)                             │
└──────┬───────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────┐
│ 3. CreateThanosStackFromPreset() [preset_deploy.go:18-167]      │
│    - Preset 로드                                                  │
│    - 시드 구문 검증 (BIP39)                                      │
│    - 계정 유도 (m/44'/60'/0'/0/*)                               │
│    - 설정 병합 (preset + 사용자 재정의)                         │
│    - 시드 구문 메모리 제거                                       │
│    - CreateThanosStack() 호출                                    │
└──────┬───────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────┐
│ 4. CreateThanosStack() [stack_lifecycle.go:101-273]             │
│                                                                   │
│    4.1 로컬 배포 사전 검사                                        │
│        checkNoActiveLocalStack()                                 │
│                                                                   │
│    4.2 트랜잭션 시작                                             │
│        tx := deploymentRepo.BeginTx()                            │
│                                                                   │
│    4.3 StackEntity 생성                                         │
│        stackRepo.CreateStackByTx(tx, stackEntity)              │
│        Status: Pending                                           │
│                                                                   │
│    4.4 DeploymentEntity 생성 (L1)                              │
│        deploymentRepo.CreateDeployment(tx, l1Deployment)       │
│        Step: deploy_l1_contracts                                │
│        Status: Pending                                           │
│                                                                   │
│    4.5 DeploymentEntity 생성 (AWS)                             │
│        deploymentRepo.CreateDeployment(tx, awsDeployment)      │
│        Step: deploy_aws_infrastructure                          │
│        Status: Pending                                           │
│                                                                   │
│    4.6 IntegrationEntity 생성                                  │
│        for each enabled module (Bridge, BlockExplorer, ...):   │
│          integrationRepo.CreateIntegration(tx, integration)     │
│                                                                   │
│    4.7 트랜잭션 커밋                                             │
│        tx.Commit()                                              │
│                                                                   │
│    4.8 배포 작업 큐잉                                            │
│        taskID := "deploy_{stackID}"                             │
│        taskManager.AddTask(taskID, deploy_callback)           │
│        (이 시점에서 함수 반환)                                  │
│                                                                   │
│    4.9 응답 반환                                                 │
│        {stack_id, deployment_id, status: queued}               │
└──────┬───────────────────────────────────────────────────────────┘
       │
       ▼ (비동기 - 즉시 반환)
┌──────────────────────────────────────────────────────────────────┐
│ 5. TaskManager 워커 처리                                         │
│                                                                   │
│    5.1 작업 채널에서 deploy 태스크 획득                         │
│        Worker goroutine: worker()                               │
│                                                                   │
│    5.2 deploy() 함수 실행 [deployment.go:31-437]              │
│        - 스택 상태 → Deploying                                 │
│        - executeDeployments() 호출                             │
│        - 에러 시: 상태 → FailedToDeploy                       │
│        - 성공 시: 상태 → Deployed                             │
└──────┬───────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────┐
│ 6. executeDeployments() [deployment.go:439-650+]                │
│                                                                   │
│    ┌─ 각 배포 단계 순차 실행 ─┐                                 │
│    │                        │                                   │
│    │ Step 1: deploy_l1_contracts                               │
│    │ ├─ DeploymentEntity[0].Status → InProgress              │
│    │ ├─ 로그 파일 생성                                         │
│    │ ├─ executeL1Deployment() (비동기)                       │
│    │ │  ├─ 계정 자금 확인                                      │
│    │ │  ├─ Admin 계정에 ETH 전송                             │
│    │ │  ├─ 컨트랙트 배포 (AddressManager, Bridge)           │
│    │ │  ├─ L1 config 저장                                    │
│    │ │  └─ Rollup config 생성                               │
│    │ ├─ statusChan 모니터링 (진행률 수신)                   │
│    │ ├─ ctx.Done() 확인 (취소 신호)                         │
│    │ └─ DeploymentEntity[0].Status → Success                │
│    │
│    │ Step 2: deploy_aws_infrastructure                        │
│    │ ├─ DeploymentEntity[1].Status → InProgress             │
│    │ ├─ 로그 파일 생성                                        │
│    │ ├─ executeAWSDeployment() (비동기)                     │
│    │ │  ├─ K8s 클러스터 생성                                 │
│    │ │  ├─ Persistent Volume 생성                           │
│    │ │  ├─ ConfigMap에서 L1 메타데이터 로드                 │
│    │ │  ├─ Pod 스펙 생성 (op-node, op-geth, ...)           │
│    │ │  ├─ Pod 실행                                         │
│    │ │  └─ 헬스 체크                                        │
│    │ ├─ statusChan 모니터링                                 │
│    │ ├─ ctx.Done() 확인                                      │
│    │ └─ DeploymentEntity[1].Status → Success               │
│    │
│    └─────────────────────┘
│
│    6.1 autoInstallIntegrations() 호출                        │
│        for each IntegrationEntity with Status = Pending:    │
│          - InstallBridge()                                  │
│          - InstallBlockExplorer()                           │
│          - InstallMonitoring()                              │
│          - 각각 비동기 실행                                  │
│
└──────┬───────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────┐
│ 7. 최종 상태 업데이트                                            │
│                                                                   │
│    모든 단계 성공 시:                                            │
│    └─ StackEntity.Status → Deployed                           │
│                                                                   │
│    실패한 단계 있을 시:                                         │
│    └─ StackEntity.Status → FailedToDeploy                    │
│       (해당 DeploymentEntity.Status → Failed)                 │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

### 시간 추정

| 단계 | 예상 시간 |
|------|----------|
| HTTP 요청 → 응답 | < 100ms |
| 트랜잭션 커밋까지 | 200-500ms |
| TaskManager 큐잉 | < 10ms |
| **총 동기 시간** | **300-600ms** |
| L1 배포 (실제) | 3-10분 (네트워크, 블록 시간) |
| AWS 인프라 생성 | 15-25분 (EKS 클러스터) |
| **전체 배포 시간** | **20-40분** |

---

## 결론

Phase 3는 사용자의 배포 요청을 안전하게 수신하고, 원자적 트랜잭션으로 데이터베이스에 저장한 후, 메모리 기반 작업 큐를 통해 비동기적으로 오케스트레이션한다. 핵심 설계 결정은:

1. **트랜잭션 기반 원자성**: 부분적 생성 상태 방지
2. **메모리 기반 큐**: 단일 머신 배포에 최적화, Redis 불필요
3. **단계별 순차 실행**: L1 → AWS 순서 보장
4. **컨텍스트 기반 취소**: Graceful shutdown 지원
5. **상세한 로깅**: 배포 경로 추적 및 디버깅 가능

이러한 설계는 단순성과 신뢰성의 균형을 맞추며, 프로덕션 환경의 Thanos 배포를 안정적으로 지원한다.
