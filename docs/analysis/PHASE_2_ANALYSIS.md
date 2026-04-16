# Phase 2: Web UI 배포 요청 (Frontend → Backend)

## 개요

Electron 앱의 Next.js 웹 UI에서 사용자가 배포 버튼을 클릭한 후, HTTP POST 요청이 trh-backend의 엔드포인트에 도달하는 전체 과정을 분석합니다. 배포는 두 가지 방식으로 시작됩니다:
1. **Preset 기반 배포** (`POST /api/v1/stacks/thanos/preset-deploy`) - 사전 정의된 템플릿 사용
2. **직접 배포** (`POST /api/v1/stacks/thanos`) - 직접 파라미터 입력

본 분석은 현재 활발히 사용되는 Preset 기반 배포 플로우를 중심으로 진행합니다.

---

## 진입점 (Frontend - Preset 마법사)

### 페이지 컴포넌트
- **파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/app/rollup/create/page.tsx`
- **컴포넌트**: `PresetWizardContent()`
- **라인**: 18-139
- **역할**: 3단계 Preset 마법사 UI를 렌더링하고 배포 트리거 제어

### 배포 트리거
- **파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/app/rollup/create/page.tsx`
- **라인**: 121 (버튼)
  ```tsx
  <Button
    onClick={goToNextStep}
    disabled={currentStep === 1 && !selectedPresetId}
    className="flex items-center gap-2"
  >
    {isLastStep ? "Deploy Rollup" : "Next"}
  </Button>
  ```
- **버튼 텍스트**: `"Deploy Rollup"` (마지막 단계에서만 표시)
- **핸들러 함수**: `goToNextStep()` (usePresetWizard 훅에서 제공)

---

## 마법사 상태 관리

### usePresetWizard 훅
- **파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/features/rollup/hooks/usePresetWizard.ts`
- **라인**: 43-177
- **핵심 상태**:
  - `currentStep`: 현재 단계 (1-3)
  - `selectedPresetId`: 선택된 Preset ID
  - `presetDetail`: 선택된 Preset 상세 정보
  - `overrides`: 사용자가 수정한 파라미터

### 3단계 흐름
1. **Step 1**: Preset 선택 (general, defi, gaming, full)
2. **Step 2**: 기본 정보 입력 (Chain name, Network, L1 RPC, Seed phrase, AWS 자격증)
3. **Step 3**: 설정 검토 및 배포

---

## 배포 함수 호출 체인

### Step 3: goToNextStep() → handleDeploy()
**파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/features/rollup/hooks/usePresetWizard.ts`

**라인**: 87-126 (handleDeploy 함수)

```typescript
const handleDeploy = useCallback(async () => {
  const basicInfo = form.getValues("presetBasicInfo");

  try {
    toast.loading("Initiating deployment...", { id: "preset-deploy" });

    const { startPresetDeployment } = await import("../services/presetService");
    const result = await startPresetDeployment({
      presetId: selectedPresetId!,
      chainName: basicInfo.chainName,
      network: basicInfo.network,
      infraProvider: basicInfo.infraProvider,
      seedPhrase: basicInfo.seedPhrase.join(" "),
      feeToken: basicInfo.feeToken,
      awsAccessKey: basicInfo.awsAccessKey ?? "",
      awsSecretKey: basicInfo.awsSecretKey ?? "",
      awsRegion: basicInfo.awsRegion ?? "",
      l1RpcUrl: basicInfo.l1RpcUrl,
      l1BeaconUrl: basicInfo.l1BeaconUrl,
      reuseDeployment: basicInfo.network === "Mainnet" ? (basicInfo.reuseDeployment ?? true) : undefined,
      overrides: overrides.length > 0 ? overrides : undefined,
    });

    toast.success("Deployment initiated!", { id: "preset-deploy" });
    setPendingDeploymentId(result.deploymentId);
    resetState();
    router.push("/rollup");
  } catch (error) {
    // 에러 처리
  }
}, [form, selectedPresetId, overrides, setPendingDeploymentId, router]);
```

