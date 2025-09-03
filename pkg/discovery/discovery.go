package discovery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/fleet-management-sync-action/pkg/config"
	"gopkg.in/yaml.v3"
)

// FindPipelines walks the filesystem starting from the RootPath in cfg and finds all YAML files, attempting to parse them as Pipeline configurations.
// If RootPath is empty, it defaults to the current working directory.
// The context can be used to cancel the discovery process.
//
// It returns a slice of all discovered pipelines, or an error if the discovery process fails.
func FindPipelines(ctx context.Context, cfg *config.Config) ([]*Pipeline, error) {
	rootPath := cfg.PipelinesRootPath

	var pipelines []*Pipeline
	pipelineNames := make(map[string]string) // for deduplication: name -> file path

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file has .yaml or .yml extension
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		// build .alloy path
		alloyPath := strings.TrimSuffix(path, ext) + ".alloy"

		// read and set pipeline contents
		contents, err := os.ReadFile(alloyPath)
		if err != nil {
			return fmt.Errorf("failed to read alloy file %s: %w", alloyPath, err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		var p Pipeline
		decoder := yaml.NewDecoder(strings.NewReader(string(data)))
		decoder.KnownFields(true)

		if err := decoder.Decode(&p); err != nil {
			return fmt.Errorf("failed to parse pipeline from %s: %w", path, err)
		}

		p.Contents = string(contents)

		// If no name provided, use filename without extension
		if p.Name == "" {
			base := filepath.Base(path)
			p.Name = strings.TrimSuffix(base, filepath.Ext(base))
		}

		// Check for duplicate pipeline names
		if existingPath, ok := pipelineNames[p.Name]; ok {
			return fmt.Errorf("duplicate pipeline name '%s' found in %s and %s", p.Name, existingPath, path)
		}
		pipelineNames[p.Name] = path

		pipelines = append(pipelines, &p)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return pipelines, nil
}


