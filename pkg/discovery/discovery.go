package discovery

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/fleet-management-sync-action/pkg/config"
	"github.com/grafana/fleet-management-sync-action/pkg/pipeline"
	"gopkg.in/yaml.v3"
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

		// Check if file has .yaml or .yml extension
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		var p pipeline.Pipeline
		if err := yaml.Unmarshal(data, &p); err != nil {
			return fmt.Errorf("failed to parse pipeline from %s: %w", path, err)
		}

		// If no name provided, use filename without extension
		if p.Name == "" {
			base := filepath.Base(path)
			p.Name = strings.TrimSuffix(base, filepath.Ext(base))
		}

		// Validate the pipeline configuration before reading contents
		if err := p.Validate(); err != nil {
			return fmt.Errorf("invalid pipeline in %s: %w", path, err)
		}

		// Resolve contents_file path if specified, this is done for consistency
		if p.ContentsFile != "" {
			resolvedPath, err := resolveContentsPath(p.ContentsFile, path, rootPath)
			if err != nil {
				return fmt.Errorf("failed to resolve contents_file for pipeline in %s: %w", path, err)
			}
			p.ContentsFile = resolvedPath
		}

		// Check for duplicate pipeline names
		if existingPath, ok := seen[p.Name]; ok {
			return fmt.Errorf("duplicate pipeline name '%s' found in %s and %s", p.Name, existingPath, path)
		}
		seen[p.Name] = path

		// Load contents from ContentsFile, if specified
		if err := p.ReadContents(); err != nil {
			return fmt.Errorf("failed to read contents for pipeline in %s: %w", path, err)
		}

		pipelines = append(pipelines, &p)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return pipelines, nil
}

// resolveContentsPath resolves a contents_file path by trying:
// 1. If absolute, return as-is
// 2. Relative to the YAML file's directory
// 3. Relative to the root path
// Returns error if the file doesn't exist in any location
func resolveContentsPath(contentsFile, yamlPath, rootPath string) (string, error) {
	// If already absolute, check if it exists
	if filepath.IsAbs(contentsFile) {
		if _, err := os.Stat(contentsFile); err != nil {
			return "", fmt.Errorf("contents file not found at absolute path %s: %w", contentsFile, err)
		}
		return contentsFile, nil
	}

	// Try relative to YAML file's directory
	yamlDir := filepath.Dir(yamlPath)
	relativeToYAML := filepath.Join(yamlDir, contentsFile)
	if _, err := os.Stat(relativeToYAML); err == nil {
		return filepath.Abs(relativeToYAML)
	}

	// Try relative to root path
	relativeToRoot := filepath.Join(rootPath, contentsFile)
	if _, err := os.Stat(relativeToRoot); err == nil {
		return filepath.Abs(relativeToRoot)
	}

	return "", fmt.Errorf("contents file %s not found relative to YAML (%s) or root (%s)",
		contentsFile, yamlDir, rootPath)
}
