package main

import (
	"flag"
	"github.com/sankalp-r/metrics-agent/pkg/agent"
	"go.uber.org/zap"
)

func main() {
	logConfig := zap.NewProductionConfig()
	logger, _ := logConfig.Build()
	zap.ReplaceGlobals(logger)

	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.Parse()
	config, err := agent.NewConfig(configPath)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	minimalAgent, err := agent.NewBuilder(config).Build()
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	zap.L().Info("starting the agent...")
	minimalAgent.Start()
}
