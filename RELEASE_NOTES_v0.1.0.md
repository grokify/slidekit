# Release Notes - v0.1.0

**Release Date**: 2026-01-26

This is the initial release of slidekit, an AI-optimized toolkit for reading, modifying, and generating slide decks across multiple presentation formats.

## Highlights

- **Canonical slide model** - Unified data types for presentations across formats
- **Marp backend** - Full reader/writer implementation for Marp Markdown
- **TOON format** - Token-optimized output for AI consumption

## Features

### Canonical Data Model

A complete set of Go types representing presentations:

- `Deck` - Complete presentation container with metadata and theme
- `Section` - Slide grouping for LMS chapter support
- `Slide` - Individual slides with 7 layout types
- `Block` - Content blocks (paragraph, bullet, numbered, code, image, quote, heading)
- `Audio` - Audio attachment metadata for TTS and video generation
- `Diff` - Change tracking with add/remove/update/move operations

### Marp Markdown Backend

Full implementation of the Marp Markdown format:

- **Reader**
  - YAML frontmatter parsing with custom fields
  - Slide splitting respecting code fences
  - Marp directive extraction (`<!-- _key: value -->`)
  - Speaker notes with SSML break markers
  - Section detection via `_class: section-divider`
  - Layout inference (title, two-column, section divider)
  - Content type detection (bullets, numbered lists, code, images, quotes)

- **Writer**
  - Generate valid Marp Markdown from canonical model
  - Preserve frontmatter, directives, and speaker notes
  - Support all block types and layouts

- **Backend Interface**
  - `Read` - Parse file to Deck
  - `Write` - Write Deck to file
  - `Plan` - Compute diff between current and desired state
  - `Apply` - Apply diff to existing file
  - `Create` - Create new file from Deck

### TOON Output Format

Token-Optimized Object Notation encoder:

- ~8x more token-efficient than JSON for AI consumption
- Human-readable hierarchical format
- Supports both deck inspection and diff visualization

## API

### Reading a Marp file

```go
reader := marp.NewReader()
deck, err := reader.ReadFile("presentation.md")
```

### Writing a Marp file

```go
writer := marp.NewWriter()
err := writer.WriteFile(deck, "output.md")
```

### TOON output

```go
encoder := format.NewTOONEncoder()
output := encoder.EncodeDeck(deck)
```

### Backend interface

```go
backend := marp.NewBackend()
ctx := context.Background()

deck, err := backend.Read(ctx, ref)
diff, err := backend.Plan(ctx, ref, desired)
err = backend.Apply(ctx, ref, diff)
```

## Testing

- 30 test cases covering reader, writer, and TOON encoder
- Real-world validation with 65-slide production deck
- Round-trip fidelity testing (parse → encode → parse)

## Requirements

- Go 1.23 or later
- No external dependencies (pure stdlib)

## Known Limitations

- CLI and MCP server are stubs (planned for future releases)
- Google Slides and Reveal.js backends not yet implemented
- LMS export functionality not yet implemented

## What's Next

See the [PRD](PRD.md) for the full roadmap:

- **Phase 2**: Google Slides API integration
- **Phase 3**: Reveal.js HTML generation
- **Phase 4**: LMS/Video integration

## Contributors

- John Wang

## License

MIT License
