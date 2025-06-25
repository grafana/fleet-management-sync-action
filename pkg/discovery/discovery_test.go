package discovery

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/fleet-management-sync-action/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestDiscoverPipelines(t *testing.T) {
	tests := []struct {
		name          string
		rootPath      string
		wantCount     int
		wantPipelines []string
		wantErr       bool
	}{
		{
			name:          "discover pipelines with explicit root",
			rootPath:      "testdata/valid_pipelines",
			wantCount:     2,
			wantPipelines: []string{"test-pipeline-1", "pipeline2"},
		},
		{
			name:     "discover with non-existent root",
			rootPath: "testdata/non-existent",
			wantErr:  true,
		},
		{
			name:     "discover with invalid yaml",
			rootPath: "testdata/invalid_yaml",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &config.Config{
				RootPath: tt.rootPath,
				Username: "test",
				Token:    "test",
			}

			pipelines, err := FindPipelines(context.Background(), input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, pipelines, tt.wantCount)

			// Check pipeline names
			var names []string
			for _, p := range pipelines {
				names = append(names, p.Name)
			}
			require.ElementsMatch(t, tt.wantPipelines, names)
		})
	}
}

func TestDiscoverPipelinesDefaultRoot(t *testing.T) {
	// Save current directory and restore after test
	oldCwd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldCwd) //nolint:errcheck

	// Create temporary directory and change to it
	tmpDir := t.TempDir()
	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Create a pipeline in the temp directory
	yamlContent := `name: cwd-pipeline
contents: test`
	err = os.WriteFile("test.yaml", []byte(yamlContent), 0644)
	require.NoError(t, err)

	// Test with empty root path
	input := &config.Config{
		Username: "test",
		Token:    "test",
	}

	pipelines, err := FindPipelines(context.Background(), input)
	require.NoError(t, err)
	require.Len(t, pipelines, 1)
	require.Equal(t, "cwd-pipeline", pipelines[0].Name)
}

func TestResolveContentsPath(t *testing.T) {
	// Use testdata paths
	rootPath := "testdata/path_resolution"
	yamlPath := filepath.Join(rootPath, "pipelines", "test.yaml")

	// Get absolute paths for comparison
	absRoot, err := filepath.Abs(rootPath)
	require.NoError(t, err)
	absYAMLDir := filepath.Join(absRoot, "pipelines")

	relativeToYAML := filepath.Join(absYAMLDir, "config.alloy")
	relativeToRoot := filepath.Join(absRoot, "shared.alloy")

	tests := []struct {
		name         string
		contentsFile string
		yamlPath     string
		rootPath     string
		want         string
		wantErr      bool
	}{
		{
			name:         "absolute path exists",
			contentsFile: relativeToYAML,
			yamlPath:     yamlPath,
			rootPath:     absRoot,
			want:         relativeToYAML,
			wantErr:      false,
		},
		{
			name:         "relative to yaml",
			contentsFile: "config.alloy",
			yamlPath:     yamlPath,
			rootPath:     absRoot,
			want:         relativeToYAML,
			wantErr:      false,
		},
		{
			name:         "relative to root",
			contentsFile: "shared.alloy",
			yamlPath:     yamlPath,
			rootPath:     absRoot,
			want:         relativeToRoot,
			wantErr:      false,
		},
		{
			name:         "file not found",
			contentsFile: "nonexistent.alloy",
			yamlPath:     yamlPath,
			rootPath:     absRoot,
			want:         "",
			wantErr:      true,
		},
		{
			name:         "absolute path not found",
			contentsFile: "/absolute/not/found.alloy",
			yamlPath:     yamlPath,
			rootPath:     absRoot,
			want:         "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveContentsPath(tt.contentsFile, tt.yamlPath, tt.rootPath)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
