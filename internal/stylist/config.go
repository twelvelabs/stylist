package stylist

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/twelvelabs/termite/conf"
)

const (
	DefaultConfigPath   string       = ".stylist/stylist.yml"
	DefaultLogLevel     LogLevel     = LogLevelWarn
	DefaultResultFormat ResultFormat = ResultFormatTty
)

type Config struct {
	ConfigPath string       `yaml:"config_path"`
	LogLevel   LogLevel     `yaml:"log_level"`
	Output     OutputConfig `yaml:"output"`

	Excludes   []string
	Processors []*Processor
}

type OutputConfig struct {
	Format ResultFormat `yaml:"format"`
}

func NewConfig() *Config {
	return &Config{
		ConfigPath: DefaultConfigPath,
		LogLevel:   DefaultLogLevel,
		Output: OutputConfig{
			Format: DefaultResultFormat,
		},
	}
}

// NewConfigFromArgs creates a new config from the given args.
// To do so, it uses a minimal, duplicate flag set to determine the path
// to the config file and the log level.
//
// The remaining flags are defined on the cobra.Command flag set. We need two
// different sets because most flags default to values from the config file,
// which we don't know the location of until we parse the config flag.
func NewConfigFromArgs(args []string) (*Config, error) {
	config := NewConfig()
	configPath := config.ConfigPath
	logLevelStr := ""

	fs := pflag.NewFlagSet("config-args", pflag.ContinueOnError)
	fs.StringVarP(&configPath, "config", "c", configPath, "")
	// Using a string flag because the enum type's flag.Value methods do input
	// validation automatically, and this conflicts w/ Cobra's completion logic.
	// Cobra knows when it's performing completion and can ignore the
	// validation error, but our duplicate flag set will not.
	fs.StringVar(&logLevelStr, "log-level", logLevelStr, "")

	// Ignore all the flags defined on the main Cobra flagset.
	fs.ParseErrorsWhitelist.UnknownFlags = true
	// Needed to suppress the default usage
	fs.Usage = func() {}

	err := fs.Parse(args)
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		return nil, fmt.Errorf("unable to parse config args: %w", err)
	}

	config, err = conf.NewLoader(config, configPath).Load()
	if err != nil {
		return nil, err
	}
	// Ensure the path from the flag takes precedence over anything in the file.
	config.ConfigPath = configPath

	if logLevelStr != "" {
		// Coerce the level string from the flag back into an enum.
		// We can safely ignore errors because the Cobra flag will catch them.
		parsed, err := ParseLogLevel(logLevelStr)
		if err == nil {
			config.LogLevel = parsed
		}
	}

	return config, nil
}
