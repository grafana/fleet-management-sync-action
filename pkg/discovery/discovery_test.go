package discovery

import (
	"context"
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
			wantPipelines: []string{"pipeline1", "pipeline2"},
		},
		{
			name:     "discover with non-existent root",
			rootPath: "testdata/non-existent",
			wantErr:  true,
		},
		{
			name:     "discover with typo in yaml",
			rootPath: "testdata/invalid_yaml_typo",
			wantErr:  true,
		},
		{
			name:     "discover otel pipeline with empty contents",
			rootPath: "testdata/invalid_otel_empty",
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

func TestDiscoverPipelinesConfigTypes(t *testing.T) {
	input := &config.Config{
		PipelinesRootPath: "testdata/valid_mixed_pipelines",
		Username:          "test",
		Token:             "test",
	}

	pipelines, err := FindPipelines(context.Background(), input)
	require.NoError(t, err)
	require.Len(t, pipelines, 2)

	byName := make(map[string]*Pipeline)
	for _, p := range pipelines {
		byName[p.Name] = p
	}

	alloy := byName["alloy-pipe"]
	require.NotNil(t, alloy)
	require.Equal(t, ConfigTypeAlloy, alloy.ConfigType)
	require.Contains(t, alloy.Contents, "logging")

	otel := byName["otel-pipe"]
	require.NotNil(t, otel)
	require.Equal(t, ConfigTypeOTel, otel.ConfigType)
	require.Contains(t, otel.Contents, "receivers:")
}




