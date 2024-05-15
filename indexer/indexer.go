package indexer

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/tokamak-network/tokamak-thanos/indexer/config"
	"github.com/tokamak-network/tokamak-thanos/indexer/database"
	"github.com/tokamak-network/tokamak-thanos/indexer/etl"
	"github.com/tokamak-network/tokamak-thanos/indexer/node"
	"github.com/tokamak-network/tokamak-thanos/indexer/processors"
	"github.com/tokamak-network/tokamak-thanos/indexer/processors/bridge"
	"github.com/tokamak-network/tokamak-thanos/op-service/httputil"
	"github.com/tokamak-network/tokamak-thanos/op-service/metrics"
)

// Indexer contains the necessary resources for
// indexing the configured L1 and L2 chains
type Indexer struct {
	log log.Logger
	db  *database.DB

	httpConfig      config.ServerConfig
	metricsConfig   config.ServerConfig
	metricsRegistry *prometheus.Registry

	L1ETL           *etl.L1ETL
	L2ETL           *etl.L2ETL
	BridgeProcessor *processors.BridgeProcessor
}

// NewIndexer initializes an instance of the Indexer
func NewIndexer(
	log log.Logger,
	db *database.DB,
	chainConfig config.ChainConfig,
	rpcsConfig config.RPCsConfig,
	httpConfig config.ServerConfig,
	metricsConfig config.ServerConfig,
) (*Indexer, error) {
	metricsRegistry := metrics.NewRegistry()

	// L1
	l1EthClient, err := node.DialEthClient(rpcsConfig.L1RPC, node.NewMetrics(metricsRegistry, "l1"))
	if err != nil {
		return nil, err
	}
	l1Cfg := etl.Config{
		LoopIntervalMsec:  chainConfig.L1PollingInterval,
		HeaderBufferSize:  chainConfig.L1HeaderBufferSize,
		ConfirmationDepth: big.NewInt(int64(chainConfig.L1ConfirmationDepth)),
		StartHeight:       big.NewInt(int64(chainConfig.L1StartingHeight)),
	}
	l1Etl, err := etl.NewL1ETL(l1Cfg, log, db, etl.NewMetrics(metricsRegistry, "l1"), l1EthClient, chainConfig.L1Contracts)
	if err != nil {
		return nil, err
	}

	// L2 (defaults to predeploy contracts)
	l2EthClient, err := node.DialEthClient(rpcsConfig.L2RPC, node.NewMetrics(metricsRegistry, "l2"))
	if err != nil {
		return nil, err
	}
	l2Cfg := etl.Config{
		LoopIntervalMsec:  chainConfig.L2PollingInterval,
		HeaderBufferSize:  chainConfig.L2HeaderBufferSize,
		ConfirmationDepth: big.NewInt(int64(chainConfig.L2ConfirmationDepth)),
	}
	l2Etl, err := etl.NewL2ETL(l2Cfg, log, db, etl.NewMetrics(metricsRegistry, "l2"), l2EthClient, chainConfig.L2Contracts)
	if err != nil {
		return nil, err
	}

	// Bridge
	bridgeProcessor, err := processors.NewBridgeProcessor(log, db, bridge.NewMetrics(metricsRegistry), l1Etl, chainConfig)
	if err != nil {
		return nil, err
	}

	indexer := &Indexer{
		log: log,
		db:  db,

		httpConfig:      httpConfig,
		metricsConfig:   metricsConfig,
		metricsRegistry: metricsRegistry,

		L1ETL:           l1Etl,
		L2ETL:           l2Etl,
		BridgeProcessor: bridgeProcessor,
	}

	return indexer, nil
}

func (i *Indexer) startHttpServer(ctx context.Context) error {
	i.log.Debug("starting http server...", "port", i.httpConfig.Host)

	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/healthz"))

	addr := net.JoinHostPort(i.httpConfig.Host, strconv.Itoa(i.httpConfig.Port))
	srv, err := httputil.StartHTTPServer(addr, r)
	if err != nil {
		return fmt.Errorf("http server failed to start: %w", err)
	}
	i.log.Info("http server started", "addr", srv.Addr())
	<-ctx.Done()
	defer i.log.Info("http server stopped")
	return srv.Stop(context.Background())
}

func (i *Indexer) startMetricsServer(ctx context.Context) error {
	i.log.Debug("starting metrics server...", "port", i.metricsConfig.Port)
	srv, err := metrics.StartServer(i.metricsRegistry, i.metricsConfig.Host, i.metricsConfig.Port)
	if err != nil {
		return fmt.Errorf("metrics server failed to start: %w", err)
	}
	i.log.Info("metrics server started", "addr", srv.Addr())
	<-ctx.Done()
	defer i.log.Info("metrics server stopped")
	return srv.Stop(context.Background())
}

// Start starts the indexing service on L1 and L2 chains
func (i *Indexer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 5)

	// if any goroutine halts, we stop the entire indexer
	processCtx, processCancel := context.WithCancel(ctx)
	runProcess := func(start func(ctx context.Context) error) {
		wg.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					i.log.Error("halting indexer on panic", "err", err)
					debug.PrintStack()
					errCh <- fmt.Errorf("panic: %v", err)
				}

				processCancel()
				wg.Done()
			}()

			errCh <- start(processCtx)
		}()
	}

	// Kick off all the dependent routines
	runProcess(i.L1ETL.Start)
	runProcess(i.L2ETL.Start)
	runProcess(i.BridgeProcessor.Start)
	runProcess(i.startMetricsServer)
	runProcess(i.startHttpServer)
	wg.Wait()

	err := <-errCh
	if err != nil {
		i.log.Error("indexer stopped", "err", err)
	} else {
		i.log.Info("indexer stopped")
	}

	return err
}
