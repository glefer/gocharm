# Architecture guide

Objectif: clarifier les responsabilités, stabiliser l’API publique et isoler les détails ANSI/mesure d’affichage.

## Constats actuels
- `core` concentre Console, Text/Markup, couleurs/ANSI, calcul de largeur, utilitaires ANSI.
- `components/table` rend un tableau et dépend de helpers ANSI/width exposés par `core`.
- L’interface `Renderable` expose `Render() string`, mais la table possède aussi un rendu streaming (`RenderTo(w)`), ce qui crée 2 approches concurrentes.
- Quelques détails d’implémentation sont exposés inutilement (`ANSIStyles`, `StringHeader`, `RenderableAdapter`).

## Cibles de refonte (idiomatique Go)
- API publique minimaliste et stable à la racine du module + sous‑packages clairs.
- Détails de bas niveau cachés (internal) et utilitaires nommés précisément.
- Rendu orienté streaming par défaut pour éviter les copies mémoire superflues, avec helpers pour récupérer une `string` si besoin.

### Layout proposé

```
github.com/glefer/gocharm
├── console/                # API haut niveau d’écriture (Console)
│   └── console.go
├── render/                 # Contrats de rendu (interfaces/helpers)
│   └── render.go           # Renderable (WriteTo), RenderString, etc.
├── text/                   # Texte simple + markup
│   ├── text.go             # New(text, color, styles...) impl WriteTo
│   └── markup.go           # Parseur markup (tags -> ANSI)
├── style/                  # Couleurs/Styles et colorize
│   └── style.go            # Color, Style, Colorize (map ANSI non exportée)
├── measure/                # Largeur affichée, padding, alignements
│   └── width.go            # RuneWidth, VisibleLen, Alignment, Pad
├── table/                  # Composant tableau (package de premier niveau)
│   ├── table.go            # Modèle + options (WithAlignments, WithPadding...)
│   ├── render.go           # Rendu (écrit via WriteTo)
│   └── border.go           # BorderStyle & presets
└── internal/
    └── ansi/               # Helpers ANSI non exportés
        ├── strip.go        # StripANSI, ExtractLeadingPrefix
        └── codes.go        # tables ANSI internes (si besoin)
```

Notes:
- `table` devient un package de premier niveau (`github.com/glefer/gocharm/table`).
- `ANSIStyles` devient interne; l’API publique expose seulement `style.Colorize(...)` et les types (`Color`, `Style`).
- Les types utilitaires aujourd’hui exportés mais propres au package (`StringHeader`, `RenderableAdapter`) deviennent non exportés.

## Contrats de rendu proposés

- Interface unique orientée streaming:
  - `type Renderable interface { WriteTo(w io.Writer) (int64, error) }`
- Helpers:
  - `func RenderString(r Renderable) string`
  - Optionnel: `type Stringer interface { String() string }` si besoin.
- Compat: conserver `Render() string` comme helper (non-interface), en déléguant à `RenderString`.

Bénéfices:
- Meilleure perf pour de grands tableaux/blocs
- Un seul contrat simplifie l’écosystème (Console, Table, Text partagent la même interface)

## Plan de migration incrémental (sans rupture)
1) Introduire le nouveau package `render` et l’interface `Renderable.WriteTo` + `RenderString`.
   - Implémenter `WriteTo` pour `Text` et `Table`.
   - Garder `Render() string` comme wrapper transitoire.
2) Extraire le parseur de markup dans `text/markup.go`.
   - `core.NewMarkup` reste un proxy temporaire (déprécié) vers `text.NewMarkup`.
3) Isoler ANSI et mesure:
   - Déplacer `StripANSI`, `ExtractLeadingPrefix` en `internal/ansi`.
   - Déplacer `VisibleLen`, `RuneWidth`, `Alignment`, `Pad` en `measure`.
   - Remplacer les imports dans `table`/`text` par `internal/ansi` et `measure`.
4) Cacher la table ANSI publique:
   - Rendre la map ANSI non exportée; exposer uniquement `style.Colorize` et types.
5) Nettoyage API:
   - Rendre `StringHeader`, `RenderableAdapter` non exportés.
   - Ajouter des options fonctionnelles cohérentes (ex: `table.WithBorder`, `table.WithPadding`).
6) Déprécations & alias:
   - Re‑export minimal depuis `core` ou la racine si nécessaire, avec commentaires `// Deprecated:`.
   - Marquer dans le README le chemin de migration recommandé.

## Lignes directrices
- Garder les packages petits, avec responsabilités nettes.
- Préférer des fonctions/constructeurs clairs et des types simples.
- Exporter le strict nécessaire; tout le reste va dans `internal/`.
- Tests: couvrir le parseur de markup, le calcul de largeur (unicode, emojis), et des golden tests pour `table`.

## Quick wins immédiats
- Ajuster `go.mod` vers une version Go existante (ex: `go 1.22` ou `1.23`).
- Renommer en non exporté: `StringHeader` -> `stringHeader`, `RenderableAdapter` -> `renderableAdapter`.
- Ajouter des commentaires de package (`// Package table ...`) et des godocs sur les options.
- Ajouter un test pour `ExtractLeadingPrefix` (cas avec multiples séquences ANSI).

## Exemple d’usage cible
```go
import (
    "github.com/glefer/gocharm/console"
    "github.com/glefer/gocharm/style"
    "github.com/glefer/gocharm/table"
)

c := console.New()
c.Println("[green bold]Hello[/]")

// table
T := table.New("Name", "City").WithPadding(2)
T.AddRowVar("Alice", "Paris")
c.Render(T) // via WriteTo/RenderString
```

---
Ce document est une proposition pragmatique, pensée pour une migration par petites PRs sans casser l’existant.