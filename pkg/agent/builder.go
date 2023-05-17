package agent

import (
	"errors"
	"github.com/sankalp-r/metrics-agent/pkg/metrics"
	"github.com/sankalp-r/metrics-agent/pkg/metrics/http"
	"github.com/sankalp-r/metrics-agent/pkg/metrics/system"
	"github.com/sankalp-r/metrics-agent/pkg/reporter"
	"time"
)

// Builder represents the agent-builder configuration
type Builder struct {
	config *Config
}

// NewBuilder creates new builder object
func NewBuilder(config *Config) *Builder {
	return &Builder{config: config}
}

// Build the agent using builder
func (b *Builder) Build() (Agent, error) {
	if b.config == nil {
		return nil, errors.New("config is nil")
	}
	metricsSources := make([]metrics.Source, 0)
	// append http-sources
	for _, httpSource := range b.config.HttpSources {
		metricsSources = append(metricsSources, http.NewHttpSource(httpSource.Endpoints, httpSource.Headers))
	}
	// append system-source
	metricsSources = append(metricsSources, system.NewSystemSource("algod"))
	fileReporter, err := reporter.NewReporter(b.config.TargetOutputFile)
	if err != nil {
		return nil, err
	}
	t := time.Duration(b.config.SampleFrequency) * time.Second
	return &minimalAgent{
		samplingDuration: t,
		metricsSources:   metricsSources,
		reporter:         fileReporter,
	}, nil
}
