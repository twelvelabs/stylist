package stylist

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/twelvelabs/termite/ioutil"
	"github.com/twelvelabs/termite/run"
	"github.com/twelvelabs/termite/ui"
)

type ctxKey string

const (
	ctxCmdClient ctxKey = "stylist.CmdClient"
	ctxConfig    ctxKey = "stylist.Config"
	ctxLogger    ctxKey = "stylist.Logger"
)

type App struct {
	IO        *ioutil.IOStreams
	Config    *Config
	Messenger *ui.Messenger
	Prompter  ui.Prompter
	CmdClient *run.Client
	Logger    *logrus.Logger
}

// InitContext returns a new context set with app values.
func (a *App) InitContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, ctxCmdClient, a.CmdClient)
	ctx = context.WithValue(ctx, ctxConfig, a.Config)
	ctx = context.WithValue(ctx, ctxLogger, a.Logger)
	return ctx
}

func AppCmdClient(ctx context.Context) *run.Client {
	return ctx.Value(ctxCmdClient).(*run.Client)
}

func AppConfig(ctx context.Context) *Config {
	return ctx.Value(ctxConfig).(*Config)
}

func AppLogger(ctx context.Context) *logrus.Logger {
	return ctx.Value(ctxLogger).(*logrus.Logger)
}

func NewApp() (*App, error) {
	startedAt := time.Now()

	config, err := NewConfigFromArgs(os.Args)
	if err != nil {
		return nil, err
	}

	ios := ioutil.System()
	logger := newLogger(ios, config.LogLevel)
	logger.Debugf(
		"Initializing app: config=%v log-level=%v",
		config.ConfigPath,
		config.LogLevel,
	)

	app := &App{
		IO:        ios,
		Config:    config,
		Messenger: ui.NewMessenger(ios),
		Prompter:  ui.NewSurveyPrompter(ios.In, ios.Out, ios.Err, ios),
		CmdClient: run.NewClient(),
		Logger:    logger,
	}

	logger.Debugf("Initialized app in %s", time.Since(startedAt))

	return app, nil
}

func NewTestApp() *App {
	config := NewConfig()

	ios := ioutil.Test()
	logger := newLogger(ios, LogLevelDebug)

	app := &App{
		IO:        ios,
		Config:    config,
		Messenger: ui.NewMessenger(ios),
		Prompter:  ui.NewPrompterMock(),
		CmdClient: run.NewClient().WithStubbing(),
		Logger:    logger,
	}

	return app
}

func newLogger(ios *ioutil.IOStreams, level LogLevel) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ios.Err)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:  ios.IsColorEnabled(),
		PadLevelText: true,
	})

	switch level {
	case LogLevelError:
		logger.SetLevel(logrus.ErrorLevel)
	case LogLevelWarn:
		logger.SetLevel(logrus.WarnLevel)
	case LogLevelInfo:
		logger.SetLevel(logrus.InfoLevel)
	case LogLevelDebug:
		logger.SetLevel(logrus.DebugLevel)
	default:
		panic(fmt.Sprintf("unknown log level: %v", level))
	}

	return logger
}