### 주요 특징
- **Toast 알림**: 로딩 중 → 성공/실패 피드백
- **Seed phrase 변환**: 배열 → 공백으로 구분된 문자열 (`basicInfo.seedPhrase.join(" ")`)
- **네비게이션**: 배포 완료 후 `/rollup` 페이지로 이동

---

## HTTP 요청: startPresetDeployment()

### 함수 정의
- **파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/features/rollup/services/presetService.ts`
- **라인**: 45-53
- **함수명**: `startPresetDeployment()`

```typescript
export const startPresetDeployment = async (
  request: DeployWithPresetRequest
): Promise<{ deploymentId: string }> => {
  const response = await apiPost<{ deploymentId: string }>(
    "stacks/thanos/preset-deploy",
    request
  );
  return response.data;
};
```

### 엔드포인트
- **상대 경로**: `stacks/thanos/preset-deploy`
- **절대 경로**: `/api/proxy/stacks/thanos/preset-deploy`
- **기본 URL**: `/api/proxy/` (Next.js 프록시를 통한 후처리)
- **HTTP 메서드**: `POST`

### 요청 페이로드 타입
**파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/features/rollup/schemas/preset.ts`

**라인**: 140-155

```typescript
export const deployWithPresetRequestSchema = z.object({
  presetId: z.string(),
  chainName: z.string().regex(/^[a-z0-9-]{3,32}$/, "Must be 3-32 lowercase alphanumeric characters or hyphens"),
  network: z.enum(["Testnet", "Mainnet"]),
  seedPhrase: z.string().min(1),
  infraProvider: z.enum(["aws", "local"]),
  awsAccessKey: z.string().optional().default(""),
  awsSecretKey: z.string().optional().default(""),
  awsRegion: z.string().optional().default(""),
  l1RpcUrl: z.string().url(),
  l1BeaconUrl: z.string().url(),
  feeToken: feeTokenSchema,
  reuseDeployment: z.boolean().optional(),
  overrides: z.array(presetFieldOverrideSchema).optional(),
});
export type DeployWithPresetRequest = z.infer<typeof deployWithPresetRequestSchema>;
```

### 실제 인프라 제공자 열거형
**파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/features/rollup/schemas/preset.ts`

**라인**: 145

```typescript
infraProvider: z.enum(["aws", "local"]),
```

**중요 주의**: 원본 명세에서는 "aws/gcp/azure or local"을 언급했지만, 실제 구현은 **"aws"와 "local"만 지원**합니다. GCP와 Azure는 구현되지 않았습니다.

### 예시 JSON 페이로드
```json
{
  "presetId": "gaming",
  "chainName": "my-chain",
  "network": "Testnet",
  "seedPhrase": "word1 word2 word3 ... word12",
  "infraProvider": "aws",
  "awsAccessKey": "AKIAIOSFODNN7EXAMPLE",
  "awsSecretKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
  "awsRegion": "us-east-1",
  "l1RpcUrl": "https://sepolia.infura.io/v3/PROJECT_ID",
  "l1BeaconUrl": "https://www.beaconcha.in/api/v1",
  "feeToken": "TON",
  "reuseDeployment": true,
  "overrides": [
    {
      "field": "l2BlockTime",
      "value": 4
    }
  ]
}
```

### 실제 네트워크 열거형
**파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/features/rollup/schemas/preset.ts`

**라인**: 143

```typescript
network: z.enum(["Testnet", "Mainnet"]),
```

**중요 주의**: 원본 명세에서는 "sepolia" 지원을 언급했지만, 실제 구현은 **Testnet과 Mainnet만 지원**합니다. sepolia는 구현되지 않았으며 일반 Testnet 환경을 사용합니다.

