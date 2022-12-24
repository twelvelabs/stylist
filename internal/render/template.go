package render

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// Template is a lightweight wrapper for template strings.
//
// The underlying [text/template.Template] struct has access to
// functions from the [github.com/Masterminds/sprig] package,
// and can be invoked via the `Render()` method.
//
// Since it implements the [encoding.TextMarshaler] and
// [encoding.TextUnmarshaler] interfaces, it can be used with
// JSON or YAML fields containing template strings.
type Template struct {
	s string
	t *template.Template
}

// Compile parses a template string and returns, if successful,
// a new Template that can be rendered.
func Compile(s string) (*Template, error) {
	t, err := template.New("render.Template").Funcs(sprig.FuncMap()).Parse(s)
	if err != nil {
		return &Template{}, err
	}
	return &Template{s: s, t: t}, nil
}

func MustCompile(s string) *Template {
	t, err := Compile(s)
	if err != nil {
		panic(err)
	}
	return t
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (ts *Template) MarshalText() ([]byte, error) {
	return []byte(ts.s), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (ts *Template) UnmarshalText(text []byte) error {
	tsp, err := Compile(string(text))
	*ts = *tsp
	return err
}

// Render renders the template using data.
func (ts *Template) Render(data map[string]any) (string, error) {
	buf := bytes.Buffer{}
	err := ts.t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
