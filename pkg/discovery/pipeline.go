package discovery

import (
	fmapi "github.com/grafana/fleet-management-api/api/gen/proto/go/pipeline/v1"
)

// ConfigType identifies the format of a pipeline's contents.
type ConfigType string

const (
	ConfigTypeAlloy ConfigType = "alloy"
	ConfigTypeOTel  ConfigType = "otel"
)

type Pipeline struct {
	Name       string     `yaml:"name"`
	Contents   string     `yaml:"contents"`
	Matchers   []string   `yaml:"matchers"`
	Enabled    bool       `yaml:"enabled"`
	ConfigType ConfigType `yaml:"config_type"`
}

// ToFleetManagementAPI converts the local Pipeline struct to a fleet-management-api.Pipeline
func (p *Pipeline) ToFleetManagementPipeline() *fmapi.Pipeline {
	apiPipeline := &fmapi.Pipeline{
		Name:       p.Name,
		Contents:   p.Contents,
		Matchers:   p.Matchers,
		Enabled:    &p.Enabled,
		ConfigType: p.apiConfigType(),
	}

	return apiPipeline
}

// apiConfigType maps the pipeline's config type to the Fleet Management API enum.
// An unset or "alloy" type maps to Alloy; "otel" maps to OTel.
func (p *Pipeline) apiConfigType() fmapi.ConfigType {
	switch p.ConfigType {
	case ConfigTypeOTel:
		return fmapi.ConfigType_CONFIG_TYPE_OTEL
	default:
		return fmapi.ConfigType_CONFIG_TYPE_ALLOY
	}
}
