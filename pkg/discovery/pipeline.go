package discovery

import (
	fmapi "github.com/grafana/fleet-management-api/api/gen/proto/go/pipeline/v1"
)

type Pipeline struct {
	Name     string   `yaml:"name"`
	Contents string   `yaml:"-"`
	Matchers []string `yaml:"matchers"`
	Enabled  bool     `yaml:"enabled"`
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
