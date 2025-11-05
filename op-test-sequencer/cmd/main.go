package main

import (
	"context"
	"io"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"

	opservice "github.com/tokamak-network/tokamak-thanos/op-service"
	"github.com/tokamak-network/tokamak-thanos/op-service/cliapp"
	"github.com/tokamak-network/tokamak-thanos/op-service/ctxinterrupt"
	oplog "github.com/tokamak-network/tokamak-thanos/op-service/log"
	"github.com/tokamak-network/tokamak-thanos/op-service/metrics/doc"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/config"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/flags"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/metrics"
	"github.com/tokamak-network/tokamak-thanos/op-test-sequencer/sequencer"
)

var (
	Version   = "v0.0.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := run(ctx, os.Stdout, os.Stderr, os.Args, fromConfig)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

func run(ctx context.Context, w io.Writer, ew io.Writer, args []string, fn sequencer.MainFn) error {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Writer = w
	app.ErrWriter = ew
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Version = opservice.FormatVersion(Version, GitCommit, GitDate, "")
	app.Name = "op-test-sequencer"
	app.Usage = "op-test-sequencer sequences blocks"
	app.Description = "op-test-sequencer sequences blocks"
	app.Action = cliapp.LifecycleCmd(sequencer.Main(app.Version, fn))
	app.Commands = []*cli.Command{
		{
			Name:        "doc",
			Subcommands: doc.NewSubcommands(metrics.NewMetrics("default")),
		},
	}
	return app.RunContext(ctx, args)
}

func fromConfig(ctx context.Context, cfg *config.Config, logger log.Logger) (cliapp.Lifecycle, error) {
	return sequencer.FromConfig(ctx, cfg, logger)
}