### 실제 인프라 제공자 열거형
**파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/features/rollup/schemas/preset.ts`

**라인**: 145

```typescript
infraProvider: z.enum(["aws", "local"]),
```

**중요 주의**: 원본 명세에서는 "aws/gcp/azure or local"을 언급했지만, 실제 구현은 **"aws"와 "local"만 지원**합니다. GCP와 Azure는 구현되지 않았습니다.

---

## 조건부 필드 요구사항

폼에는 다른 필드의 값에 따라 필수 또는 동작이 달라지는 여러 필드가 있습니다:

### 1. AWS 제공자 기반 필드
- **awsAccessKey, awsSecretKey, awsRegion**: `infraProvider = "aws"`일 때만 필수입니다.
- `infraProvider = "local"`일 때, 이 필드들은 무시됩니다.
- **Backend 검증** (thanos.go, 라인 563-564):
  ```go
  if r.AwsAccessKey == "" || r.AwsSecretKey == "" || r.AwsRegion == "" {
    return fmt.Errorf("awsAccessKey, awsSecretKey, and awsRegion are required for aws provider")
  }
  ```
- **Frontend 검증** (create-rollup.ts, 라인 214-224):
  ```typescript
  if (data.infraProvider === "aws") {
    if (!data.awsAccessKey) {
      ctx.addIssue({ ... message: "AWS Access Key is required" ... });
    }
    if (!data.awsSecretKey) {
      ctx.addIssue({ ... message: "AWS Secret Key is required" ... });
    }
    if (!data.awsRegion) {
      ctx.addIssue({ ... message: "AWS region is required" ... });
    }
  }
  ```

### 2. 네트워크 + 제공자 충돌
- **네트워크 = "Mainnet" + 제공자 = "local"은 무효**합니다.
- Backend가 이 조합을 거부합니다.
- **Backend 검증** (thanos.go, 라인 570-572):
  ```go
  case "local":
    if r.Network == "Mainnet" {
      return fmt.Errorf("local deployment is not supported for Mainnet")
    }
  ```
- **Frontend 검증** (create-rollup.ts, 라인 226-228):
  ```typescript
  if (data.infraProvider === "local" && data.network === "Mainnet") {
    ctx.addIssue({ 
      code: z.ZodIssueCode.custom, 
      message: "Local deployment is not supported for Mainnet", 
      path: ["network"] 
    });
  }
  ```

### 3. L1 URL 필드
- **l1RpcUrl, l1BeaconUrl**: 항상 필수입니다.
- 두 필드 모두 유효한 URL 형식이어야 합니다.
- **Frontend 검증** (preset.ts, 라인 149-150):
  ```typescript
  l1RpcUrl: z.string().url(),
  l1BeaconUrl: z.string().url(),
  ```
- **Backend DTO** (thanos.go, 라인 552-553): 바인딩 required 태그 있음

### 4. reuseDeployment 필드
- **Mainnet에만 적용**됩니다.
- **Frontend 로직** (usePresetWizard.ts, 라인 83):
  ```typescript
  reuseDeployment: basicInfo.network === "Mainnet" ? (basicInfo.reuseDeployment ?? true) : undefined,
  ```
- Testnet일 때는 값이 `undefined`로 전송되어 무시됩니다.

---

## HTTP 클라이언트 구성

### API 클라이언트
- **파일**: `/Users/theo/workspace_tokamak/trh-platform-ui/src/lib/api.ts`
- **라인**: 1-195
- **HTTP 라이브러리**: axios

### 요청 설정

| 항목 | 값 |
|------|-----|
| **기본 URL** | `/api/proxy/` |
| **타임아웃** | 10초 (기본값) |
| **Content-Type** | `application/json` |
| **Authorization 헤더** | `Bearer {accessToken}` (localStorage에서 로드) |
| **ngrok 헤더** | `ngrok-skip-browser-warning: 69420` |

### Request Interceptor
- **라인**: 26-41
- **기능**: 
  - localStorage에서 `accessToken` 추출
  - Authorization 헤더에 `Bearer {token}` 형식으로 추가
  - 토큰이 없으면 헤더 미추가

```typescript
apiClient.interceptors.request.use(
  (config) => {
    const token = typeof window !== "undefined"
      ? localStorage.getItem("accessToken")
      : null;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);
```

### Response Interceptor
- **라인**: 44-64
- **기능**:
  - 401 에러 시: localStorage에서 accessToken 제거 (로그아웃)
  - 403 에러: 접근 거부 로그 기록
  - 기타 에러: 그대로 전달

### API POST 함수
- **라인**: 121-136

```typescript
export const apiPost = async <T = any>(
  url: string,
  data?: any,
  config?: AxiosRequestConfig
): Promise<ApiResponse<T>> => {
  try {
    const response: AxiosResponse<ApiResponse<T>> = await apiClient.post(
      url,
      data,
      config
    );
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError);
  }
};
```

---

## 응답 처리

### 성공 응답 (HTTP 200)

**타입**: `ApiResponse<{ deploymentId: string }>`

```typescript
export interface ApiResponse<T = any> {
  data: T;
  message?: string;
  success: boolean;
}
```

**예시 응답**:
```json
{
  "success": true,
  "message": "ok",
  "data": {
    "deploymentId": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### 응답 처리 (usePresetWizard)
- **라인**: 113-115

```typescript
toast.success("Deployment initiated!", { id: "preset-deploy" });
setPendingDeploymentId(result.deploymentId);
resetState();
router.push("/rollup");
```

### 실패 응답 (HTTP 4xx, 5xx)

**에러 처리**: handleApiError() 함수
- **라인**: 81-102

```typescript
export const handleApiError = (error: AxiosError): ApiError => {
  if (error.response) {
    // 서버 응답 에러
    return {
      message: (error.response.data as any)?.message || "An error occurred",
      status: error.response.status,
      errors: (error.response.data as any)?.errors,
    };
  } else if (error.request) {
    // 요청 했지만 응답 없음
    return {
      message: "No response from server",
      status: 0,
    };
  } else {
    // 기타 에러
    return {
      message: error.message || "An unexpected error occurred",
      status: 0,
    };
  }
};
```

**에러 UI 피드백** (usePresetWizard, 라인 117-124):
```typescript
catch (error) {
  const message = error instanceof Error
    ? error.message
    : typeof error === "object" && error !== null && "message" in error
      ? String((error as { message: unknown }).message)
      : "Failed to initiate deployment";
  toast.error(message, { id: "preset-deploy" });
}
```

---

## 백엔드 핸들러

### Preset 배포 핸들러
- **파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/api/handlers/thanos/presets.go`
- **함수**: `PresetDeploy()`
- **라인**: 78-104
- **HTTP 핸들러 시그니처**: `func (h *ThanosDeploymentHandler) PresetDeploy(c *gin.Context)`

```go
// @Summary		Deploy with Preset
// @Description	Create a Thanos stack from a preset definition (Admin only)
// @Tags			Thanos Presets
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dtos.PresetDeployRequest	true	"Preset Deploy Request"
// @Success		200		{object}	entities.Response
// @Failure		400		{object}	map[string]interface{}
// @Failure		401		{object}	map[string]interface{}
// @Failure		403		{object}	map[string]interface{}
// @Router			/stacks/thanos/preset-deploy [post]
func (h *ThanosDeploymentHandler) PresetDeploy(c *gin.Context) {
  var request dtos.PresetDeployRequest
  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, &entities.Response{
      Status:  http.StatusBadRequest,
      Message: err.Error(),
    })
    return
  }

  // Validate preset ID
  presetSvc := presets.NewService()
  if _, err := presetSvc.GetByID(request.PresetID); err != nil {
    c.JSON(http.StatusBadRequest, &entities.Response{
      Status:  http.StatusBadRequest,
      Message: err.Error(),
    })
    return
  }

  response, err := h.ThanosDeploymentService.CreateThanosStackFromPreset(c.Request.Context(), request)
  if err != nil {
    logger.Error("failed to deploy preset stack", zap.Error(err))
  }

  c.JSON(int(response.Status), response)
}
```

### 직접 배포 핸들러 (비교 참고)
- **파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/api/handlers/thanos/deployment.go`
- **함수**: `Deploy()`
- **라인**: 32-120
- **엔드포인트**: `POST /stacks/thanos`

