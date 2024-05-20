package core

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	ZKHosts   []string `env:"ZKHOSTS, required"`
	ZKTimeout int      `env:"ZKTIMEOUT, required"`
	RestPort  string   `env:"RESTPORT, required"`
}

func prepareConfig() AppConfig {
	ctx := context.Background()
	var c AppConfig
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	return c
}

var (
	Config AppConfig = prepareConfig()
)
