package agent

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		config   *Config
		expected error
	}{
		{
			desc: "valid config",
			config: &Config{
				HttpSources: []HttpSource{
					{
						Endpoints: "test.com/metrics",
						Headers:   map[string]string{"test-api-key": "api-key"},
					},
				},
				SampleFrequency:  1,
				TargetOutputFile: "test_output_file",
			},
			expected: nil,
		},
		{
			desc: "invalid sample frequency config",
			config: &Config{
				HttpSources: []HttpSource{
					{
						Endpoints: "test.com/metrics",
						Headers:   map[string]string{"test-api-key": "api-key"},
					},
				},
				SampleFrequency:  0,
				TargetOutputFile: "test_output_file",
			},
			expected: errors.New(invalidFrequencyErr),
		},
		{
			desc:     "empty config",
			config:   nil,
			expected: errors.New(emptyConfigErr),
		},
		{
			desc: "invalid config with no output file",
			config: &Config{
				HttpSources: []HttpSource{
					{
						Endpoints: "test.com/metrics",
						Headers:   map[string]string{"test-api-key": "api-key"},
					},
				},
				SampleFrequency: 1,
			},
			expected: errors.New(emptyOutputTargetErr),
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			err := validateConfig(test.config)
			assert.Equal(t, test.expected, err)

		})
	}
}
