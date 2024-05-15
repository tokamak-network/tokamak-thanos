package etl

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/indexer/node"
)

type Config struct {
	LoopIntervalMsec uint
	HeaderBufferSize uint

	StartHeight       *big.Int
	ConfirmationDepth *big.Int
}

type ETL struct {
	log     log.Logger
	metrics Metricer

	loopInterval     time.Duration
	headerBufferSize uint64
	headerTraversal  *node.HeaderTraversal

	contracts  []common.Address
	etlBatches chan ETLBatch

	EthClient node.EthClient
}

type ETLBatch struct {
	Logger log.Logger

	Headers   []types.Header
	HeaderMap map[common.Hash]*types.Header

	Logs           []types.Log
	HeadersWithLog map[common.Hash]bool
}

func (etl *ETL) Start(ctx context.Context) error {
	done := ctx.Done()
	pollTicker := time.NewTicker(etl.loopInterval)
	defer pollTicker.Stop()

	// A reference that'll stay populated between intervals
	// in the event of failures in order to retry.
	var headers []types.Header

	etl.log.Info("starting etl...")
	for {
		select {
		case <-done:
			etl.log.Info("stopping etl")
			return nil

		case <-pollTicker.C:
			done := etl.metrics.RecordInterval()
			if len(headers) > 0 {
				etl.log.Info("retrying previous batch")
			} else {
				newHeaders, err := etl.headerTraversal.NextHeaders(etl.headerBufferSize)
				if err != nil {
					etl.log.Error("error querying for headers", "err", err)
				} else if len(newHeaders) == 0 {
					etl.log.Warn("no new headers. etl at head?")
				} else {
					headers = newHeaders
				}

				latestHeader := etl.headerTraversal.LatestHeader()
				if latestHeader != nil {
					etl.metrics.RecordLatestHeight(latestHeader.Number)
				}
			}

			// only clear the reference if we were able to process this batch
			err := etl.processBatch(headers)
			if err == nil {
				headers = nil
			}

			done(err)
		}
	}
}

func (etl *ETL) processBatch(headers []types.Header) error {
	if len(headers) == 0 {
		return nil
	}

	firstHeader, lastHeader := headers[0], headers[len(headers)-1]
	batchLog := etl.log.New("batch_start_block_number", firstHeader.Number, "batch_end_block_number", lastHeader.Number)
	batchLog.Info("extracting batch", "size", len(headers))

	headerMap := make(map[common.Hash]*types.Header, len(headers))
	for i := range headers {
		header := headers[i]
		headerMap[header.Hash()] = &header
	}

	headersWithLog := make(map[common.Hash]bool, len(headers))
	filterQuery := ethereum.FilterQuery{FromBlock: firstHeader.Number, ToBlock: lastHeader.Number, Addresses: etl.contracts}
	logs, err := etl.EthClient.FilterLogs(filterQuery)
	if err != nil {
		batchLog.Info("failed to extract logs", "err", err)
		return err
	}

	if logs.ToBlockHeader.Number.Cmp(lastHeader.Number) != 0 {
		// Warn and simply wait for the provider to synchronize state
		batchLog.Warn("mismatch in FilterLog#ToBlock number", "queried_to_block_number", lastHeader.Number, "reported_to_block_number", logs.ToBlockHeader.Number)
		return fmt.Errorf("mismatch in FilterLog#ToBlock number")
	} else if logs.ToBlockHeader.Hash() != lastHeader.Hash() {
		batchLog.Error("mismatch in FitlerLog#ToBlock block hash!!!", "queried_to_block_hash", lastHeader.Hash().String(), "reported_to_block_hash", logs.ToBlockHeader.Hash().String())
		return fmt.Errorf("mismatch in FitlerLog#ToBlock block hash!!!")
	}

	if len(logs.Logs) > 0 {
		batchLog.Info("detected logs", "size", len(logs.Logs))
	}

	for i := range logs.Logs {
		log := logs.Logs[i]
		headersWithLog[log.BlockHash] = true
		if _, ok := headerMap[log.BlockHash]; !ok {
			// NOTE. Definitely an error state if the none of the headers were re-orged out in between
			// the blocks and logs retrieval operations. Unlikely as long as the confirmation depth has
			// been appropriately set or when we get to natively handling reorgs.
			batchLog.Error("log found with block hash not in the batch", "block_hash", logs.Logs[i].BlockHash, "log_index", logs.Logs[i].Index)
			return errors.New("parsed log with a block hash not in the batch")
		}
	}

	// ensure we use unique downstream references for the etl batch
	headersRef := headers
	etl.etlBatches <- ETLBatch{Logger: batchLog, Headers: headersRef, HeaderMap: headerMap, Logs: logs.Logs, HeadersWithLog: headersWithLog}
	return nil
}
