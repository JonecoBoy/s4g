package markdown

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFetch(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "hello.md"), `---
title: Hello World
slug: hello
---
# Hello

This is a test.
`)

	writeFile(t, filepath.Join(dir, "no-fm.md"), `No front-matter here.`)

	src := New(dir)
	pages, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(pages))
	}

	// Find the "hello" page
	var helloPage *struct{ Slug, Title, Body string }
	for _, p := range pages {
		if p.Slug == "hello" {
			cp := struct{ Slug, Title, Body string }{p.Slug, p.Title, p.Body}
			helloPage = &cp
		}
	}
	if helloPage == nil {
		t.Fatal("expected page with slug 'hello'")
	}
	if helloPage.Title != "Hello World" {
		t.Errorf("title = %q, want %q", helloPage.Title, "Hello World")
	}
	if helloPage.Body == "" {
		t.Error("body should not be empty")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
