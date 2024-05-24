package flags

import (
	"github.com/urfave/cli/v2"

	"github.com/tokamak-network/tokamak-thanos/op-node/flags"
	opflags "github.com/tokamak-network/tokamak-thanos/op-service/flags"
	oplog "github.com/tokamak-network/tokamak-thanos/op-service/log"
	opmetrics "github.com/tokamak-network/tokamak-thanos/op-service/metrics"
	oprpc "github.com/tokamak-network/tokamak-thanos/op-service/rpc"
)

const envVarPrefix = "OP_BOOTNODE"

var Flags = []cli.Flag{
	opflags.CLINetworkFlag(envVarPrefix, ""),
	opflags.CLIRollupConfigFlag(envVarPrefix, ""),
}

func init() {
	Flags = append(Flags, flags.P2PFlags(envVarPrefix)...)
	Flags = append(Flags, opmetrics.CLIFlags(envVarPrefix)...)
	Flags = append(Flags, oplog.CLIFlags(envVarPrefix)...)
	Flags = append(Flags, oprpc.CLIFlags(envVarPrefix)...)
}
