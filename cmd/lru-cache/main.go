package main

import (
	"context"
	"flag"
	"log"
	"lru-cache/internal/srv"

	"github.com/caarlos0/env/v11"
)

func main() {
	var cfg srv.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse env vars: %v", err)
	}
	flag.StringVar(&cfg.HostPort, "server-host-port", cfg.HostPort, "Server host and port")
	flag.IntVar(&cfg.CacheSize, "cache-size", cfg.CacheSize, "Cache size")
	flag.DurationVar(&cfg.DefaultTTL, "default-cache-ttl", cfg.DefaultTTL, "Default cache TTL")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Log level")
	flag.Parse()

	srv, err := srv.New(cfg)
	if err != nil {
		log.Fatalf("Failed to configure a server: %v", err)
	}
	err = srv.Run(context.Background())
	if err != nil {
		log.Fatalf("Failed to run a server: %v", err)
	}

}