---

## 요청 DTO (Backend)

### PresetDeployRequest
- **파일**: `/Users/theo/workspace_tokamak/trh-backend/pkg/api/dtos/thanos.go`
- **라인**: 539-577

```go
type PresetDeployRequest struct {
  PresetID      string                     `json:"presetId"      binding:"required"`
  ChainName     string                     `json:"chainName"     binding:"required"`
  Network       entities.DeploymentNetwork `json:"network"       binding:"required"`
  SeedPhrase    string                     `json:"seedPhrase"    binding:"required"`
  InfraProvider string                     `json:"infraProvider" binding:"required"`
  AwsAccessKey  string                     `json:"awsAccessKey"`
  AwsSecretKey  string                     `json:"awsSecretKey"`
  AwsRegion     string                     `json:"awsRegion"`
  L1RpcUrl      string                     `json:"l1RpcUrl"      binding:"required"`
  L1BeaconUrl   string                     `json:"l1BeaconUrl"   binding:"required"`
  FeeToken         string                     `json:"feeToken"`
  ReuseDeployment  *bool                      `json:"reuseDeployment,omitempty"`
  Overrides        []PresetFieldOverride      `json:"overrides,omitempty"`
}

// Validate checks provider-specific required fields.
func (r *PresetDeployRequest) ValidateProvider() error {
  switch r.InfraProvider {
  case "aws":
    if r.AwsAccessKey == "" || r.AwsSecretKey == "" || r.AwsRegion == "" {
      return fmt.Errorf("awsAccessKey, awsSecretKey, and awsRegion are required for aws provider")
    }
  case "local":
    if r.Network == "Mainnet" {
      return fmt.Errorf("local deployment is not supported for Mainnet")
    }
  default:
    return fmt.Errorf("invalid infraProvider %q: must be \"aws\" or \"local\"", r.InfraProvider)
  }
  return nil
}
```

