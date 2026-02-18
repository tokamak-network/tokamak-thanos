package metrics

import (
	"github.com/tokamak-network/tokamak-thanos/op-service/eth"
	"github.com/tokamak-network/tokamak-thanos/op-service/event"
	opmetrics "github.com/tokamak-network/tokamak-thanos/op-service/metrics"
	"github.com/tokamak-network/tokamak-thanos/op-supervisor/supervisor/types"
)

type noopMetrics struct {
	opmetrics.NoopRPCMetrics
	event.NoopMetrics
}

var NoopMetrics Metricer = new(noopMetrics)

func (*noopMetrics) Document() []opmetrics.DocumentedMetric { return nil }

func (*noopMetrics) RecordInfo(version string) {}
func (*noopMetrics) RecordUp()                 {}

func (m *noopMetrics) RecordCrossUnsafe(_ eth.ChainID, _ types.BlockSeal) {}
func (m *noopMetrics) RecordCrossSafe(_ eth.ChainID, _ types.BlockSeal)   {}
func (m *noopMetrics) RecordLocalSafe(_ eth.ChainID, _ types.BlockSeal)   {}
func (m *noopMetrics) RecordLocalUnsafe(_ eth.ChainID, _ types.BlockSeal) {}

func (m *noopMetrics) CacheAdd(_ eth.ChainID, _ string, _ int, _ bool) {}
func (m *noopMetrics) CacheGet(_ eth.ChainID, _ string, _ bool)        {}

func (m *noopMetrics) RecordDBEntryCount(_ eth.ChainID, _ string, _ int64) {}
func (m *noopMetrics) RecordDBSearchEntriesRead(_ eth.ChainID, _ int64)    {}

func (m *noopMetrics) RecordAccessListVerifyFailure(_ eth.ChainID) {}
