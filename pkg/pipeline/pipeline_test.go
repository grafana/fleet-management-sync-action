package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadContents(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.alloy")
	testContent := `logging { level="info" }`
	
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err, "Failed to create test file")

	tests := []struct {
		name         string
		pipeline     Pipeline
		wantContents string
		wantErr      bool
	}{
		{
			name: "read contents from file",
			pipeline: Pipeline{
				Name:         "test-pipeline",
				ContentsFile: testFile,
			},
			wantContents: testContent,
			wantErr:      false,
		},
		{
			name: "no contents file specified",
			pipeline: Pipeline{
				Name:     "test-pipeline",
				Contents: "existing content",
			},
			wantContents: "existing content",
			wantErr:      false,
		},
		{
			name: "file does not exist",
			pipeline: Pipeline{
				Name:         "test-pipeline",
				ContentsFile: "/non/existent/file.alloy",
			},
			wantContents: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.pipeline
			err := p.ReadContents()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantContents, p.Contents)
		})
	}
}

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
			name: "valid with contents_file",
			pipeline: Pipeline{
				Name:         "test",
				ContentsFile: "/path/to/file",
			},
			wantErr: nil,
		},
		{
			name: "invalid with both contents and contents_file",
			pipeline: Pipeline{
				Name:         "test",
				Contents:     "some content",
				ContentsFile: "/path/to/file",
			},
			wantErr: ErrBothContentFields,
		},
		{
			name: "invalid with neither contents nor contents_file",
			pipeline: Pipeline{
				Name: "test",
			},
			wantErr: ErrNoContentFields,
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
