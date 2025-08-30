package pipeline

import (
	"errors"
	"fmt"

	fmapi "github.com/grafana/fleet-management-api/api/gen/proto/go/pipeline/v1"
)

// Validation errors
var (
	ErrNoContent = errors.New("must have 'contents' field")
)

type Pipeline struct {
	Name     string   `yaml:"name"`
	Contents string   `yaml:"contents"`
	Matchers []string `yaml:"matchers"`
	Enabled  bool     `yaml:"enabled"`
}

// Validate checks that the Pipeline configuration is valid
func (p *Pipeline) Validate() error {
	if p.Contents == "" {
		return fmt.Errorf("pipeline '%s': %w", p.Name, ErrNoContent)
	}

	return nil
}

// ReadContents is a no-op and is kept for backwards compatibility with the discovery process.
func (p *Pipeline) ReadContents() error {
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
