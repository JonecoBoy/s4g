# s4g Core Architecture & Extension Guide

This guide explains how `s4g` works internally and how you can extend it with new file types and data sources.

## 1. Core Files (The Heart of the SSG)

The main files are in `internal/core/`. They define the "contract" between where data comes from and where it goes.

### `internal/core/content.go`
This is the most important file. It defines the **`Content`** struct, which is the universal representation of a page within the system.
- It doesn't matter if the data comes from a Markdown file or a SQL Database; it will be converted into `Content`.
- It also defines the **interfaces** you must implement to extend the system:
    - **`DataSource`**: Defines how to fetch data (`Fetch`).
    - **`Renderer`**: Defines how to transform data into output files (`Render`).

### `internal/core/generator.go`
This is the "engine" that orchestrates everything.
1. It iterates through all registered `Sources` and calls `Fetch()`.
2. It receives a list of `[]Content`.
3. For each registered `Renderer`, it iterates through this list and calls `Render()` for each page.

---

## 2. How to create a new Renderer (e.g., RSS or JSON)

To create a new renderer, follow these steps:

### Step 1: Create the package
Create a new folder, for example `internal/renderer/json/`.

### Step 2: Implement the `Renderer` interface
Create a `renderer.go` file and implement the `Name()` and `Render()` methods.

```go
package json

import (
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    "github.com/user/s4g/internal/core"
)

type Renderer struct {
    OutputDir string
}

func (r *Renderer) Name() string { return "json" }

func (r *Renderer) Render(ctx context.Context, all []core.Content, c core.Content) error {
    outPath := filepath.Join(r.OutputDir, c.Slug + ".json")
    f, _ := os.Create(outPath)
    defer f.Close()
    
    return json.NewEncoder(f).Encode(c)
}
```

### Step 3: Register in `cmd/build.go`
Add the new renderer to the switch case inside the `RunE` function of `build.go`.

---

## 3. How to add a new Source (e.g., SQL or JSON)

### Step 1: Create the package
Create a folder like `internal/source/json/`.

### Step 2: Implement the `DataSource` interface
Create a `source.go` file and implement `Name()` and `Fetch()`.

```go
package json

import (
    "context"
    "github.com/user/s4g/internal/core"
)

type Source struct {
    FilePath string
}

func (s *Source) Name() string { return "json" }

func (s *Source) Fetch(ctx context.Context) ([]core.Content, error) {
    // 1. Read the JSON file
    // 2. Convert to a slice of core.Content
    // 3. Return the slice
    return []core.Content{
        {Title: "Example", Body: "Hello from JSON", Slug: "example-json"},
    }, nil
}
```

### Step 3: Register in `cmd/build.go`
Like with the renderer, add the logic to instantiate your new Source based on the configuration.

---

## Extension Flow Summary
1. **Define the data** (`Source`) -> Converts to `core.Content`.
2. **Define the output** (`Renderer`) -> Receives `core.Content` and writes to disk.
3. **Connect** in `cmd/build.go`.
