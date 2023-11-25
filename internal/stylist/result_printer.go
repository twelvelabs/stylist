package stylist

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/twelvelabs/termite/ui"

	"github.com/twelvelabs/stylist/internal/checkstyle"
)

// ResultPrinter is the interface that wraps the Print method.
type ResultPrinter interface {
	// Print writes the results to Stdout.
	Print(results []*Result) error
}

// NewResultPrinter returns the appropriate printer for the given format.
func NewResultPrinter(ios *ui.IOStreams, config *Config) ResultPrinter { //nolint:ireturn
	format := config.Output.Format
	switch format {
	case ResultFormatCheckstyle:
		return &CheckstylePrinter{ios: ios, config: config}
	case ResultFormatJson:
		return &JSONPrinter{ios: ios, config: config}
	case ResultFormatSarif:
		return &SarifPrinter{ios: ios, config: config}
	case ResultFormatTty:
		return &TtyPrinter{ios: ios, config: config}
	default:
		panic(fmt.Sprintf("unknown result format: %s", format))
	}
}

/*
* CheckstylePrinter
**/

// CheckstylePrinter generates Checkstyle formatted output.
type CheckstylePrinter struct {
	ios    *ui.IOStreams
	config *Config
}

// Print writes the Checkstyle formatted results to Stdout.
func (p *CheckstylePrinter) Print(results []*Result) error {
	files := map[string]*checkstyle.CSFile{}
	paths := []string{}

	for _, result := range results {
		path := result.Location.Path

		if _, ok := files[path]; !ok {
			files[path] = &checkstyle.CSFile{
				Name: path,
			}
			paths = append(paths, path)
		}

		files[path].Errors = append(files[path].Errors, &checkstyle.CSError{
			Column:   result.Location.StartColumn,
			Line:     result.Location.StartLine,
			Message:  fmt.Sprintf("%s [%s]", result.Rule.Description, result.Rule.ID),
			Severity: result.Level.String(),
			Source:   result.Source,
		})
	}

	csr := &checkstyle.CSResult{Version: "4.3"}
	for _, path := range paths {
		csr.Files = append(csr.Files, files[path])
	}

	buf, err := xml.Marshal(csr)
	if err != nil {
		return err
	}

	doc := xml.Header + string(buf) + "\n"
	_, err = fmt.Fprint(p.ios.Out, doc)
	return err
}

/*
* JSONPrinter
**/

// JSONPrinter generates JSON formatted output.
type JSONPrinter struct {
	ios    *ui.IOStreams
	config *Config
}

// Print writes the JSON formatted results to Stdout.
func (p *JSONPrinter) Print(results []*Result) error {
	if !p.config.Output.ShowContext {
		for _, r := range results {
			r.ContextLines = nil
			r.ContextLang = ""
		}
	}

	buf, err := json.Marshal(results)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(p.ios.Out, string(buf)+"\n")
	return err
}

/*
* SarifPrinter
**/

// SarifPrinter generates SARIF formatted output.
type SarifPrinter struct {
	ios    *ui.IOStreams
	config *Config
}

// Print writes the SARIF formatted results to Stdout.
func (p *SarifPrinter) Print(results []*Result) error {
	// Group results by source so we can create a SARIF run for each.
	resultsBySource := map[string][]*Result{}
	for _, r := range results {
		if _, ok := resultsBySource[r.Source]; !ok {
			resultsBySource[r.Source] = []*Result{}
		}
		resultsBySource[r.Source] = append(resultsBySource[r.Source], r)
	}

	// create a top-level SARIF report.
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return err
	}

	// Create a SARIF run for each source.
	for sourceName, resultsForSource := range resultsBySource {
		run := sarif.NewRun(*sarif.NewSimpleTool(sourceName))
		for _, r := range resultsForSource {
			// Create a new rule for each rule id.
			rule := run.AddRule(r.Rule.ID).
				WithName(r.Rule.Name)
			if r.Rule.URI != "" {
				rule.WithHelpURI(r.Rule.URI)
			}

			// Add the location as a unique artifact for the entire run.
			artifact := run.AddDistinctArtifact(r.Location.Path)
			if p.config.Output.ShowContext {
				artifact.WithSourceLanguage(r.ContextLang)
			}

			// Create a region w/ the location and range for this particular result.
			region := sarif.NewRegion().
				WithStartLine(r.Location.StartLine).
				WithStartColumn(r.Location.StartColumn).
				WithEndLine(r.Location.EndLine).
				WithEndColumn(r.Location.EndColumn).
				WithSourceLanguage(r.ContextLang)
			if p.config.Output.ShowContext {
				region = region.WithSourceLanguage(r.ContextLang)
				if len(r.ContextLines) > 0 {
					region = region.WithSnippet(
						sarif.NewArtifactContent().WithText(
							strings.Join(r.ContextLines, "\n"),
						),
					)
				}
			}

			// Add the result to the run.
			run.CreateResultForRule(r.Rule.ID).
				WithLevel(r.Level.String()).
				WithMessage(sarif.NewTextMessage(r.Rule.Description)).
				AddLocation(
					sarif.NewLocationWithPhysicalLocation(
						sarif.NewPhysicalLocation().
							WithArtifactLocation(
								sarif.NewSimpleArtifactLocation(r.Location.Path),
							).
							WithRegion(region),
					),
				)
		}
		report.AddRun(run)
	}

	return report.Write(p.ios.Out)
}

