package clients

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tokamak-network/tokamak-thanos/op-ufm/pkg/metrics"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	signer "github.com/tokamak-network/tokamak-thanos/op-service/signer"
	optls "github.com/tokamak-network/tokamak-thanos/op-service/tls"
)

type InstrumentedSignerClient struct {
	c            *signer.SignerClient
	providerName string
}

func NewSignerClient(providerName string, logger log.Logger, endpoint string, tlsConfig optls.CLIConfig) (*InstrumentedSignerClient, error) {
	start := time.Now()
	c, err := signer.NewSignerClient(logger, endpoint, tlsConfig)
	if err != nil {
		metrics.RecordErrorDetails(providerName, "signer.NewSignerClient", err)
		return nil, err
	}
	metrics.RecordRPCLatency(providerName, "signer", "NewSignerClient", time.Since(start))
	return &InstrumentedSignerClient{c: c, providerName: providerName}, nil
}

func (i *InstrumentedSignerClient) SignTransaction(ctx context.Context, chainId *big.Int, from *common.Address, tx *types.Transaction) (*types.Transaction, error) {
	start := time.Now()
	tx, err := i.c.SignTransaction(ctx, chainId, *from, tx)
	if err != nil {
		metrics.RecordErrorDetails(i.providerName, "signer.SignTransaction", err)
		return nil, err
	}
	metrics.RecordRPCLatency(i.providerName, "signer", "SignTransaction", time.Since(start))
	return tx, err
}
