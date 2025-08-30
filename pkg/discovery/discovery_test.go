package discovery

import (
	"context"
	"os"
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
			name:          "discover alloy pipeline",
			rootPath:      "testdata/valid_alloy_pipeline",
			wantCount:     1,
			wantPipelines: []string{"alloy_pipeline"},
		},
		{
			name:          "ignore alloy pipeline without metadata",
			rootPath:      "testdata/invalid_alloy_pipeline",
			wantCount:     0,
			wantPipelines: []string{},
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
				PipelinesRootPath: tt.rootPath,
				Username:          "test",
				Token:             "test",
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
