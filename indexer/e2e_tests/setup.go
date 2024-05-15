package e2e_tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tokamak-network/tokamak-thanos/indexer"
	"github.com/tokamak-network/tokamak-thanos/indexer/api"
	"github.com/tokamak-network/tokamak-thanos/indexer/client"
	"github.com/tokamak-network/tokamak-thanos/indexer/config"
	"github.com/tokamak-network/tokamak-thanos/indexer/database"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	op_e2e "github.com/tokamak-network/tokamak-thanos/op-e2e"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

/*
	NOTE - Most of the current bridge tests fetch chain data via direct database queries. These could all
	be transitioned to use the API client instead to better simulate/validate real-world usage.
	Supporting this would potentially require adding new API endpoints for the specific query lookup types.
*/

type E2ETestSuite struct {
	t *testing.T

	// API
	Client *client.Client
	API    *api.API

	// Indexer
	DB      *database.DB
	Indexer *indexer.Indexer

	// Rollup
	OpCfg *op_e2e.SystemConfig
	OpSys *op_e2e.System

	// Clients
	L1Client *ethclient.Client
	L2Client *ethclient.Client
}

// createE2ETestSuite ... Create a new E2E test suite
func createE2ETestSuite(t *testing.T) E2ETestSuite {
	dbUser := os.Getenv("DB_USER")
	dbName := setupTestDatabase(t)

	// Rollup System Configuration. Unless specified,
	// omit logs emitted by the various components. Maybe
	// we can eventually dump these logs to a temp file
	log.Root().SetHandler(log.DiscardHandler())
	opCfg := op_e2e.DefaultSystemConfig(t)
	if len(os.Getenv("ENABLE_ROLLUP_LOGS")) == 0 {
		t.Log("set env 'ENABLE_ROLLUP_LOGS' to show rollup logs")
		for name, logger := range opCfg.Loggers {
			t.Logf("discarding logs for %s", name)
			logger.SetHandler(log.DiscardHandler())
		}
	}

	// Rollup Start
	opSys, err := opCfg.Start(t)
	require.NoError(t, err)
	t.Cleanup(func() { opSys.Close() })

	// Indexer Configuration and Start
	indexerCfg := config.Config{
		DB: config.DBConfig{
			Host: "127.0.0.1",
			Port: 5432,
			Name: dbName,
			User: dbUser,
		},
		RPCs: config.RPCsConfig{
			L1RPC: opSys.EthInstances["l1"].HTTPEndpoint(),
			L2RPC: opSys.EthInstances["sequencer"].HTTPEndpoint(),
		},
		Chain: config.ChainConfig{
			L1PollingInterval: uint(opCfg.DeployConfig.L1BlockTime) * 1000,
			L2PollingInterval: uint(opCfg.DeployConfig.L2BlockTime) * 1000,
			L2Contracts:       config.L2ContractsFromPredeploys(),
			L1Contracts: config.L1Contracts{
				AddressManager:              opCfg.L1Deployments.AddressManager,
				SystemConfigProxy:           opCfg.L1Deployments.SystemConfigProxy,
				OptimismPortalProxy:         opCfg.L1Deployments.OptimismPortalProxy,
				L2OutputOracleProxy:         opCfg.L1Deployments.L2OutputOracleProxy,
				L1CrossDomainMessengerProxy: opCfg.L1Deployments.L1CrossDomainMessengerProxy,
				L1StandardBridgeProxy:       opCfg.L1Deployments.L1StandardBridgeProxy,
				L1ERC721BridgeProxy:         opCfg.L1Deployments.L1ERC721BridgeProxy,
			},
		},
		HTTPServer:    config.ServerConfig{Host: "127.0.0.1", Port: 0},
		MetricsServer: config.ServerConfig{Host: "127.0.0.1", Port: 0},
	}

	// E2E tests can run on the order of magnitude of minutes. Once
	// the system is running, mark this test for Parallel execution
	t.Parallel()

	// provide a DB for the unit test. disable logging
	silentLog := testlog.Logger(t, log.LvlInfo)
	silentLog.SetHandler(log.DiscardHandler())
	db, err := database.NewDB(silentLog, indexerCfg.DB)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	indexerLog := testlog.Logger(t, log.LvlInfo).New("role", "indexer")
	indexer, err := indexer.NewIndexer(indexerLog, db, indexerCfg.Chain, indexerCfg.RPCs, indexerCfg.HTTPServer, indexerCfg.MetricsServer)
	require.NoError(t, err)

	indexerCtx, indexerStop := context.WithCancel(context.Background())
	go func() {
		err := indexer.Run(indexerCtx)
		if err != nil { // panicking here ensures that the test will exit
			// during service failure. Using t.Fail() wouldn't be caught
			// until all awaiting routines finish which would never happen.
			panic(err)
		}
	}()

	apiLog := testlog.Logger(t, log.LvlInfo).New("role", "indexer_api")

	apiCfg := config.ServerConfig{
		Host: "127.0.0.1",
		Port: 0,
	}

	mCfg := config.ServerConfig{
		Host: "127.0.0.1",
		Port: 0,
	}

	api := api.NewApi(apiLog, db.BridgeTransfers, apiCfg, mCfg)
	apiCtx, apiStop := context.WithCancel(context.Background())
	go func() {
		err := api.Run(apiCtx)
		if err != nil {
			panic(err)
		}
	}()

	t.Cleanup(func() {
		apiStop()
		indexerStop()
	})

	// Wait for the API to start listening
	time.Sleep(1 * time.Second)

	client, err := client.NewClient(&client.Config{
		PaginationLimit: 100,
		BaseURL:         fmt.Sprintf("http://%s:%d", apiCfg.Host, api.Port()),
	})

	require.NoError(t, err)

	return E2ETestSuite{
		t:        t,
		Client:   client,
		DB:       db,
		Indexer:  indexer,
		OpCfg:    &opCfg,
		OpSys:    opSys,
		L1Client: opSys.Clients["l1"],
		L2Client: opSys.Clients["sequencer"],
	}
}

func setupTestDatabase(t *testing.T) string {
	user := os.Getenv("DB_USER")
	require.NotEmpty(t, user, "DB_USER env variable expected to instantiate test database")

	pg, err := sql.Open("pgx", fmt.Sprintf("postgres://%s@localhost:5432?sslmode=disable", user))
	require.NoError(t, err)
	require.NoError(t, pg.Ping())

	// create database
	dbName := fmt.Sprintf("indexer_test_%d", time.Now().UnixNano())
	_, err = pg.Exec("CREATE DATABASE " + dbName)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := pg.Exec("DROP DATABASE " + dbName)
		require.NoError(t, err)
		pg.Close()
	})

	dbConfig := config.DBConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		Name:     dbName,
		User:     user,
		Password: "",
	}

	silentLog := log.New()
	silentLog.SetHandler(log.DiscardHandler())
	db, err := database.NewDB(silentLog, dbConfig)
	require.NoError(t, err)
	defer db.Close()

	err = db.ExecuteSQLMigration("../migrations")
	require.NoError(t, err)

	t.Logf("database %s setup and migrations executed", dbName)
	return dbName
}
