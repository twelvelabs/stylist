package stylist

import (
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
)

// NewContextLineAnalyzer returns a new ContextLineAnalyzer.
func NewContextLineAnalyzer() *ContextLineAnalyzer {
	return &ContextLineAnalyzer{}
}

// ContextLineAnalyzer analyzes context lines and returns relevant metadata.
type ContextLineAnalyzer struct {
}

// DetectLanguage analyzes the given context lines and returns
// the detected programming language.
func (s *ContextLineAnalyzer) DetectLanguage(path string, lines []string) string {
	l := lexers.Match(path)
	if l == nil {
		l = lexers.Analyse(strings.Join(lines, "\n")) // cspell:words Analyse
	}
	if l == nil {
		l = lexers.Get("plaintext")
	}
	return strings.ToLower(l.Config().Name)
}