### PresetFieldOverride
**라인**: 해당 타입은 Frontend 전용이지만, Backend는 generic slice로 처리

```go
Overrides []PresetFieldOverride `json:"overrides,omitempty"`
```

---

## 데이터 흐름 및 변환

### Frontend → Backend 매핑

| Frontend (form) | Backend (DTO) | 설명 |
|----------------|---------------|------|
| `form.getValues("presetBasicInfo.chainName")` | `PresetDeployRequest.ChainName` | 체인 이름 |
| `form.getValues("presetBasicInfo.network")` | `PresetDeployRequest.Network` | Testnet/Mainnet |
| `form.getValues("presetBasicInfo.seedPhrase").join(" ")` | `PresetDeployRequest.SeedPhrase` | 시드 프레이즈 (공백 분리) |
| `form.getValues("presetBasicInfo.infraProvider")` | `PresetDeployRequest.InfraProvider` | aws/local |
| `form.getValues("presetBasicInfo.awsAccessKey")` | `PresetDeployRequest.AwsAccessKey` | AWS 접근 키 |
| `form.getValues("presetBasicInfo.awsSecretKey")` | `PresetDeployRequest.AwsSecretKey` | AWS 비밀 키 |
| `form.getValues("presetBasicInfo.awsRegion")` | `PresetDeployRequest.AwsRegion` | AWS 리전 |
| `form.getValues("presetBasicInfo.l1RpcUrl")` | `PresetDeployRequest.L1RpcUrl` | Layer 1 RPC URL |
| `form.getValues("presetBasicInfo.l1BeaconUrl")` | `PresetDeployRequest.L1BeaconUrl` | Layer 1 Beacon URL |
| `form.getValues("presetBasicInfo.feeToken")` | `PresetDeployRequest.FeeToken` | TON/ETH/USDT/USDC |
| `selectedPresetId` | `PresetDeployRequest.PresetID` | Preset ID (general/defi/gaming/full) |
| `overrides` | `PresetDeployRequest.Overrides` | 파라미터 재정의 |
| `form.watch("presetBasicInfo.network") === "Mainnet" ? reuseDeployment : undefined` | `PresetDeployRequest.ReuseDeployment` | Mainnet만 적용 |

---

