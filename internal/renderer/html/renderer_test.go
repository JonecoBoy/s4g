package html

import (
	"context"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JonecoBoy/s4g/internal/core"
)

func TestRender(t *testing.T) {
	tmplDir := t.TempDir()
	outDir := t.TempDir()

	// Write minimal templates
	layout := `{{define "layout.html"}}<!DOCTYPE html>
<html><head><title>{{.SiteTitle}} - {{.Page.Title}}</title></head>
<body>{{template "content" .}}</body>
</html>{{end}}`

	page := `{{define "page.html"}}{{template "layout.html" .}}{{end}}
{{define "content"}}<h1>{{.Page.Title}}</h1>{{.Page.Body | safeHTML}}{{end}}`

	// parse manually for the test (safeHTML won't be registered, so use a simple page)
	simplePage := `{{define "page.html"}}<!DOCTYPE html>
<html><head><title>{{.SiteTitle}} | {{.Page.Title}}</title></head>
<body><h1>{{.Page.Title}}</h1></body>
</html>{{end}}`

	_ = layout
	_ = page

	if err := os.WriteFile(filepath.Join(tmplDir, "page.html"), []byte(simplePage), 0644); err != nil {
		t.Fatal(err)
	}

	r := New(tmplDir, outDir, "Test Site")
	if err := r.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	c := core.Content{Slug: "hello", Title: "Hello World", Body: "<p>hi</p>"}
	if err := r.Render(context.Background(), []core.Content{c}, c); err != nil {
		t.Fatalf("Render: %v", err)
	}

	outFile := filepath.Join(outDir, "hello.html")
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not found: %v", err)
	}

	got := string(data)
	if !strings.Contains(got, "Hello World") {
		t.Errorf("output does not contain page title; got:\n%s", got)
	}
	if !strings.Contains(got, "Test Site") {
		t.Errorf("output does not contain site title; got:\n%s", got)
	}
}

// keep template import used
var _ = template.New
