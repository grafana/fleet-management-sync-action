package discovery

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/fleet-management-sync-action/pkg/config"
	"github.com/grafana/fleet-management-sync-action/pkg/pipeline"
	"gopkg.in/yaml.v3"
)

const (
	pipelineMetadataHeader = "/* fleet-management"
	pipelineMetadataFooter = "*/"
)

// FindPipelines walks the filesystem starting from the RootPath in cfg and finds all YAML files, attempting to parse them as Pipeline configurations.
// If RootPath is empty, it defaults to the current working directory.
// The context can be used to cancel the discovery process.
//
// It returns a slice of all discovered pipelines, or an error if the discovery process fails.
func FindPipelines(ctx context.Context, cfg *config.Config) ([]*pipeline.Pipeline, error) {
	rootPath := cfg.PipelinesRootPath
	if rootPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
		rootPath = cwd
	}

	var pipelines []*pipeline.Pipeline
	seen := make(map[string]string) // for deduplication: name -> file path

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

		// Check if file has a supported extension
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" && ext != ".alloy" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		var p pipeline.Pipeline
		if ext == ".alloy" {
			var err error
			p, err = parseAlloyMetadata(data)
			if err != nil {
				// This file is not a pipeline, so we can ignore it.
				return nil
			}
			p.Contents = string(data)
		} else {
			if err := yaml.Unmarshal(data, &p); err != nil {
				return fmt.Errorf("failed to parse pipeline from %s: %w", path, err)
			}
		}

		// If no name provided, use filename without extension
		if p.Name == "" {
			base := filepath.Base(path)
			p.Name = strings.TrimSuffix(base, filepath.Ext(base))
		}

		// Validate the pipeline configuration
		if err := p.Validate(); err != nil {
			return fmt.Errorf("invalid pipeline in %s: %w", path, err)
		}

		// Check for duplicate pipeline names
		if existingPath, ok := seen[p.Name]; ok {
			return fmt.Errorf("duplicate pipeline name '%s' found in %s and %s", p.Name, existingPath, path)
		}
		seen[p.Name] = path

		pipelines = append(pipelines, &p)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return pipelines, nil
}

// parseAlloyMetadata parses the fleet management metadata from an alloy file.
func parseAlloyMetadata(data []byte) (pipeline.Pipeline, error) {
	var p pipeline.Pipeline

	trimmedData := bytes.TrimSpace(data)
	if !bytes.HasPrefix(trimmedData, []byte(pipelineMetadataHeader)) {
		return p, fmt.Errorf("no fleet-management metadata block found at the beginning of the file")
	}

	headerEnd := bytes.Index(trimmedData, []byte(pipelineMetadataFooter))
	if headerEnd == -1 {
		return p, fmt.Errorf("could not find closing tag for fleet-management metadata block")
	}

	metadata := trimmedData[len(pipelineMetadataHeader):headerEnd]
	if err := yaml.Unmarshal(metadata, &p); err != nil {
		return p, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return p, nil
}
