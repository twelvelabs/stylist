package stylist

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewConfig()
	assert.Equal(t, ".stylist.yml", config.ConfigPath)
	assert.Equal(t, LogLevelWarn, config.LogLevel)
	assert.Equal(t, ResultFormatTty, config.Output.Format)
	assert.Equal(t, true, config.Output.ShowContext)
}

func TestNewConfigFromArgs(t *testing.T) {
	validConfigPath := filepath.Join(".", "testdata", "config", "valid.yml")
	invalidConfigPath := filepath.Join(".", "testdata", "config", "invalid.yml")

	tests := []struct {
		desc        string
		args        []string
		expectation func(t *testing.T, config *Config)
		err         string
	}{
		{
			desc: "returns defaults when no args",
			args: []string{},
			expectation: func(t *testing.T, config *Config) {
				t.Helper()
				defaultConfig := NewConfig()
				assert.Equal(t, defaultConfig.ConfigPath, config.ConfigPath)
				assert.Equal(t, defaultConfig.LogLevel, config.LogLevel)
			},
		},
		{
			desc: "returns error when unable to parse args",
			args: []string{"---"},
			err:  "unable to parse",
		},
		{
			desc: "uses a custom config path if supplied",
			args: []string{
				"--config=" + validConfigPath,
			},
			expectation: func(t *testing.T, config *Config) {
				t.Helper()
				assert.Equal(t, validConfigPath, config.ConfigPath)
				assert.Equal(t, LogLevelDebug, config.LogLevel)
			},
		},
		{
			desc: "uses a custom log level if supplied",
			args: []string{
				"--config=" + validConfigPath,
				"--log-level=info",
			},
			expectation: func(t *testing.T, config *Config) {
				t.Helper()
				assert.Equal(t, validConfigPath, config.ConfigPath)
				assert.Equal(t, LogLevelInfo, config.LogLevel)
			},
		},
		{
			desc: "ignores unknown log level args",
			args: []string{
				"--log-level=nope",
			},
			expectation: func(t *testing.T, config *Config) {
				t.Helper()
				defaultConfig := NewConfig()
				assert.Equal(t, defaultConfig.ConfigPath, config.ConfigPath)
				assert.Equal(t, defaultConfig.LogLevel, config.LogLevel)
			},
		},
		{
			desc: "returns an error when invalid config file",
			args: []string{
				"--config=" + invalidConfigPath,
			},
			err: "yaml: unmarshal errors",
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

			if tt.expectation != nil {
				tt.expectation(t, actual)
			}
		})
	}
}
