package stylist

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/twelvelabs/termite/conf"
	"github.com/twelvelabs/termite/ioutil"
	"github.com/twelvelabs/termite/run"
	"github.com/twelvelabs/termite/ui"
)

type ctxKey string

const (
	ctxCmdClient    ctxKey = "stylist.CmdClient"
	ctxConfigLoader ctxKey = "stylist.ConfigLoader"
	ctxLogger       ctxKey = "stylist.Logger"
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
	ctx = context.WithValue(ctx, ctxCmdClient, a.CmdClient)
	ctx = context.WithValue(ctx, ctxConfigLoader, a.ConfigLoader)
	ctx = context.WithValue(ctx, ctxLogger, a.Logger)
	return ctx
}

func AppCmdClient(ctx context.Context) *run.Client {
	return ctx.Value(ctxCmdClient).(*run.Client)
}

func AppConfigLoader(ctx context.Context) *conf.Loader[*Config] {
	return ctx.Value(ctxConfigLoader).(*conf.Loader[*Config])
}

func AppLogger(ctx context.Context) *logrus.Logger {
	return ctx.Value(ctxLogger).(*logrus.Logger)
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
	ios := ioutil.Test()

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
