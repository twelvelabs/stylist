package stylist

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromArgs(t *testing.T) {
	validConfigPath := filepath.Join(".", "testdata", "config", "valid.yml")
	invalidConfigPath := filepath.Join(".", "testdata", "config", "invalid.yml")
	tests := []struct {
		desc     string
		args     []string
		expected *Config
		err      string
	}{
		{
			desc: "returns defaults when no args",
			args: []string{},
			expected: &Config{
				ConfigPath: DefaultConfigPath,
				LogLevel:   DefaultLogLevel,
			},
		},
		{
			desc:     "returns error when unable to parse args",
			args:     []string{"---"},
			expected: nil,
			err:      "unable to parse",
		},
		{
			desc: "uses a custom config path if supplied",
			args: []string{
				"--config=" + validConfigPath,
			},
			expected: &Config{
				ConfigPath: validConfigPath,
				LogLevel:   LogLevelDebug,
			},
		},
		{
			desc: "uses a custom log level if supplied",
			args: []string{
				"--config=" + validConfigPath,
				"--log-level=info",
			},
			expected: &Config{
				ConfigPath: validConfigPath,
				LogLevel:   LogLevelInfo,
			},
		},
		{
			desc: "ignores unknown log level args",
			args: []string{
				"--log-level=nope",
			},
			expected: &Config{
				ConfigPath: DefaultConfigPath,
				LogLevel:   DefaultLogLevel,
			},
		},
		{
			desc: "returns an error when invalid config file",
			args: []string{
				"--config=" + invalidConfigPath,
			},
			expected: nil,
			err:      "yaml: unmarshal errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual, err := NewConfigFromArgs(tt.args)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}

			assert.Equal(t, tt.expected, actual)
		})
	}
}
