package metrics

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	contractMetrics "github.com/tokamak-network/tokamak-thanos/op-challenger/game/fault/contracts/metrics"

	txmetrics "github.com/tokamak-network/tokamak-thanos/op-service/txmgr/metrics"
)

type NoopMetricsImpl struct {
	txmetrics.NoopTxMetrics
	contractMetrics.NoopMetrics
}

func (i *NoopMetricsImpl) StartBalanceMetrics(l log.Logger, client *ethclient.Client, account common.Address) io.Closer {
	return nil
}

var NoopMetrics Metricer = new(NoopMetricsImpl)

func (*NoopMetricsImpl) RecordInfo(version string) {}
func (*NoopMetricsImpl) RecordUp()                 {}

func (*NoopMetricsImpl) RecordGameMove()        {}
func (*NoopMetricsImpl) RecordGameStep()        {}
func (*NoopMetricsImpl) RecordGameL2Challenge() {}

func (*NoopMetricsImpl) RecordActedL1Block(_ uint64) {}

func (*NoopMetricsImpl) RecordPreimageChallenged()      {}
func (*NoopMetricsImpl) RecordPreimageChallengeFailed() {}

func (*NoopMetricsImpl) RecordBondClaimFailed()   {}
func (*NoopMetricsImpl) RecordBondClaimed(uint64) {}

func (*NoopMetricsImpl) RecordCannonExecutionTime(t float64)   {}
func (*NoopMetricsImpl) RecordAsteriscExecutionTime(t float64) {}
func (*NoopMetricsImpl) RecordClaimResolutionTime(t float64)   {}
func (*NoopMetricsImpl) RecordGameActTime(t float64)           {}

func (*NoopMetricsImpl) RecordGamesStatus(inProgress, defenderWon, challengerWon int) {}

func (*NoopMetricsImpl) RecordGameUpdateScheduled() {}
func (*NoopMetricsImpl) RecordGameUpdateCompleted() {}

func (*NoopMetricsImpl) IncActiveExecutors() {}
func (*NoopMetricsImpl) DecActiveExecutors() {}
func (*NoopMetricsImpl) IncIdleExecutors()   {}
func (*NoopMetricsImpl) DecIdleExecutors()   {}

func (*NoopMetricsImpl) CacheAdd(_ string, _ int, _ bool) {}
func (*NoopMetricsImpl) CacheGet(_ string, _ bool)        {}
