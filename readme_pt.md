# s4g Core Architecture & Extension Guide

Esta guia explica como o `s4g` funciona internamente e como você pode estendê-lo com novos tipos de arquivos e fontes de dados.

## 1. Arquivos Core (O Coração do SSG)

Os arquivos principais estão em `internal/core/`. Eles definem o "contrato" entre de onde os dados vêm e para onde eles vão.

### `internal/core/content.go`
Este é o arquivo mais importante. Ele define a estrutura **`Content`**, que é a representação universal de uma página dentro do sistema.
- Não importa se o dado veio de um Markdown ou de um Banco de Dados SQL, ele será convertido para um `Content`.
- Ele também define as **interfaces** que você deve implementar para estender o sistema:
    - **`DataSource`**: Define como buscar dados (`Fetch`).
    - **`Renderer`**: Define como transformar dados em arquivos de saída (`Render`).

### `internal/core/generator.go`
Este é o "motor" que orquestra tudo.
1. Ele percorre todas as `Sources` cadastradas e chama `Fetch()`.
2. Recebe uma lista de `[]Content`.
3. Para cada `Renderer` cadastrado, ele percorre essa lista e chama `Render()` para cada página.

---

## 2. Como criar um novo Renderer (Ex: RSS ou JSON)

Para criar um novo renderer, siga estes passos:

### Passo 1: Criar o pacote
Crie uma nova pasta, por exemplo `internal/renderer/json/`.

### Passo 2: Implementar a interface `Renderer`
Crie um arquivo `renderer.go` e implemente os métodos `Name()` e `Render()`.

```go
package json

import (
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    "github.com/JonecoBoy/s4g/internal/core"
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

### Passo 3: Registrar no `cmd/build.go`
Adicione o novo renderer no switch case dentro da função `RunE` do `build.go`.

---

## 3. Como adicionar um novo Source (Ex: SQL ou JSON)

### Passo 1: Criar o pacote
Crie uma pasta como `internal/source/json/`.

### Passo 2: Implementar a interface `DataSource`
Crie um arquivo `source.go` e implemente `Name()` e `Fetch()`.

```go
package json

import (
    "context"
    "github.com/JonecoBoy/s4g/internal/core"
)

type Source struct {
    FilePath string
}

func (s *Source) Name() string { return "json" }

func (s *Source) Fetch(ctx context.Context) ([]core.Content, error) {
    // 1. Leia o arquivo JSON
    // 2. Converta para uma slice de core.Content
    // 3. Retorne a slice
    return []core.Content{
        {Title: "Exemplo", Body: "Olá do JSON", Slug: "exemplo-json"},
    }, nil
}
```

### Passo 3: Registrar no `cmd/build.go`
Assim como no renderer, adicione a lógica para instanciar seu novo Source baseada na configuração.

---

## Resumo do Fluxo de Extensão
1. **Defina o dado** (`Source`) -> Converte para `core.Content`.
2. **Defina a saída** (`Renderer`) -> Recebe `core.Content` e grava no disco.
3. **Conecte** no `cmd/build.go`.
