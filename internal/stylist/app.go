package stylist

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/twelvelabs/termite/conf"
	"github.com/twelvelabs/termite/ioutil"
	"github.com/twelvelabs/termite/run"
	"github.com/twelvelabs/termite/ui"
)

type CtxKey string

const (
	CtxLogger    CtxKey = "Logger"
	CtxCmdClient CtxKey = "CmdClient"
)

type App struct {
	IO           *ioutil.IOStreams
	ConfigLoader *conf.Loader[*Config]
	Messenger    *ui.Messenger
	CmdClient    *run.Client
	Logger       *logrus.Logger
}

// InitContext returns a new context set with app values.
func (a *App) InitContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, CtxCmdClient, a.CmdClient)
	ctx = context.WithValue(ctx, CtxLogger, a.Logger)
	return ctx
}

func (a *App) SetLogLevel(level logrus.Level) {
	if level >= logrus.TraceLevel {
		level = logrus.TraceLevel
	}
	a.Logger.SetLevel(level)
	a.Logger.Debug("Set log level to " + level.String())
}

func AppLogger(ctx context.Context) *logrus.Logger {
	return ctx.Value(CtxLogger).(*logrus.Logger)
}

func AppCmdClient(ctx context.Context) *run.Client {
	return ctx.Value(CtxCmdClient).(*run.Client)
}

func NewApp() (*App, error) {
	ios := ioutil.System()
	loader := conf.NewLoader(&Config{}, ".stylist/stylist.yml")

	app := &App{
		IO:           ios,
		ConfigLoader: loader,
		Messenger:    ui.NewMessenger(ios),
		CmdClient:    run.NewClient(),
		Logger:       newLogger(ios),
	}

	return app, nil
}

func NewTestApp() *App {
	ios := ioutil.System()

	// TODO: use fixture config file
	loader := conf.NewLoader(&Config{}, ".stylist/stylist.yml")

	app := &App{
		IO:           ios,
		ConfigLoader: loader,
		Messenger:    ui.NewMessenger(ios),
		CmdClient:    run.NewClient().WithStubbing(),
		Logger:       newLogger(ios),
	}

	return app
}

func newLogger(ios *ioutil.IOStreams) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ios.Err)
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:  true,
		PadLevelText: true,
	})
	return logger
}
