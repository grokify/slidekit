# slidekit

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

An AI-optimized toolkit for reading, modifying, and generating slide decks across multiple presentation formats.

## Overview

slidekit provides deterministic CLI and MCP (Model Context Protocol) interfaces for AI assistants, backed by reusable Go libraries. It enables programmatic control over presentations with a canonical data model that works across formats.

## Features

- 📦 **Canonical data model** - Unified representation for slides, sections, blocks, and audio metadata
- 🔄 **Multi-format support** - Marp Markdown (implemented), Google Slides, Reveal.js (planned)
- ⚡ **TOON output** - Token-Optimized Object Notation for efficient AI consumption (~8x smaller than JSON)
- 🔁 **Lossless round-tripping** - Parse and regenerate without data loss
- 🎤 **Speaker notes** - Full support for presenter notes with SSML markers
- 🎓 **LMS integration** - Section-based structure for educational platform export (Udemy, Teachable)
- ✅ **Plan/Apply workflow** - Safe, reviewable changes before mutation

## Installation

```bash
go install github.com/grokify/slidekit/cli/cmd/slidekit@latest
```

Or add as a dependency:

```bash
go get github.com/grokify/slidekit
```

## Quick Start

### Parse a Marp Markdown file

```go
package main

import (
    "fmt"
    "github.com/grokify/slidekit/backends/marp"
    "github.com/grokify/slidekit/format"
)

func main() {
    // Parse a Marp Markdown file
    reader := marp.NewReader()
    deck, err := reader.ReadFile("presentation.md")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Title: %s\n", deck.Title)
    fmt.Printf("Sections: %d\n", len(deck.Sections))
    fmt.Printf("Total slides: %d\n", deck.SlideCount())

    // Output in TOON format (token-optimized)
    encoder := format.NewTOONEncoder()
    output := encoder.EncodeDeck(deck)
    fmt.Println(output)
}
```

### Write a deck to Marp Markdown

```go
writer := marp.NewWriter()
err := writer.WriteFile(deck, "output.md")
```

### Use the Backend interface

```go
import "context"

backend := marp.NewBackend()
ctx := context.Background()

// Read
ref := model.Ref{Backend: "marp", Path: "deck.md"}
deck, err := backend.Read(ctx, ref)

// Plan changes
diff, err := backend.Plan(ctx, ref, modifiedDeck)

// Apply changes
err = backend.Apply(ctx, ref, diff)
```

## Data Model

### Core Types

| Type | Description |
|------|-------------|
| `Deck` | Complete presentation with metadata, sections, and theme |
| `Section` | Groups slides (maps to LMS chapters) |
| `Slide` | Individual slide with layout, content, and notes |
| `Block` | Content unit (paragraph, bullet, code, image, quote, heading) |
| `Audio` | Audio attachment for TTS/video generation |
| `Diff` | Change tracking between deck states |

### Slide Layouts

- `title` - Title slide
- `title_body` - Title with body content
- `title_two_col` - Title with two columns
- `section` - Section divider
- `blank` - No predefined structure
- `image` - Full-bleed image
- `comparison` - Side-by-side comparison

### Block Kinds

- `paragraph` - Plain text
- `bullet` - Bulleted list item
- `numbered` - Numbered list item
- `code` - Code block with language
- `image` - Image with URL and alt text
- `quote` - Block quote
- `heading` - Subheading (levels 2-6)

## TOON Output Format

TOON (Token-Optimized Object Notation) provides a compact, human-readable format optimized for AI token efficiency:

```
deck AI & Dev Productivity
meta author John Wang
meta date 2026-01-22

section intro
  slide s1 title
    title AI & Dev Productivity
    subtitle Best Practices for 2026

section fundamentals
  slide s2 title_body
    title Why AI Matters
    bullet Faster iteration
    bullet Lower cognitive load
    note Emphasize the productivity gains
```

## Roadmap

- [x] **Phase 1**: Marp Markdown reader/writer, TOON format, canonical model
- [ ] **Phase 2**: Google Slides integration (read/write/sync)
- [ ] **Phase 3**: Reveal.js HTML generation
- [ ] **Phase 4**: LMS/Video integration (Udemy export, audio assignment)

## Development

### Requirements

- Go 1.23 or later

### Build

```bash
go build ./...
```

### Test

```bash
go test -v ./...
```

### Lint

```bash
golangci-lint run
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome. Please ensure:

1. All tests pass (`go test -v ./...`)
2. Linting passes (`golangci-lint run`)
3. New code includes appropriate tests
4. Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)

## References

- [Marp](https://marp.app/) - Markdown presentation ecosystem
- [Google Slides API](https://developers.google.com/slides/api)
- [Reveal.js](https://revealjs.com/) - HTML presentation framework
- [PRD](PRD.md) - Full product requirements document

 [build-status-svg]: https://github.com/grokify/slidekit/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/slidekit/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/slidekit/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/slidekit/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/slidekit
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/slidekit
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/slidekit
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/slidekit
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fslidekit
 [loc-svg]: https://tokei.rs/b1/github/grokify/slidekit
 [repo-url]: https://github.com/grokify/slidekit
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/slidekit/blob/master/LICENSE
