# Guide d'Architecture s4g et Extension

Ce guide explique le fonctionnement interne de `s4g` et comment vous pouvez l'étendre avec de nouveaux types de fichiers et sources de données.

## 1. Fichiers Core (Le Cœur du SSG)

Les fichiers principaux se trouvent dans `internal/core/`. Ils définissent le « contrat » entre la provenance des données et leur destination.

### `internal/core/content.go`
C'est le fichier le plus important. Il définit la structure **`Content`**, qui est la représentation universelle d'une page au sein du système.
- Peu importe si les données proviennent d'un fichier Markdown ou d'une base de données SQL ; elles seront converties en `Content`.
- Il définit également les **interfaces** que vous devez implémenter pour étendre le système :
    - **`DataSource`** : Définit comment récupérer les données (`Fetch`).
    - **`Renderer`** : Définit comment transformer les données en fichiers de sortie (`Render`).

### `internal/core/generator.go`
C'est le « moteur » qui orchestre tout.
1. Il parcourt toutes les `Sources` enregistrées et appelle `Fetch()`.
2. Il reçoit une liste de `[]Content`.
3. Pour chaque `Renderer` enregistré, il parcourt cette liste et appelle `Render()` pour chaque page.

---

## 2. Comment créer un nouveau Renderer (ex. : RSS ou JSON)

Pour créer un nouveau moteur de rendu, suivez ces étapes :

### Étape 1 : Créer le paquet
Créez un nouveau dossier, par exemple `internal/renderer/json/`.

### Étape 2 : Implémenter l'interface `Renderer`
Créez un fichier `renderer.go` et implémentez les méthodes `Name()` et `Render()`.

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

### Étape 3 : S'enregistrer dans `cmd/build.go`
Ajoutez le nouveau moteur de rendu au switch case à l'intérieur de la fonction `RunE` de `build.go`.

---

## 3. Comment ajouter une nouvelle Source (ex. : SQL ou JSON)

### Étape 1 : Créer le paquet
Créez un dossier comme `internal/source/json/`.

### Étape 2 : Implémenter l'interface `DataSource`
Créez un fichier `source.go` et implémentez `Name()` et `Fetch()`.

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
    // 1. Lire le fichier JSON
    // 2. Convertir en une slice de core.Content
    // 3. Retourner la slice
    return []core.Content{
        {Title: "Exemple", Body: "Bonjour du JSON", Slug: "exemple-json"},
    }, nil
}
```

### Étape 3 : S'enregistrer dans `cmd/build.go`
Comme pour le moteur de rendu, ajoutez la logique pour instancier votre nouvelle Source en fonction de la configuration.

---

## Résumé du Flux d'Extension
1. **Définir la donnée** (`Source`) -> Conversion en `core.Content`.
2. **Définir la sortie** (`Renderer`) -> Reçoit `core.Content` et écrit sur le disque.
3. **Connecter** dans `cmd/build.go`.
