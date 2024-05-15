package op_challenger

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/config"
	"github.com/tokamak-network/tokamak-thanos/op-service/testlog"
)

func TestMainShouldReturnErrorWhenConfigInvalid(t *testing.T) {
	cfg := &config.Config{}
	err := Main(context.Background(), testlog.Logger(t, log.LvlInfo), cfg)
	require.ErrorIs(t, err, cfg.Check())
}