## 시퀀스 다이어그램 (텍스트 표현)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ User clicks "Deploy Rollup" button (Step 3)                                 │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ page.tsx: goToNextStep() triggered                                          │
│ → currentStep === 3: calls handleDeploy()                                   │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ usePresetWizard.ts: handleDeploy()                                          │
│ • Extract form values (chainName, network, infraProvider, etc.)             │
│ • Join seedPhrase array to string                                           │
│ • Construct DeployWithPresetRequest object                                  │
│ • Show toast: "Initiating deployment..."                                    │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ presetService.ts: startPresetDeployment(request)                            │
│ → Call apiPost<{ deploymentId: string }>(                                   │
│     "stacks/thanos/preset-deploy",                                          │
│     request                                                                  │
│   )                                                                          │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ api.ts: apiPost() - Axios HTTP request                                      │
│ • Method: POST                                                              │
│ • URL: /api/proxy/stacks/thanos/preset-deploy                               │
│ • Headers:                                                                   │
│   - Content-Type: application/json                                          │
│   - Authorization: Bearer {accessToken}                                     │
│   - ngrok-skip-browser-warning: 69420                                       │
│ • Body: JSON serialized DeployWithPresetRequest                             │
│ • Timeout: 10000ms                                                          │
│                                                                              │
│ Request Interceptor:                                                         │
│ • Adds Authorization header with token from localStorage                    │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
        ┌────────────────────┐
        │ Next.js /api/proxy │
        │ Proxies request to │
        │ trh-backend        │
        │ (actual URL        │
        │  localhost:8000)   │
        └────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ trh-backend: presets.go: PresetDeploy()                                      │
│ • Parse JSON request body: c.ShouldBindJSON(&request)                       │
│ • Validate:                                                                  │
│   - PresetID exists in preset service                                       │
│ • Call service: ThanosDeploymentService.CreateThanosStackFromPreset()       │
│ • Return: entities.Response with Status code and Data                       │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ HTTP Response (200 OK)                                                      │
│ {                                                                            │
│   "success": true,                                                          │
│   "message": "ok",                                                          │
│   "data": {                                                                 │
│     "deploymentId": "550e8400-e29b-41d4-a716-446655440000"                 │
│   }                                                                          │
│ }                                                                            │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ api.ts: Response Interceptor + apiPost() returns response.data               │
│ → Unwrap ApiResponse<{ deploymentId }>                                      │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ usePresetWizard.ts: handleDeploy() continues                                │
│ • const result = await startPresetDeployment(...)                           │
│ • toast.success("Deployment initiated!", { id: "preset-deploy" })          │
│ • setPendingDeploymentId(result.deploymentId)                               │
│ • resetState()                                                              │
│ • router.push("/rollup")                                                    │
└────────────────┬────────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ User navigated to /rollup page                                              │
│ Deployment process starts on backend asynchronously                         │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 에러 처리 시나리오

### 시나리오 1: JSON 바인딩 실패
**상황**: 잘못된 JSON 형식
**HTTP 상태코드**: 400 Bad Request
**응답**:
```json
{
  "status": 400,
  "message": "json: cannot unmarshal string into Go value of type int"
}
```
**Frontend**: toast.error(error.message)

### 시나리오 2: 필수 필드 누락
**상황**: PresetID 누락
**HTTP 상태코드**: 400 Bad Request
**응답**:
```json
{
  "status": 400,
  "message": "Key: 'PresetDeployRequest.PresetID' Error:Field validation for 'PresetID' failed on the 'required' tag"
}
```
**Frontend**: toast.error()

### 시나리오 3: Preset ID 유효하지 않음
**상황**: 존재하지 않는 Preset ID
**HTTP 상태코드**: 400 Bad Request
**응답**:
```json
{
  "status": 400,
  "message": "Preset not found: invalid-preset-id"
}
```
**Frontend**: toast.error()

### 시나리오 4: 네트워크 오류
**상황**: 서버 응답 없음
**HTTP 상태코드**: 0 (네트워크 에러)
**Frontend Response**:
```typescript
{
  message: "No response from server",
  status: 0
}
```
**Frontend**: toast.error("No response from server")

### 시나리오 5: 인증 실패
**상황**: 유효하지 않은 토큰
**HTTP 상태코드**: 401 Unauthorized
**Frontend 액션**:
- localStorage에서 accessToken 제거
- 사용자는 로그인 페이지로 리다이렉트 (별도 미들웨어)

