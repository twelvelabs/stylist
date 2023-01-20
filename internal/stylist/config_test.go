package stylist

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twelvelabs/termite/testutil"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewConfig()
	assert.Equal(t, ".stylist.yml", config.ConfigPath)
	assert.Equal(t, LogLevelWarn, config.LogLevel)
	assert.Equal(t, ResultFormatTty, config.Output.Format)
	assert.Equal(t, true, config.Output.ShowContext)
}

func TestNewConfigFromArgs(t *testing.T) {
	configFixturePath := func(name string) string {
		return filepath.Join(".", "testdata", "config", name+".yml")
	}

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
				"--config=" + configFixturePath("valid"),
			},
			expectation: func(t *testing.T, config *Config) {
				t.Helper()
				assert.Equal(t, configFixturePath("valid"), config.ConfigPath)
				assert.Equal(t, LogLevelDebug, config.LogLevel)
			},
		},
		{
			desc: "uses a custom log level if supplied",
			args: []string{
				"--config=" + configFixturePath("valid"),
				"--log-level=info",
			},
			expectation: func(t *testing.T, config *Config) {
				t.Helper()
				assert.Equal(t, configFixturePath("valid"), config.ConfigPath)
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
				"--config=" + configFixturePath("invalid"),
			},
			err: "yaml: unmarshal errors",
		},
		{
			desc: "resolves processor presets",
			args: []string{
				"--config=" + configFixturePath("valid-presets"),
			},
			expectation: func(t *testing.T, config *Config) {
				t.Helper()
				assert.Equal(t, 2, len(config.Processors))
				assert.Equal(t, "golangci-lint", config.Processors[0].Name)
				assert.Equal(t, "markdownlint", config.Processors[1].Name)
			},
		},
		{
			desc: "returns an error when referencing an unknown preset name",
			args: []string{
				"--config=" + configFixturePath("invalid-presets"),
			},
			err: "unknown preset",
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

func TestWriteConfig(t *testing.T) {
	config := &Config{
		LogLevel: LogLevelWarn,
	}

	testutil.InTempDir(t, func(dir string) {
		configPath := filepath.Join(dir, "config.yml")
		assert.NoFileExists(t, configPath)

		err := WriteConfig(config, configPath)
		assert.NoError(t, err)
		assert.FileExists(t, configPath)

		configData, err := os.ReadFile(configPath)
		assert.NoError(t, err)
		assert.Equal(t, "log_level: warn\n", string(configData))
	})
}

func TestWriteConfig_WhenInvalidPath(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		configPath := filepath.Join(dir, "unknown", "sub-dir", "config.yml")
		assert.NoFileExists(t, configPath)

		err := WriteConfig(&Config{}, configPath)
		assert.ErrorContains(t, err, "no such file or directory")
		assert.NoFileExists(t, configPath)
	})
}

func TestCommentOutConfigPresets(t *testing.T) {
	uncommented, _ := os.ReadFile(filepath.Join("testdata", "config", "uncommented.yml"))
	commented, _ := os.ReadFile(filepath.Join("testdata", "config", "commented.yml"))

	testutil.InTempDir(t, func(dir string) {
		configPath := filepath.Join(dir, "config.yml")
		testutil.WriteFile(t, configPath, uncommented, 0600)

		err := CommentOutConfigPresets(configPath)
		assert.NoError(t, err)

		testutil.AssertFilePath(t, configPath, string(commented))
	})
}

func TestCommentOutConfigPresets_WhenInvalidPath(t *testing.T) {
	testutil.InTempDir(t, func(dir string) {
		configPath := filepath.Join(dir, "config.yml")
		assert.NoFileExists(t, configPath)

		err := CommentOutConfigPresets(configPath)
		assert.ErrorContains(t, err, "no such file or directory")
	})
}
