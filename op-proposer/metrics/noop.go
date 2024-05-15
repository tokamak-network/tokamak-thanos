package metrics

import (
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
	opmetrics "github.com/tokamak-network/tokamak-thanos/op-service/metrics"
	txmetrics "github.com/tokamak-network/tokamak-thanos/op-service/txmgr/metrics"
)

type noopMetrics struct {
	opmetrics.NoopRefMetrics
	txmetrics.NoopTxMetrics
}

var NoopMetrics Metricer = new(noopMetrics)

func (*noopMetrics) RecordInfo(version string) {}
func (*noopMetrics) RecordUp()                 {}

func (*noopMetrics) RecordL2BlocksProposed(l2ref eth.L2BlockRef) {}
