package interop

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
	oplog "github.com/tokamak-network/tokamak-thanos/op-service/log"
)

func init() {
	oplog.SetGlobalLogHandler(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelDebug, true))
}
