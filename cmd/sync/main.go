package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/grafana/fleet-management-sync-action/pkg/config"
	"github.com/grafana/fleet-management-sync-action/pkg/discovery"
	"github.com/grafana/fleet-management-sync-action/pkg/fleetmanagement"
)

func main() {
	// Get inputs from GitHub Action environment variables
	cfg, err := config.NewFromEnv()
	if err != nil {
		slog.Error("Invalid config", "error", err)
		os.Exit(1)
	}

	cfg.SetupLogging()
	slog.Info("Starting fm-sync",
		"root_path", cfg.RootPath,
		"timeout", cfg.Timeout.String())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	pipelines, err := discovery.FindPipelines(ctx, cfg)
	if err != nil {
		slog.Error("Failed to discover pipelines", "error", err)
		os.Exit(1)
	}

	slog.Info("Discovered pipelines", "count", len(pipelines))

	// Sync pipelines with Fleet Management
	if err := fleetmanagement.SyncPipelines(cfg, pipelines); err != nil {
		slog.Error("Failed to sync pipelines", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully completed sync")
}
