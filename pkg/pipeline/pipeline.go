package pipeline

import (
	"errors"
	"fmt"
	"os"

	fmapi "github.com/grafana/fleet-management-api/api/gen/proto/go/pipeline/v1"
)

// Validation errors
var (
	ErrBothContentFields = errors.New("cannot have both 'contents' and 'contents_file' fields")
	
	ErrNoContentFields = errors.New("must have either 'contents' or 'contents_file' field")
)

type Pipeline struct {
	Name         string   `yaml:"name"`
	Contents     string   `yaml:"contents"`
	ContentsFile string   `yaml:"contents_file"`
	Matchers     []string `yaml:"matchers"`
	Enabled      bool     `yaml:"enabled"`
}

// Validate checks that the Pipeline configuration is valid
func (p *Pipeline) Validate() error {
	// Check that contents and contents_file are mutually exclusive
	if p.Contents != "" && p.ContentsFile != "" {
		return fmt.Errorf("pipeline '%s': %w", p.Name, ErrBothContentFields)
	}
	
	// Check that at least one content source is provided
	if p.Contents == "" && p.ContentsFile == "" {
		return fmt.Errorf("pipeline '%s': %w", p.Name, ErrNoContentFields)
	}
	
	return nil
}

// ReadContents reads the contents from ContentsFile if specified.
// If ContentsFile is empty, this is a no-op.
func (p *Pipeline) ReadContents() error {
	if p.ContentsFile == "" {
		return nil
	}

	data, err := os.ReadFile(p.ContentsFile)
	if err != nil {
		return fmt.Errorf("failed to read contents file '%s' for pipeline '%s': %w", p.ContentsFile, p.Name, err)
	}

	p.Contents = string(data)
	return nil
}

// ToFleetManagementAPI converts the local Pipeline struct to a fleet-management-api.Pipeline
func (p *Pipeline) ToFleetManagementPipeline() *fmapi.Pipeline {
	apiPipeline := &fmapi.Pipeline{
		Name:     p.Name,
		Contents: p.Contents,
		Matchers: p.Matchers,
		Enabled:  &p.Enabled,
	}
	
	return apiPipeline
}
