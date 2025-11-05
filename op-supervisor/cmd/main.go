package main

import (
	"context"
	"os"

	"github.com/tokamak-network/tokamak-thanos/op-supervisor/config"
	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"

	opservice "github.com/tokamak-network/tokamak-thanos/op-service"
	"github.com/tokamak-network/tokamak-thanos/op-service/cliapp"
	"github.com/tokamak-network/tokamak-thanos/op-service/ctxinterrupt"
	oplog "github.com/tokamak-network/tokamak-thanos/op-service/log"
	"github.com/tokamak-network/tokamak-thanos/op-service/metrics/doc"
	"github.com/tokamak-network/tokamak-thanos/op-supervisor/flags"
	"github.com/tokamak-network/tokamak-thanos/op-supervisor/metrics"
	"github.com/tokamak-network/tokamak-thanos/op-supervisor/supervisor"
)

var (
	Version   = "v0.0.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := run(ctx, os.Args, fromConfig)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

func run(ctx context.Context, args []string, fn supervisor.MainFn) error {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "op-supervisor"
	app.Usage = "op-supervisor monitors cross-L2 interop messaging"
	app.Description = "The op-supervisor monitors cross-L2 interop messaging by pre-fetching events and then resolving the cross-L2 dependencies to answer safety queries."
	app.Action = cliapp.LifecycleCmd(supervisor.Main(app.Version, fn))
	app.Commands = []*cli.Command{
		{
			Name:        "doc",
			Subcommands: doc.NewSubcommands(metrics.NewMetrics("default")),
		},
	}
	return app.RunContext(ctx, args)
}

func fromConfig(ctx context.Context, cfg *config.Config, logger log.Logger) (cliapp.Lifecycle, error) {
	return supervisor.SupervisorFromConfig(ctx, cfg, logger)
}
