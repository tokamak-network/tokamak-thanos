package op_e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCreateAndInitializePoolIfNecessary_1(t *testing.T) {
	// 시스템 설정
	cfg := DefaultSystemConfig(t)
	t.Log("시스템 설정 완료")

	// 시스템 시작
	sys, err := cfg.Start(t)
	require.Nil(t, err, "시스템 시작 오류")
	defer sys.Close()
	t.Log("시스템 시작 성공")

	// L2 클라이언트 설정
	l2Client := sys.Clients["sequencer"]
	t.Log("L2 클라이언트 설정 완료")

	// TransactOpts 설정 (Alice 계정 사용)
	auth, err := bind.NewKeyedTransactorWithChainID(sys.cfg.Secrets.Alice, cfg.L2ChainIDBig())
	require.NoError(t, err, "키드 트랜잭터 생성 오류")
	t.Log("TransactOpts 설정 완료")

	// 가스 한도 설정
	auth.GasLimit = uint64(5000000) // 가스 한도를 더 높게 설정

	// nonfungiblePositionManager 인스턴스 생성 (L2)
	nonfungiblePositionManager, err := bindings.NewNonfungiblePositionManager(common.HexToAddress("0x4200000000000000000000000000000000000506"), l2Client)
	require.NoError(t, err, "비대체성 포지션 매니저 생성 오류")
	t.Log("NonfungiblePositionManager 인스턴스 생성 완료")

	// createAndInitializePoolIfNecessary 호출 (L2)
	tx, err := nonfungiblePositionManager.CreateAndInitializePoolIfNecessary(
		auth,
		common.HexToAddress("0x4200000000000000000000000000000000000006"), // Token0 주소
		common.HexToAddress("0x4200000000000000000000000000000000000486"), // Token1 주소
		big.NewInt(3000),        // Fee 값 (예: 3000)
		big.NewInt(42951287400), // sqrtPriceX96 값
	)
	require.NoError(t, err, "풀 생성 및 초기화 오류")
	t.Logf("트랜잭션 해시 (L2): %s", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), l2Client, tx)
	require.NoError(t, err, "트랜잭션 마이닝 오류")

	// 트랜잭션 상태 코드 확인
	t.Logf("트랜잭션 상태 코드: %d", receipt.Status)
	if receipt.Status != 1 {
		t.Fatalf("트랜잭션 실패, 상태 코드: %d", receipt.Status)
	}

	require.Equal(t, receipt.Status, uint64(1), "트랜잭션 실패")

	// UniswapV3Factory 인스턴스 생성
	uniswapV3Factory, err := bindings.NewUniswapV3Factory(common.HexToAddress("0x4200000000000000000000000000000000000504"), l2Client) // 정확한 주소 사용
	require.NoError(t, err, "Uniswap V3 팩토리 생성 오류")
	t.Log("UniswapV3Factory 인스턴스 생성 완료")

	// getPool 호출
	poolAddress, err := uniswapV3Factory.GetPool(&bind.CallOpts{},
		common.HexToAddress("0x4200000000000000000000000000000000000006"),
		common.HexToAddress("0x4200000000000000000000000000000000000486"),
		big.NewInt(3000))
	require.NoError(t, err, "풀 주소 조회 오류")
	t.Logf("풀 주소: %s", poolAddress.Hex())
}
