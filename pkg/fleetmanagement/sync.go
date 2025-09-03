package fleetmanagement

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/grafana/fleet-management-sync-action/pkg/config"
	"github.com/grafana/fleet-management-sync-action/pkg/discovery"
)

// SyncPipelines syncs the discovered pipelines with the Fleet Management API.
//
// If cfg.DryRun is set, this will log the pipelines that would have been synced but
// not actually sync them.
func SyncPipelines(cfg *config.Config, pipelines []*discovery.Pipeline) error {
	slog.Info("starting pipeline sync",
		"username", cfg.Username,
		"pipeline_count", len(pipelines),
		"dry_run", cfg.DryRun)

	if !cfg.DryRun {
		// TODO: Implement actual sync with Fleet Management API
		slog.Error("actual sync not yet implemented")
		return fmt.Errorf("sync implementation pending")
	}

	// Run dry-run operations and log the pipelines that would have been synced
	for _, p := range pipelines {
		apiPipeline := p.ToFleetManagementPipeline()

		// Marshal to JSON for logging
		data, err := json.MarshalIndent(apiPipeline, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal pipeline %s: %w", p.Name, err)
		}

		slog.Info("would sync pipeline", "name", p.Name, "enabled", p.Enabled)
		slog.Debug("pipeline payload", "payload", string(data))
	}

	slog.Info("pipeline sync completed (dry run)")
	return nil
}
