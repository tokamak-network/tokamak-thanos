package op_e2e

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

func TestCreateAndInitializePoolIfNecessary(t *testing.T) {
	// 로컬 devnet 설정
	devnetL1URL := "http://localhost:9545"
	chainID := big.NewInt(901)

	// 로컬 devnet에 연결
	rpcClient, err := rpc.Dial(devnetL1URL)
	require.NoError(t, err, "Failed to connect to devnetL1")

	l2Client := ethclient.NewClient(rpcClient)

	// 올바른 개인 키 설정 (이 예제에서는 임의의 개인 키를 사용)
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	require.NoError(t, err, "Failed to load private key")

	// TransactOpts 설정 (개인 키 사용)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err, "Failed to create keyed transactor")

	// 가스 한도 설정
	auth.GasLimit = uint64(5000000) // 필요에 따라 적절히 조정

	// nonfungiblePositionManager 인스턴스 생성 (L2)
	nonfungiblePositionManager, err := bindings.NewNonfungiblePositionManager(common.HexToAddress("0x4200000000000000000000000000000000000506"), l2Client)
	require.NoError(t, err, "Failed to create nonfungible position manager")

	// createAndInitializePoolIfNecessary 호출 (L2)
	tx, err := nonfungiblePositionManager.CreateAndInitializePoolIfNecessary(
		auth,
		common.HexToAddress("0x4200000000000000000000000000000000000006"), // Token0 주소
		common.HexToAddress("0x4200000000000000000000000000000000000486"), // Token1 주소
		big.NewInt(3000),        // Fee 값 (예: 3000)
		big.NewInt(42951287400), // sqrtPriceX96 값
	)
	require.NoError(t, err, "Failed to create and initialize pool")

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
