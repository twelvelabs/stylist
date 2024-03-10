package lsp

import (
	"strings"
	"time"

	protocol "github.com/tliron/glsp/protocol_3_16"
	"go.lsp.dev/uri"
)

type Document struct {
	URI         string
	Path        string
	Content     string
	Debouncer   *Debouncer
	Diagnostics []protocol.Diagnostic

	lines []string
}

// ApplyChanges updates the content of the document from LSP textDocument/didChange events.
func (d *Document) ApplyChanges(events []any) {
	for _, event := range events {
		switch c := event.(type) {
		case protocol.TextDocumentContentChangeEvent:
			startIndex, endIndex := c.Range.IndexesIn(d.Content)
			d.Content = d.Content[:startIndex] + c.Text + d.Content[endIndex:]
		case protocol.TextDocumentContentChangeEventWhole:
			d.Content = c.Text
		}
	}
	d.lines = nil
}

// ContentAtRange returns the document text at given range.
func (d *Document) ContentAtRange(r protocol.Range) string {
	startIndex, endIndex := r.IndexesIn(d.Content)
	return d.Content[startIndex:endIndex]
}

// GetLine returns the line at the given index.
func (d *Document) GetLine(index int) (string, bool) {
	lines := d.GetLines()
	if index < 0 || index > len(lines) {
		return "", false
	}
	return lines[index], true
}

// GetLines returns all the lines in the document.
func (d *Document) GetLines() []string {
	if d.lines == nil {
		d.lines = strings.Split(d.Content, "\n")
	}
	return d.lines
}

// WordAt returns the word found at the given location.
func (d *Document) WordAt(pos protocol.Position) string {
	line, ok := d.GetLine(int(pos.Line))
	if !ok {
		return ""
	}
	char := int(pos.Character)
	words := strings.Fields(line[char:])
	if len(words) > char {
		return ""
	}
	return words[char]
}

func NewDocumentStore() *DocumentStore {
	return &DocumentStore{
		documents: map[string]*Document{},
	}
}

type DocumentStore struct {
	documents map[string]*Document
}

func (s *DocumentStore) Open(pathOrURI string, content string) *Document {
	uri := s.toURI(pathOrURI)
	path := s.toPath(pathOrURI)
	doc := &Document{
		URI:       uri,
		Path:      path,
		Content:   content,
		Debouncer: NewDebouncer(1 * time.Second),
	}
	s.documents[path] = doc
	return doc
}

func (s *DocumentStore) Close(pathOrURI string) {
	path := s.toPath(pathOrURI)
	delete(s.documents, path)
}

func (s *DocumentStore) Get(pathOrURI string) (*Document, bool) {
	path := s.toPath(pathOrURI)
	d, ok := s.documents[path]
	return d, ok
}

func (s *DocumentStore) toURI(pathOrURI string) string {
	return string(uri.New(pathOrURI))
}

func (s *DocumentStore) toPath(pathOrURI string) string {
	return uri.New(pathOrURI).Filename()
}
