package stylist

import (
	"github.com/twelvelabs/termite/conf"
	"github.com/twelvelabs/termite/ioutil"
	"github.com/twelvelabs/termite/run"
	"github.com/twelvelabs/termite/ui"
)

type App struct {
	IO           *ioutil.IOStreams
	ConfigLoader *conf.Loader[*Config]
	Messenger    *ui.Messenger
	CmdClient    *run.Client
}

func NewApp() (*App, error) {
	ios := ioutil.System()
	loader := conf.NewLoader(&Config{}, ".stylist/stylist.yml")

	app := &App{
		IO:           ios,
		ConfigLoader: loader,
		Messenger:    ui.NewMessenger(ios),
		CmdClient:    run.NewClient(),
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
	}

	return app
}
