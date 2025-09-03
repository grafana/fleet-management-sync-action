package discovery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToFleetManagementPipeline(t *testing.T) {
	tests := []struct {
		name     string
		pipeline Pipeline
	}{
		{
			name: "with explicit name",
			pipeline: Pipeline{
				Name:     "my-pipeline",
				Contents: "config content",
				Matchers: []string{"env=prod"},
				Enabled:  true,
			},
		},
		{
			name: "without explicit name",
			pipeline: Pipeline{
				Contents: "config content",
				Matchers: []string{"env=dev"},
				Enabled:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.pipeline
			apiPipeline := p.ToFleetManagementPipeline()

			require.NotNil(t, apiPipeline)
			require.Equal(t, p.Name, apiPipeline.Name)
			require.Equal(t, p.Contents, apiPipeline.Contents)
			require.Equal(t, p.Matchers, apiPipeline.Matchers)
			require.NotNil(t, apiPipeline.Enabled)
			require.Equal(t, p.Enabled, *apiPipeline.Enabled)
		})
	}
}
