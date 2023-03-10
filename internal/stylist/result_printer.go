package stylist

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/twelvelabs/termite/ioutil"
)

// ResultPrinter is the interface that wraps the Print method.
type ResultPrinter interface {
	// Print writes the results to Stdout.
	Print(results []*Result) error
}

// NewResultPrinter returns the appropriate printer for the given format.
func NewResultPrinter(ios *ioutil.IOStreams, config *Config) ResultPrinter { //nolint:ireturn
	format := config.Output.Format
	switch format {
	case ResultFormatSarif:
		return &SarifPrinter{ios: ios, config: config}
	case ResultFormatTty:
		return &TtyPrinter{ios: ios, config: config}
	default:
		panic(fmt.Sprintf("unknown result format: %s", format))
	}
}

/*
* SarifPrinter
**/

// SarifPrinter generates SARIF formatted output.
type SarifPrinter struct {
	ios    *ioutil.IOStreams
	config *Config
}

// Print writes the SARIF formatted results to Stdout.
func (p *SarifPrinter) Print(results []*Result) error {
	return nil
}

/*
* TtyPrinter
**/

// TtyPrinter generates TTY formatted output.
// The output will contain ANSI color codes if the terminal allows them.
type TtyPrinter struct {
	ios    *ioutil.IOStreams
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

func (p *TtyPrinter) printLocation(result *Result, formatter *ioutil.Formatter) {
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
		severity = ""
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
func (p *TtyPrinter) printUnderLinePointer(result *Result, formatter *ioutil.Formatter) {
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