/*
* TtyPrinter
**/

// TtyPrinter generates TTY formatted output.
// The output will contain ANSI color codes if the terminal allows them.
type TtyPrinter struct {
	ios    *ui.IOStreams
	config *Config
}

// Print writes the TTY formatted results to Stdout.
func (p *TtyPrinter) Print(results []*Result) error {
	formatter := p.ios.Formatter()
	for _, result := range results {
		p.printLocation(result, formatter)
		if p.config.Output.ShowContext {
			p.printContext(result)
			p.printUnderLinePointer(result, formatter)
		}
	}
	return nil
}

func (p *TtyPrinter) printLocation(result *Result, formatter *ui.Formatter) {
	severity := result.Level.String()
	switch result.Level {
	case ResultLevelError:
		severity = formatter.Red(severity) + ": "
	case ResultLevelWarning:
		severity = formatter.Yellow(severity) + ": "
	case ResultLevelInfo:
		severity = formatter.Cyan(severity) + ": "
	case ResultLevelNone:
		severity = formatter.Gray(severity) + ": "
	default:
	}

	source := result.Source
	if source != "" {
		source = formatter.Underline(source) + ": "
	}

	msg := result.Rule.Description
	if msg != "" {
		if !(strings.HasSuffix(msg, ".") || strings.HasSuffix(msg, "!")) {
			msg += "."
		}
		msg += " "
	}

	rule := ""
	if result.Rule.ID != "" {
		rule = fmt.Sprintf("[%s]", result.Rule.ID)
	}
	if result.Rule.URI != "" && p.config.Output.ShowURL {
		rule = fmt.Sprintf("%s(%s)", rule, result.Rule.URI)
	}

	fmt.Fprintf(
		p.ios.Out,
		"%s: %s%s%s%s\n",
		formatter.Bold(result.Location.String()),
		severity,
		source,
		msg,
		rule,
	)
}

func (p *TtyPrinter) printContext(result *Result) {
	if len(result.ContextLines) == 0 {
		return
	}

	contextLines := strings.Join(result.ContextLines, "\n") + "\n"
	if p.config.Output.SyntaxHighlight && p.ios.IsColorEnabled() {
		contextLines, _ = p.syntaxHighlight(
			contextLines, result.Location.Path, result.ContextLang,
		)
	}

	fmt.Fprint(p.ios.Out, contextLines)
}

func (p *TtyPrinter) syntaxHighlight(text, path, lang string) (string, error) {
	// seems the chroma author uses UK english :/
	// cspell:words Analyse Tokenise

	// Resolve the lexer.
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Match(path)
	}
	if l == nil {
		l = lexers.Analyse(text)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Resolve the formatter.
	f := formatters.TTY256

	// Resolve the style.
	s := styles.Get("emacs")
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, text)
	if err != nil {
		return text, err
	}

	buf := &bytes.Buffer{}
	err = f.Format(buf, s, it)
	if err != nil {
		return text, err
	}

	return buf.String(), nil
}

// Copied from golangci-lint.
func (p *TtyPrinter) printUnderLinePointer(result *Result, formatter *ui.Formatter) {
	// StartColumn == 0 means "unknown".
	if len(result.ContextLines) != 1 || result.Location.StartColumn == 0 {
		return
	}

	startCol0 := result.Location.StartColumn - 1
	endCol0 := result.Location.EndColumn - 1
	line := result.ContextLines[0]

	prefixRunes := make([]rune, 0, len(line))
	for j := 0; j < len(line) && j < startCol0; j++ {
		if line[j] == '\t' {
			prefixRunes = append(prefixRunes, '\t')
		} else {
			prefixRunes = append(prefixRunes, ' ')
		}
	}

	indicatorCols := 1
	if endCol0 > startCol0 && endCol0 <= len(line) {
		indicatorCols = endCol0 - startCol0
	}
	indicatorRunes := make([]rune, 0, indicatorCols)
	for j := 0; j < indicatorCols; j++ {
		indicatorRunes = append(indicatorRunes, '^')
	}

	fmt.Fprintf(p.ios.Out, "%s%s\n", string(prefixRunes), formatter.Yellow(string(indicatorRunes)))
}
