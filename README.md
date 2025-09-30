# gocharm

Affichage enrichi pour la console en Go : couleurs, styles, et markup minimal.

## Installation

```sh
go get github.com/glefer/gocharm
```

## Utilisation

```go
import "github.com/glefer/gocharm/core"

console := core.NewConsole()
console.Println("Hello, world!")
```

## Fonctionnalités
- Console avec gestion des couleurs et styles
- Texte enrichi via markup
- API simple et extensible

## Tests

```sh
go test ./...
```

## Licence
MIT
