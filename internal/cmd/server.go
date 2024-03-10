package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/twelvelabs/stylist/internal/lsp"
	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewServerCmd(app *stylist.App) *cobra.Command {
	action := NewServerAction(app)

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the Stylist language server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := action.Validate(args); err != nil {
				return err
			}
			return action.Run(cmd.Context())
		},
	}

	// LSP clients will use these flags to indicate the communication channel.
	cmd.Flags().StringVar(&action.ClientProcessID, "clientProcessId", action.ClientProcessID, "")
	cmd.Flags().BoolVar(&action.NodeIPC, "node-ipc", action.NodeIPC, "Use node IPC communication between the client and the server.") //nolint: lll
	cmd.Flags().IntVar(&action.SocketPort, "socket", action.SocketPort, "Use a TCP socket port as the communication channel.")        //nolint: lll
	cmd.Flags().BoolVar(&action.Stdio, "stdio", action.Stdio, "Use stdio as the communication channel.")                              //nolint: lll

	return cmd
}

func NewServerAction(app *stylist.App) *ServerAction {
	return &ServerAction{
		App: app,
	}
}

type ServerAction struct {
	*stylist.App

	Stdio           bool
	NodeIPC         bool
	SocketPort      int
	ClientProcessID string
}

func (a *ServerAction) Validate(_ []string) error {
	return nil
}

func (a *ServerAction) Run(_ context.Context) error {
	server, err := lsp.NewServer(a.App)
	if err != nil {
		return err
	}

	if a.Stdio { //nolint: gocritic
		return server.RunStdio()
	} else if a.NodeIPC {
		return server.RunNodeJs()
	} else if a.SocketPort > 0 {
		address := fmt.Sprintf(":%d", a.SocketPort)
		return server.RunTCP(address)
	}

	return errors.New("undefined communication channel")
}
