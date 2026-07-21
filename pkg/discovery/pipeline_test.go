package discovery

import (
	"testing"

	fmapi "github.com/grafana/fleet-management-api/api/gen/proto/go/pipeline/v1"
	"github.com/stretchr/testify/require"
)

func TestToFleetManagementPipeline(t *testing.T) {
	tests := []struct {
		name           string
		pipeline       Pipeline
		wantConfigType fmapi.ConfigType
	}{
		{
			name: "with explicit name",
			pipeline: Pipeline{
				Name:     "my-pipeline",
				Contents: "config content",
				Matchers: []string{"env=prod"},
				Enabled:  true,
			},
			wantConfigType: fmapi.ConfigType_CONFIG_TYPE_ALLOY,
		},
		{
			name: "without explicit name",
			pipeline: Pipeline{
				Contents: "config content",
				Matchers: []string{"env=dev"},
				Enabled:  false,
			},
			wantConfigType: fmapi.ConfigType_CONFIG_TYPE_ALLOY,
		},
		{
			name: "explicit alloy config type",
			pipeline: Pipeline{
				Name:       "alloy-pipeline",
				Contents:   "config content",
				ConfigType: ConfigTypeAlloy,
			},
			wantConfigType: fmapi.ConfigType_CONFIG_TYPE_ALLOY,
		},
		{
			name: "otel config type",
			pipeline: Pipeline{
				Name:       "otel-pipeline",
				Contents:   "receivers: {}",
				ConfigType: ConfigTypeOTel,
			},
			wantConfigType: fmapi.ConfigType_CONFIG_TYPE_OTEL,
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
			require.Equal(t, tt.wantConfigType, apiPipeline.ConfigType)
		})
	}
}
