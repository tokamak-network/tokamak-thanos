package main

import (
	"context"
	"os"

	"github.com/tokamak-network/tokamak-thanos/op-service/ctxinterrupt"

	opservice "github.com/tokamak-network/tokamak-thanos/op-service"
	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"
	"github.com/tokamak-network/tokamak-thanos/op-dripper/dripper"
	"github.com/tokamak-network/tokamak-thanos/op-dripper/flags"
	"github.com/tokamak-network/tokamak-thanos/op-dripper/metrics"
	"github.com/tokamak-network/tokamak-thanos/op-service/cliapp"
	oplog "github.com/tokamak-network/tokamak-thanos/op-service/log"
	"github.com/tokamak-network/tokamak-thanos/op-service/metrics/doc"
)

var (
	Version   = "v0.0.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "op-dripper"
	app.Usage = "Drippie Executor"
	app.Description = "Service for executing Drippie drips"
	app.Action = cliapp.LifecycleCmd(dripper.Main(Version))
	app.Commands = []*cli.Command{
		{
			Name:        "doc",
			Subcommands: doc.NewSubcommands(metrics.NewMetrics("default")),
		},
	}

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}