### 시나리오 6: 권한 부족
**상황**: 일반 사용자 (Admin 권한 필요)
**HTTP 상태코드**: 403 Forbidden
**Frontend**: Console에 "Access denied" 로그 기록

### 시나리오 7: Mainnet + Local 제공자 충돌
**상황**: Mainnet과 Local 제공자 조합 선택
**HTTP 상태코드**: 400 Bad Request
**Backend 검증** (thanos.go, 라인 570-572):
```go
case "local":
  if r.Network == "Mainnet" {
    return fmt.Errorf("local deployment is not supported for Mainnet")
  }
```
**응답**:
```json
{
  "status": 400,
  "message": "local deployment is not supported for Mainnet"
}
```
**Frontend**: toast.error() - 사용자에게 에러 메시지 표시

### 시나리오 8: 유효하지 않은 L1 RPC/Beacon URL
**상황**: L1RpcUrl 또는 L1BeaconUrl이 잘못된 URL 형식
**Frontend 검증** (preset.ts, 라인 149-150):
```typescript
l1RpcUrl: z.string().url("Must be a valid URL"),
l1BeaconUrl: z.string().url("Must be a valid URL"),
```
**응답** (Frontend Zod 검증 실패):
```json
{
  "errors": [
    {
      "code": "invalid_string",
      "expected": "url",
      "message": "Must be a valid URL",
      "path": ["l1RpcUrl"]
    }
  ]
}
```
**Frontend**: 폼 레벨에서 검증 에러 표시 (Backend 호출 없음)

### 시나리오 9: Backend 배포 서비스 오류
**상황**: CreateThanosStackFromPreset() 서비스 호출 실패 (데이터베이스, 외부 API 등)
**Code 위치**: presets.go, 라인 98-103
```go
response, err := h.ThanosDeploymentService.CreateThanosStackFromPreset(c.Request.Context(), request)
if err != nil {
  logger.Error("failed to deploy preset stack", zap.Error(err))
}
c.JSON(int(response.Status), response)  // 위험: response가 nil이면 panic 발생 가능
```
**문제점**: 에러 발생 시 response가 nil이 될 수 있는데도 continue 실행
**HTTP 상태코드**: 500 Internal Server Error (또는 패닉)
**응답** (의도하지 않은 동작):
```json
{
  "status": 500,
  "message": "internal server error"
}
```
**Frontend**: toast.error() - Backend 에러 메시지 수신

---

## 알려진 문제점 (Known Issues)

### 1. Backend PresetDeploy 핸들러 에러 경로 버그
**코드 위치**: `/Users/theo/workspace_tokamak/trh-backend/pkg/api/handlers/thanos/presets.go`, 라인 98-103

```go
response, err := h.ThanosDeploymentService.CreateThanosStackFromPreset(c.Request.Context(), request)
if err != nil {
  logger.Error("failed to deploy preset stack", zap.Error(err))
  // Missing: return statement with error response
}
c.JSON(int(response.Status), response)
```

**문제**: CreateThanosStackFromPreset() 실패 시 에러만 로깅하고 계속 진행. response가 nil이면 `response.Status` 접근 시 nil pointer dereference 발생.

**영향**: 
- 배포 실패 상황이 제대로 클라이언트에 보고되지 않음
- 잠재적 서버 패닉 발생

**해결책**: 필수 - 에러 발생 시 HTTP 상태 코드와 함께 응답 반환
```go
if err != nil {
  logger.Error("failed to deploy preset stack", zap.Error(err))
  c.JSON(http.StatusInternalServerError, &entities.Response{
    Status:  http.StatusInternalServerError,
    Message: err.Error(),
  })
  return
}
```

### 2. 원본 명세의 미이행 사항

#### a. Network 옵션 - Sepolia 미지원
- **명세**: "network (enum: mainnet/testnet/sepolia)"
- **실제 구현**: Testnet/Mainnet만 지원
- **파일**: preset.ts, 라인 143

#### b. InfraProvider 옵션 - GCP/Azure 미지원
- **명세**: "infraProvider (enum: aws/gcp/azure or local)"
- **실제 구현**: aws/local만 지원
- **파일**: preset.ts, 라인 145
- **Backend 검증**: thanos.go, 라인 574

