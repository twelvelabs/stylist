package lsp

import (
	"context"
	"fmt"
	"os"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/twelvelabs/stylist/internal/stylist"
)

func NewDiagnosticService(app *stylist.App, pipeline *stylist.Pipeline) *DiagnosticService {
	return &DiagnosticService{
		app:      app,
		pipeline: pipeline,
	}
}

type DiagnosticService struct {
	app      *stylist.App
	pipeline *stylist.Pipeline
}

func (d *DiagnosticService) Calculate(doc *Document) error {
	cwd, _ := os.Getwd()
	// Execute the stylist pipeline.
	ctx := d.app.InitContext(context.Background())
	results, err := d.pipeline.Check(ctx, cwd, []string{doc.Path})
	if err != nil {
		return fmt.Errorf("diagnostic calculate: %w", err)
	}

	// Map the results to diagnostics
	diagnostics := []protocol.Diagnostic{}
	for _, result := range results {
		diagnostics = append(diagnostics, d.toDiagnostic(result))
	}
	doc.Diagnostics = diagnostics

	return nil
}

func (d *DiagnosticService) Publish(doc *Document, ctx *glsp.Context) error {
	if err := d.Calculate(doc); err != nil {
		return err
	}

	ctx.Notify(
		protocol.ServerTextDocumentPublishDiagnostics,
		&protocol.PublishDiagnosticsParams{
			URI:         doc.URI,
			Diagnostics: doc.Diagnostics,
		},
	)

	return nil
}

func (d *DiagnosticService) toSeverity(level stylist.ResultLevel) *protocol.DiagnosticSeverity {
	var severity protocol.DiagnosticSeverity
	switch level {
	case stylist.ResultLevelError:
		severity = protocol.DiagnosticSeverityError
	case stylist.ResultLevelWarning:
		severity = protocol.DiagnosticSeverityWarning
	case stylist.ResultLevelInfo:
		severity = protocol.DiagnosticSeverityInformation
	default:
		return nil
	}
	return &severity
}

func (d *DiagnosticService) toDiagnostic(result *stylist.Result) protocol.Diagnostic {
	diagnostic := protocol.Diagnostic{
		Range:    protocol.Range{},
		Severity: d.toSeverity(result.Level),
		Code: &protocol.IntegerOrString{
			Value: result.Rule.ID,
		},
		Source:  &result.Source,
		Message: result.Rule.Description,
	}

	if result.Location.StartLine > 0 {
		diagnostic.Range.Start.Line = uint32(result.Location.StartLine - 1)
	}
	if result.Location.StartColumn > 0 {
		diagnostic.Range.Start.Character = uint32(result.Location.StartColumn - 1)
	}

	if result.Location.EndLine > 0 {
		diagnostic.Range.End.Line = uint32(result.Location.EndLine - 1)
	} else {
		diagnostic.Range.End.Line = diagnostic.Range.Start.Line
	}
	if result.Location.EndColumn > 0 {
		diagnostic.Range.End.Character = uint32(result.Location.EndColumn - 1)
	} else {
		diagnostic.Range.End.Character = diagnostic.Range.Start.Character
	}

	if result.Rule.URI != "" {
		diagnostic.CodeDescription = &protocol.CodeDescription{
			HRef: result.Rule.URI,
		}
	}

	return diagnostic
}
