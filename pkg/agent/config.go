package agent

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	emptyConfigErr       = "empty config"
	invalidFrequencyErr  = "invalid sample frequency"
	emptyOutputTargetErr = "target output file not specified"
)

// Config represents configuration supplied by user
type Config struct {
	// HttpSources holds all the http-sources
	HttpSources []HttpSource `yaml:"httpsources"`
	// SampleFrequency of the agent
	SampleFrequency int `yaml:"sampleFrequency"`
	// TargetOutputFile name
	TargetOutputFile string `yaml:"targetOutputFile"`
}

// HttpSource represents http-source configuration
type HttpSource struct {
	// Endpoints of htpp-source
	Endpoints string `yaml:"endpoints"`
	// Headers holds are http-headers supplied
	Headers map[string]string `yaml:"headers"`
}

// NewConfig reads the config file and return config object
func NewConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)

	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}
	if err = validateConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

// validateConfig supplied by user
func validateConfig(config *Config) error {
	if config == nil {
		return errors.New(emptyConfigErr)
	}
	if config.SampleFrequency < 1 {
		return errors.New(invalidFrequencyErr)
	}
	if config.TargetOutputFile == "" {
		return errors.New(emptyOutputTargetErr)
	}
	return nil
}
