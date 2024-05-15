package etl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/indexer/config"
	"github.com/tokamak-network/tokamak-thanos/indexer/database"
	"github.com/tokamak-network/tokamak-thanos/indexer/node"
	"github.com/tokamak-network/tokamak-thanos/op-service/retry"
)

type L1ETL struct {
	ETL

	db        *database.DB
	mu        *sync.Mutex
	listeners []chan interface{}
}

// NewL1ETL creates a new L1ETL instance that will start indexing from different starting points
// depending on the state of the database and the supplied start height.
func NewL1ETL(cfg Config, log log.Logger, db *database.DB, metrics Metricer, client node.EthClient, contracts config.L1Contracts) (*L1ETL, error) {
	log = log.New("etl", "l1")

	zeroAddr := common.Address{}
	l1Contracts := []common.Address{}
	if err := contracts.ForEach(func(name string, addr common.Address) error {
		// Since we dont have backfill support yet, we want to make sure all expected
		// contracts are specified to ensure consistent behavior. Once backfill support
		// is ready, we can relax this requirement.
		if addr == zeroAddr && !strings.HasPrefix(name, "Legacy") {
			log.Error("address not configured", "name", name)
			return errors.New("all L1Contracts must be configured")
		}

		log.Info("configured contract", "name", name, "addr", addr)
		l1Contracts = append(l1Contracts, addr)
		return nil
	}); err != nil {
		return nil, err
	}

	latestHeader, err := db.Blocks.L1LatestBlockHeader()
	if err != nil {
		return nil, err
	}

	// Determine the starting height for traversal
	var fromHeader *types.Header
	if latestHeader != nil {
		log.Info("detected last indexed block", "number", latestHeader.Number, "hash", latestHeader.Hash)
		fromHeader = latestHeader.RLPHeader.Header()
	} else if cfg.StartHeight.BitLen() > 0 {
		log.Info("no indexed state starting from supplied L1 height", "height", cfg.StartHeight.String())
		header, err := client.BlockHeaderByNumber(cfg.StartHeight)
		if err != nil {
			return nil, fmt.Errorf("could not fetch starting block header: %w", err)
		}

		fromHeader = header
	} else {
		log.Info("no indexed state, starting from genesis")
	}

	// NOTE - The use of un-buffered channel here assumes that downstream consumers
	// will be able to keep up with the rate of incoming batches
	etlBatches := make(chan ETLBatch)
	etl := ETL{
		loopInterval:     time.Duration(cfg.LoopIntervalMsec) * time.Millisecond,
		headerBufferSize: uint64(cfg.HeaderBufferSize),

		log:             log,
		metrics:         metrics,
		headerTraversal: node.NewHeaderTraversal(client, fromHeader, cfg.ConfirmationDepth),
		contracts:       l1Contracts,
		etlBatches:      etlBatches,

		EthClient: client,
	}

	return &L1ETL{ETL: etl, db: db, mu: new(sync.Mutex)}, nil
}

func (l1Etl *L1ETL) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- l1Etl.ETL.Start(ctx)
	}()

	for {
		select {
		case err := <-errCh:
			return err

		// Index incoming batches (only L1 blocks that have an emitted log)
		case batch := <-l1Etl.etlBatches:
			l1BlockHeaders := make([]database.L1BlockHeader, 0, len(batch.Headers))
			for i := range batch.Headers {
				if _, ok := batch.HeadersWithLog[batch.Headers[i].Hash()]; ok {
					l1BlockHeaders = append(l1BlockHeaders, database.L1BlockHeader{BlockHeader: database.BlockHeaderFromHeader(&batch.Headers[i])})
				}
			}

			if len(l1BlockHeaders) == 0 {
				batch.Logger.Info("no l1 blocks with logs in batch")
				continue
			}

			l1ContractEvents := make([]database.L1ContractEvent, len(batch.Logs))
			for i := range batch.Logs {
				timestamp := batch.HeaderMap[batch.Logs[i].BlockHash].Time
				l1ContractEvents[i] = database.L1ContractEvent{ContractEvent: database.ContractEventFromLog(&batch.Logs[i], timestamp)}
				l1Etl.ETL.metrics.RecordIndexedLog(batch.Logs[i].Address)
			}

			// Continually try to persist this batch. If it fails after 10 attempts, we simply error out
			retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
			if _, err := retry.Do[interface{}](ctx, 10, retryStrategy, func() (interface{}, error) {
				if err := l1Etl.db.Transaction(func(tx *database.DB) error {
					if err := tx.Blocks.StoreL1BlockHeaders(l1BlockHeaders); err != nil {
						return err
					}
					// we must have logs if we have l1 blocks
					if err := tx.ContractEvents.StoreL1ContractEvents(l1ContractEvents); err != nil {
						return err
					}
					return nil
				}); err != nil {
					batch.Logger.Error("unable to persist batch", "err", err)
					return nil, err
				}

				l1Etl.ETL.metrics.RecordIndexedHeaders(len(l1BlockHeaders))
				l1Etl.ETL.metrics.RecordIndexedLatestHeight(l1BlockHeaders[len(l1BlockHeaders)-1].Number)

				// a-ok!
				return nil, nil
			}); err != nil {
				return err
			}

			batch.Logger.Info("indexed batch")

			// Notify Listeners
			l1Etl.mu.Lock()
			for i := range l1Etl.listeners {
				select {
				case l1Etl.listeners[i] <- struct{}{}:
				default:
					// do nothing if the listener hasn't picked
					// up the previous notif
				}
			}
			l1Etl.mu.Unlock()
		}
	}
}

// Notify returns a channel that'll receive a value every time new data has
// been persisted by the L1ETL
func (l1Etl *L1ETL) Notify() <-chan interface{} {
	receiver := make(chan interface{})
	l1Etl.mu.Lock()
	defer l1Etl.mu.Unlock()

	l1Etl.listeners = append(l1Etl.listeners, receiver)
	return receiver
}
