package main

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/log"
	op_challenger "github.com/tokamak-network/tokamak-thanos/op-challenger"
	opservice "github.com/tokamak-network/tokamak-thanos/op-service"
	"github.com/urfave/cli/v2"

	"github.com/tokamak-network/tokamak-thanos/op-challenger/config"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/flags"
	"github.com/tokamak-network/tokamak-thanos/op-challenger/version"
	"github.com/tokamak-network/tokamak-thanos/op-service/cliapp"
	oplog "github.com/tokamak-network/tokamak-thanos/op-service/log"
)

var (
	GitCommit = ""
	GitDate   = ""
)

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = opservice.FormatVersion(version.Version, GitCommit, GitDate, version.Meta)

func main() {
	args := os.Args
	if err := run(args, op_challenger.Main); err != nil {
		log.Crit("Application failed", "err", err)
	}
}

type ConfigAction func(ctx context.Context, log log.Logger, config *config.Config) error

func run(args []string, action ConfigAction) error {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Version = VersionWithMeta
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Name = "op-challenger"
	app.Usage = "Challenge outputs"
	app.Description = "Ensures that on chain outputs are correct."
	app.Action = func(ctx *cli.Context) error {
		logger, err := setupLogging(ctx)
		if err != nil {
			return err
		}
		logger.Info("Starting op-challenger", "version", VersionWithMeta)

		cfg, err := flags.NewConfigFromCLI(ctx)
		if err != nil {
			return err
		}
		return action(ctx.Context, logger, cfg)
	}
	return app.Run(args)
}

func setupLogging(ctx *cli.Context) (log.Logger, error) {
	logCfg := oplog.ReadCLIConfig(ctx)
	logger := oplog.NewLogger(oplog.AppOut(ctx), logCfg)
	oplog.SetGlobalLogHandler(logger.GetHandler())
	return logger, nil
}
