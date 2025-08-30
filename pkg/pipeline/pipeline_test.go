package pipeline

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		pipeline Pipeline
		wantErr  error
	}{
		{
			name: "valid with contents",
			pipeline: Pipeline{
				Name:     "test",
				Contents: "some content",
			},
			wantErr: nil,
		},
		{
			name: "invalid with no contents",
			pipeline: Pipeline{
				Name: "test",
			},
			wantErr: ErrNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pipeline.Validate()

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Contains(t, err.Error(), tt.pipeline.Name)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToFleetManagementPipeline(t *testing.T) {
	tests := []struct {
		name     string
		pipeline Pipeline
		wantName string
	}{
		{
			name: "with explicit name",
			pipeline: Pipeline{
				Name:     "my-pipeline",
				Contents: "config content",
				Matchers: []string{"env=prod"},
				Enabled:  true,
			},
			wantName: "my-pipeline",
		},
		{
			name: "without explicit name",
			pipeline: Pipeline{
				Contents: "config content",
				Matchers: []string{"env=dev"},
				Enabled:  false,
			},
			wantName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.pipeline
			apiPipeline := p.ToFleetManagementPipeline()

			require.NotNil(t, apiPipeline)
			require.Equal(t, tt.wantName, apiPipeline.Name)
			require.Equal(t, p.Contents, apiPipeline.Contents)
			require.Equal(t, p.Matchers, apiPipeline.Matchers)
			require.NotNil(t, apiPipeline.Enabled)
			require.Equal(t, p.Enabled, *apiPipeline.Enabled)
		})
	}
}
