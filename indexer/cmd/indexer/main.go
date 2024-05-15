package main

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/log"
	oplog "github.com/tokamak-network/tokamak-thanos/op-service/log"
	"github.com/tokamak-network/tokamak-thanos/op-service/opio"
)

var (
	GitCommit = ""
	GitDate   = ""
)

func main() {
	// This is the most root context, used to propagate
	// cancellations to all spawned application-level goroutines
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		opio.BlockOnInterrupts()
		cancel()
	}()

	oplog.SetupDefaults()
	app := newCli(GitCommit, GitDate)
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error("application failed", "err", err)
		os.Exit(1)
	}
}
