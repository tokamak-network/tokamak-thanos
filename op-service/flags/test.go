package flags

import (
	"flag"

	"github.com/tokamak-network/tokamak-thanos/op-service/log"
)

type TestConfig struct {
	LogConfig log.CLIConfig
}

func ReadTestConfig() TestConfig {
	flag.Parse()

	cfg := log.ReadTestCLIConfig()

	return TestConfig{
		LogConfig: cfg,
	}
}
