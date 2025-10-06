package fleetmanagement

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	pipelinev1 "github.com/grafana/fleet-management-api/api/gen/proto/go/pipeline/v1"
	"github.com/grafana/fleet-management-api/api/gen/proto/go/pipeline/v1/pipelinev1connect"
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

	if cfg.DryRun {
		// Run dry-run operations and log the pipelines that would have been synced
		return performDryRun(pipelines)
	}

	err := performSync(cfg, pipelines)
	if err != nil {
		slog.Error("failed to sync pipelines", "error", err)
		return err
	}

	return nil
}

func performDryRun(pipelines []*discovery.Pipeline) error {
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

func performSync(cfg *config.Config, pipelines []*discovery.Pipeline) error {
	fmPipelines := make([]*pipelinev1.Pipeline, 0, len(pipelines))
	for _, p := range pipelines {
		// Add global matcher if configured and non-empty
		if cfg.GlobalMatcher != "" {
			p.Matchers = append(p.Matchers, cfg.GlobalMatcher)
		}
		fmPipelines = append(fmPipelines, p.ToFleetManagementPipeline())
	}

	req := &connect.Request[pipelinev1.SyncPipelinesRequest]{
		Msg: &pipelinev1.SyncPipelinesRequest{
			Source: &pipelinev1.PipelineSource{
				Type:      pipelinev1.PipelineSource_SOURCE_TYPE_GIT,
				Namespace: cfg.Namespace,
			},
			Pipelines: fmPipelines,
		},
	}

	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	buf := fmt.Appendf(nil, "%s:%s", cfg.Username, cfg.Token)
	credentials := base64.StdEncoding.EncodeToString(buf)

	client := pipelinev1connect.NewPipelineServiceClient(
		httpClient,
		cfg.FleetManagementURL,
		connect.WithInterceptors(basicAuthInterceptor(credentials)),
	)

	ctx := context.Background()
	_, err := client.SyncPipelines(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to sync pipelines: %w", err)
	}

	slog.Info("pipeline sync completed")

	return nil
}

func basicAuthInterceptor(credentials string) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Set("Authorization", fmt.Sprintf("Basic %s", credentials))
			return next(ctx, req)
		}
	})
}