---

## 보안 및 민감 정보 처리

### 시드 프레이즈 (Seed Phrase)
- **Frontend**: 사용자 입력 받음 (12개 단어)
- **전송**: 공백 분리 문자열로 POST 바디에 포함
- **Backend**: 요청 수신 후 즉시 비우기 (로깅 방지)

**Backend Code** (deployment.go, 라인 86):
```go
// Sanitize seed phrase — clear it from the request after use so it is never
// written to logs or persisted in raw form.
request.SeedPhrase = ""
```

### AWS 자격증
- **Frontend**: localStorage에 저장되지 않음 (폼 입력만)
- **전송**: POST 바디에 포함
- **Backend**: 검증 후 환경에 따라 처리

### 인증 토큰
- **Frontend**: localStorage에서 "accessToken" 키로 저장
- **전송**: Authorization 헤더 (`Bearer {token}`)
- **Backend**: JWT 미들웨어에서 검증 (별도 구현)

---

## 주요 파일 목록

| 레이어 | 파일 경로 | 주요 함수/컴포넌트 |
|--------|----------|------------------|
| UI | `/trh-platform-ui/src/app/rollup/create/page.tsx` | `PresetWizardContent()` |
| 훅 | `/trh-platform-ui/src/features/rollup/hooks/usePresetWizard.ts` | `usePresetWizard()`, `handleDeploy()` |
| 서비스 | `/trh-platform-ui/src/features/rollup/services/presetService.ts` | `startPresetDeployment()` |
| API 클라이언트 | `/trh-platform-ui/src/lib/api.ts` | `apiPost()` |
| 스키마 | `/trh-platform-ui/src/features/rollup/schemas/preset.ts` | `DeployWithPresetRequest` |
| 백엔드 핸들러 | `/trh-backend/pkg/api/handlers/thanos/presets.go` | `PresetDeploy()` |
| 백엔드 DTO | `/trh-backend/pkg/api/dtos/thanos.go` | `PresetDeployRequest`, `PresetFieldOverride` |

---

## 추가 참고: 직접 배포 (DeployThanosRequest)

Frontend에서는 주로 Preset 기반을 사용하지만, 직접 배포도 가능합니다.

### 엔드포인트
- **URL**: `POST /api/v1/stacks/thanos`
- **Frontend**: `apiPost("stacks/thanos", request)`

### 요청 DTO
- **파일**: `/trh-backend/pkg/api/dtos/thanos.go`
- **라인**: 56-85

### Frontend 호출
- **파일**: `/trh-platform-ui/src/features/rollup/services/rollupService.ts`
- **함수**: `deployRollup()`
- **라인**: 29-34

```typescript
export const deployRollup = async (request: RollupDeploymentRequest) => {
  const response = await apiPost<{ id: string }>("stacks/thanos", request, {
    timeout: 30000, // Increase timeout for deployment
  });
  return response;
};
```

### 타입 정의
- **파일**: `/trh-platform-ui/src/features/rollup/schemas/create-rollup.ts`
- **타입**: `RollupDeploymentRequest`
- **라인**: 254-290

이 방식은 사용자가 모든 파라미터를 직접 입력해야 하므로 복잡도가 높습니다.

---

## 요약

**Phase 2의 핵심**: Preset 마법사 UI에서 사용자 입력을 수집 → 폼 검증 → JSON 페이로드 구성 → Axios를 통한 POST 요청 → Next.js 프록시 → trh-backend 핸들러 → 배포 서비스 호출

**주요 기술 스택**:
- Frontend: React Hook Form (폼 관리), Zod (유효성 검사), Axios (HTTP), React Query (상태 관리)
- Backend: Gin Framework (HTTP 라우팅), Go (구현)
- 통신: JSON over HTTP/HTTPS, Bearer Token 인증

**보안 고려사항**:
- Seed Phrase는 Backend에서 즉시 비우기
- AWS 자격증은 전송 중 TLS/SSL로 암호화
- 인증은 JWT 토큰 기반
- 401 응답 시 토큰 제거 및 재로그인 유도
