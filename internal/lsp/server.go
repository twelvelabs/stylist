package lsp

import (
	"fmt"
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspSrv "github.com/tliron/glsp/server"

	"github.com/twelvelabs/stylist/internal/stylist"

	_ "github.com/tliron/commonlog/simple"
)

func NewServer(app *stylist.App) (*Server, error) {
	config, err := stylist.NewConfigFromPath(stylist.DefaultConfigPath)
	if err != nil {
		return nil, fmt.Errorf("language server: %w", err)
	}

	pipeline := stylist.NewPipeline(
		config.Processors,
		config.Excludes,
	)

	s := &Server{
		app:         app,
		config:      config,
		pipeline:    pipeline,
		documents:   NewDocumentStore(),
		diagnostics: NewDiagnosticService(app, pipeline),
	}
	s.server = s.newGlspServer()

	dir, _ := os.Getwd()
	s.server.Log.Infof("CURRENT WORKING DIR: %s", dir)

	return s, nil
}

type Server struct {
	app         *stylist.App
	config      *stylist.Config
	pipeline    *stylist.Pipeline
	server      *glspSrv.Server
	documents   *DocumentStore
	diagnostics *DiagnosticService
}

func (s *Server) RunNodeJs() error {
	if err := s.server.RunNodeJs(); err != nil {
		return fmt.Errorf("language server: %w", err)
	}
	return nil
}

func (s *Server) RunStdio() error {
	if err := s.server.RunStdio(); err != nil {
		return fmt.Errorf("language server: %w", err)
	}
	return nil
}

func (s *Server) RunTCP(address string) error {
	if err := s.server.RunTCP(address); err != nil {
		return fmt.Errorf("language server: %w", err)
	}
	return nil
}

func (s *Server) newGlspServer() *glspSrv.Server {
	name := "stylist"
	version := "0.0.1"

	// TODO: make configurable
	commonlog.Configure(2, nil)

	handler := protocol.Handler{}
	handler.Initialize = func(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
		if params.Trace != nil {
			protocol.SetTraceValue(*params.Trace)
		}
		capabilities := handler.CreateServerCapabilities()
		return protocol.InitializeResult{
			Capabilities: capabilities,
			ServerInfo: &protocol.InitializeResultServerInfo{
				Name:    name,
				Version: &version,
			},
		}, nil
	}
	handler.Initialized = func(context *glsp.Context, params *protocol.InitializedParams) error {
		return nil
	}
	handler.Shutdown = func(context *glsp.Context) error {
		protocol.SetTraceValue(protocol.TraceValueOff)
		return nil
	}
	handler.SetTrace = func(context *glsp.Context, params *protocol.SetTraceParams) error {
		protocol.SetTraceValue(params.Value)
		return nil
	}
	handler.TextDocumentDidOpen = func(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error { //nolint: lll
		uri := params.TextDocument.URI
		content := params.TextDocument.Text
		s.documents.Open(uri, content)
		return nil
	}
	handler.TextDocumentDidClose = func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error { //nolint: lll
		uri := params.TextDocument.URI
		s.documents.Close(uri)
		return nil
	}
	handler.TextDocumentDidChange = func(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error { //nolint: lll
		uri := params.TextDocument.URI
		doc, found := s.documents.Get(uri)
		if !found {
			return nil
		}
		doc.ApplyChanges(params.ContentChanges)

		doc.Debouncer.Run(func() {
			err := s.diagnostics.Publish(doc, context)
			if err != nil {
				s.server.Log.Error(err.Error())
			}
		})

		return nil
	}
	handler.TextDocumentDidSave = func(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error { //nolint: lll
		uri := params.TextDocument.URI
		doc, found := s.documents.Get(uri)
		if !found {
			return nil
		}

		doc.Debouncer.Run(func() {
			err := s.diagnostics.Publish(doc, context)
			if err != nil {
				s.server.Log.Error(err.Error())
			}
		})

		return nil
	}

	return glspSrv.NewServer(&handler, name, true)
}
