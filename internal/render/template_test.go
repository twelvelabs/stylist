package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestCompile(t *testing.T) {
	_, err := Compile(``)
	assert.NoError(t, err)

	_, err = Compile(`{{}`)
	assert.ErrorContains(t, err, `unexpected "}" in command`)
}

func TestMustCompile(t *testing.T) {
	_ = MustCompile(``)

	assert.Panics(t, func() {
		_ = MustCompile(`{{}`)
	})
}

func TestTemplate_Render(t *testing.T) {
	ts, _ := Compile(`Hello, {{ .Name }}`)
	rendered, err := ts.Render(map[string]any{
		"Name": "World",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World", rendered)

	ts, _ = Compile(`Hello, {{ fail "boom" }}`)
	rendered, err = ts.Render(nil)
	assert.ErrorContains(t, err, "fail: boom")
	assert.Equal(t, "", rendered)
}

func TestTemplate_MarshalText(t *testing.T) {
	ts, err := Compile(`Hello, {{ .Name }}`)
	assert.NoError(t, err)

	m, err := ts.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, `Hello, {{ .Name }}`, string(m))
}

func TestTemplate_UnmarshalText(t *testing.T) {
	ts := &Template{}
	err := ts.UnmarshalText([]byte(`Hello, {{ .Name }}`))
	assert.NoError(t, err)

	m, _ := ts.MarshalText()
	assert.Equal(t, "Hello, {{ .Name }}", string(m))

	rendered, err := ts.Render(map[string]any{
		"Name": "World",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World", rendered)
}

func TestTemplate_UnmarshalFromYAML(t *testing.T) {
	s := `greeting: "Hello, {{ .Name }}"`
	mapping := struct {
		Greeting *Template `yaml:"greeting"`
	}{}
	err := yaml.Unmarshal([]byte(s), &mapping)
	assert.NoError(t, err)

	rendered, err := mapping.Greeting.Render(map[string]any{
		"Name": "World",
	})
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World", rendered)
}
