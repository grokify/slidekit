# slidekit PRD

**Version**: 0.1.0-draft
**Last Updated**: 2026-01-22

## Overview

slidekit is an AI-optimized toolkit for reading, modifying, and generating slide decks across formats. It provides deterministic CLI and MCP interfaces for AI assistants, backed by generic, reusable Go libraries.

## Problem Statement

Creating and managing presentations across multiple platforms (Marp, Google Slides, Reveal.js, etc.) requires manual conversion and lacks programmatic control. AI assistants need structured, token-efficient interfaces to work with presentations safely and deterministically.

Additionally, educational content creators need to generate multi-section video courses from slide decks with per-section audio for LMS platforms like Udemy.

## Goals

### Phase 1: Foundation

- Define canonical slide model supporting sections, speaker notes, and audio metadata
- Implement Marp Markdown reader/writer (current workflow)
- CLI with TOON output format (token-optimized)
- Basic MCP server with inspection tools

### Phase 2: Google Slides Integration

- Read existing Google Slides presentations
- Create new presentations from canonical model
- Update existing presentations (plan/apply workflow)
- Bidirectional sync: Google Slides ↔ Marp Markdown

### Phase 3: Reveal.js Integration

- Generate Reveal.js HTML with dynamic transitions
- Support Reveal.js-specific features (fragments, vertical slides)
- Live preview server

### Phase 4: LMS/Video Integration

- Section-based audio assignment
- Video segment generation per section
- LMS metadata export (Udemy, Teachable format)
- Course outline generation

### Future Phases

- Gamma integration (API pending)
- Canva integration (API pending)
- PPTX export/import
- PDF generation improvements

## Architecture

```
slidekit/
  model/              ← Canonical slide model
    deck.go
    slide.go
    section.go
    audio.go

  backends/
    marp/             ← Marp Markdown reader/writer
    gslides/          ← Google Slides API
    reveal/           ← Reveal.js HTML generator

  ops/                ← Semantic operations
    diff.go
    patch.go
    validate.go

  cli/
    cmd/slidekit/
    format/           ← TOON/JSON serialization

  mcp/
    server/
    tools/

  lms/                ← LMS/Video integration
    course.go
    video.go
    udemy.go

  internal/
    auth/
```

## Data Model

### Core Types

```go
// Deck represents a complete presentation
type Deck struct {
    ID       string            `json:"id"`
    Title    string            `json:"title"`
    Meta     Meta              `json:"meta"`
    Sections []Section         `json:"sections"`
    Theme    *Theme            `json:"theme,omitempty"`
}

// Meta contains presentation metadata
type Meta struct {
    Author      string            `json:"author,omitempty"`
    Date        string            `json:"date,omitempty"`
    Description string            `json:"description,omitempty"`
    Keywords    []string          `json:"keywords,omitempty"`
    Custom      map[string]string `json:"custom,omitempty"`
}

// Section groups related slides (maps to LMS sections/chapters)
type Section struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Slides []Slide `json:"slides"`
    Audio  *Audio  `json:"audio,omitempty"`  // Section-level audio for LMS
}

// Slide represents a single slide
type Slide struct {
    ID         string   `json:"id"`
    Layout     Layout   `json:"layout"`
    Title      string   `json:"title,omitempty"`
    Subtitle   string   `json:"subtitle,omitempty"`
    Body       []Block  `json:"body,omitempty"`
    Notes      []Block  `json:"notes,omitempty"`      // Speaker notes
    Audio      *Audio   `json:"audio,omitempty"`      // Slide-level audio
    Transition *string  `json:"transition,omitempty"` // Reveal.js transitions
    Background *string  `json:"background,omitempty"`
}

// Block represents content within a slide
type Block struct {
    Kind  BlockKind `json:"kind"`
    Text  string    `json:"text,omitempty"`
    Level int       `json:"level,omitempty"` // Bullet nesting level
    Lang  string    `json:"lang,omitempty"`  // Code language
    URL   string    `json:"url,omitempty"`   // Image/link URL
    Alt   string    `json:"alt,omitempty"`   // Image alt text
}

// BlockKind identifies content type
type BlockKind string

const (
    BlockParagraph BlockKind = "paragraph"
    BlockBullet    BlockKind = "bullet"
    BlockNumbered  BlockKind = "numbered"
    BlockCode      BlockKind = "code"
    BlockImage     BlockKind = "image"
    BlockQuote     BlockKind = "quote"
    BlockHeading   BlockKind = "heading"
)

// Layout identifies slide layout type
type Layout string

const (
    LayoutTitle        Layout = "title"          // Title slide
    LayoutTitleBody    Layout = "title_body"     // Title + body content
    LayoutTitleTwoCol  Layout = "title_two_col"  // Title + two columns
    LayoutSection      Layout = "section"        // Section divider
    LayoutBlank        Layout = "blank"          // No predefined structure
    LayoutImage        Layout = "image"          // Full-bleed image
    LayoutComparison   Layout = "comparison"     // Side-by-side comparison
)

// Audio represents audio attachment for LMS/video generation
type Audio struct {
    Source   AudioSource `json:"source"`
    Path     string      `json:"path,omitempty"`     // Local file path
    URL      string      `json:"url,omitempty"`      // Remote URL
    Script   string      `json:"script,omitempty"`   // Text for TTS generation
    Duration Duration    `json:"duration,omitempty"` // Explicit duration
    Voice    string      `json:"voice,omitempty"`    // TTS voice ID
}

// AudioSource identifies audio origin
type AudioSource string

const (
    AudioSourceFile   AudioSource = "file"   // Pre-recorded audio file
    AudioSourceURL    AudioSource = "url"    // Remote audio URL
    AudioSourceTTS    AudioSource = "tts"    // Generate via TTS from script
    AudioSourceNotes  AudioSource = "notes"  // Generate via TTS from speaker notes
)

// Theme represents presentation styling
type Theme struct {
    Name       string            `json:"name,omitempty"`
    Primary    string            `json:"primary,omitempty"`
    Secondary  string            `json:"secondary,omitempty"`
    Background string            `json:"background,omitempty"`
    Font       string            `json:"font,omitempty"`
    Custom     map[string]string `json:"custom,omitempty"`
}
```

### Diff Types

```go
// Diff represents changes between two deck states
type Diff struct {
    DeckID  string       `json:"deck_id"`
    Changes []Change     `json:"changes"`
}

// Change represents a single modification
type Change struct {
    Op        ChangeOp `json:"op"`
    Path      string   `json:"path"`      // e.g., "sections/0/slides/2/title"
    SlideID   string   `json:"slide_id,omitempty"`
    SectionID string   `json:"section_id,omitempty"`
    OldValue  any      `json:"old_value,omitempty"`
    NewValue  any      `json:"new_value,omitempty"`
}

// ChangeOp identifies the type of change
type ChangeOp string

const (
    ChangeAdd    ChangeOp = "add"
    ChangeRemove ChangeOp = "remove"
    ChangeUpdate ChangeOp = "update"
    ChangeMove   ChangeOp = "move"
)
```

## Backend Interface

```go
// Backend defines the interface for presentation backends
type Backend interface {
    // Info returns backend metadata
    Info() BackendInfo

    // Read loads a presentation from the backend
    Read(ctx context.Context, ref Ref) (*Deck, error)

    // Plan computes changes needed to reach desired state
    Plan(ctx context.Context, ref Ref, desired *Deck) (*Diff, error)

    // Apply executes a diff against the backend
    Apply(ctx context.Context, ref Ref, diff *Diff) error

    // Create creates a new presentation
    Create(ctx context.Context, deck *Deck) (Ref, error)
}

// Ref identifies a presentation in a backend
type Ref struct {
    Backend string `json:"backend"` // "marp", "gslides", "reveal"
    ID      string `json:"id"`      // Backend-specific identifier
    Path    string `json:"path"`    // File path (for file-based backends)
}

// BackendInfo describes backend capabilities
type BackendInfo struct {
    Name         string   `json:"name"`
    Version      string   `json:"version"`
    Capabilities []string `json:"capabilities"`
    // Capabilities: "read", "write", "plan", "apply", "create",
    //               "transitions", "audio", "sections"
}
```

## CLI Specification

### Command Structure

```
slidekit <backend> <command> [options]
slidekit marp inspect --file deck.md
slidekit gslides inspect --deck <id>
slidekit gslides plan --deck <id> --from deck.md
slidekit gslides apply --deck <id> --from deck.md --yes
slidekit reveal generate --from deck.md --out ./dist
slidekit course export --from deck.md --format udemy --out course/
```

### Global Options

| Option | Description |
|--------|-------------|
| `--format toon` | Output format: toon (default), json |
| `--quiet` | Suppress non-essential output |
| `--verbose` | Verbose output |
| `--config` | Config file path |

### Backend-Specific Commands

#### `slidekit marp`

| Command | Description |
|---------|-------------|
| `inspect --file <path>` | Parse and display deck structure |
| `validate --file <path>` | Validate Marp Markdown |
| `convert --file <path> --to <format>` | Convert to HTML, PDF (wraps marp CLI) |

#### `slidekit gslides`

| Command | Description |
|---------|-------------|
| `auth` | Authenticate with Google |
| `list` | List accessible presentations |
| `inspect --deck <id>` | Display deck structure |
| `plan --deck <id> --from <file>` | Show planned changes |
| `apply --deck <id> --from <file> [--yes]` | Apply changes |
| `create --from <file> [--title <title>]` | Create new presentation |
| `export --deck <id> --to <file>` | Export to Marp Markdown |

#### `slidekit reveal`

| Command | Description |
|---------|-------------|
| `generate --from <file> --out <dir>` | Generate Reveal.js HTML |
| `serve --from <file> [--port 8080]` | Live preview server |
| `export --from <file> --format pdf` | Export to PDF |

#### `slidekit course`

| Command | Description |
|---------|-------------|
| `inspect --from <file>` | Display course structure |
| `export --from <file> --format <fmt> --out <dir>` | Export for LMS |
| `audio assign --from <file> --audio-dir <dir>` | Assign audio files to sections |
| `video generate --from <file> --out <dir>` | Generate video segments |

### TOON Output Format

Default output format optimized for AI token efficiency.

**Deck inspection:**
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
    bullet Better documentation
    note Emphasize the productivity gains
    audio tts voice=en-US-Neural2-J

  slide s3 title_body
    title Key Concepts
    bullet Foundation models
    bullet Prompt engineering
    bullet Tool use
```

**Plan output:**
```
plan deck abc123
~ section/intro/slide/s1/title
  - AI & Dev Productivity
  + AI & Developer Productivity
+ section/advanced
  + slide s10 title_body
    title Advanced Topics
```

## MCP Tool Specification

### Tool: `slidekit.inspect`

Inspect a presentation from any backend.

```json
{
  "name": "slidekit.inspect",
  "description": "Inspect presentation structure",
  "inputSchema": {
    "type": "object",
    "properties": {
      "backend": { "type": "string", "enum": ["marp", "gslides", "reveal"] },
      "ref": { "type": "string", "description": "File path or deck ID" },
      "format": { "type": "string", "enum": ["toon", "json"], "default": "toon" }
    },
    "required": ["backend", "ref"]
  }
}
```

### Tool: `slidekit.gslides.list`

List accessible Google Slides presentations.

```json
{
  "name": "slidekit.gslides.list",
  "description": "List Google Slides presentations",
  "inputSchema": {
    "type": "object",
    "properties": {
      "limit": { "type": "integer", "default": 20 },
      "format": { "type": "string", "enum": ["toon", "json"], "default": "toon" }
    }
  }
}
```

### Tool: `slidekit.gslides.plan`

Plan changes to a Google Slides presentation.

```json
{
  "name": "slidekit.gslides.plan",
  "description": "Plan changes to Google Slides presentation",
  "inputSchema": {
    "type": "object",
    "properties": {
      "deck_id": { "type": "string" },
      "source": { "type": "string", "description": "Marp Markdown content or file path" },
      "format": { "type": "string", "enum": ["toon", "json"], "default": "toon" }
    },
    "required": ["deck_id", "source"]
  }
}
```

### Tool: `slidekit.gslides.apply`

Apply planned changes to a Google Slides presentation.

```json
{
  "name": "slidekit.gslides.apply",
  "description": "Apply changes to Google Slides presentation",
  "inputSchema": {
    "type": "object",
    "properties": {
      "deck_id": { "type": "string" },
      "diff": { "type": "string", "description": "Diff in TOON or JSON format" },
      "confirm": { "type": "boolean", "default": false }
    },
    "required": ["deck_id", "diff", "confirm"]
  }
}
```

### Tool: `slidekit.gslides.create`

Create a new Google Slides presentation.

```json
{
  "name": "slidekit.gslides.create",
  "description": "Create new Google Slides presentation",
  "inputSchema": {
    "type": "object",
    "properties": {
      "source": { "type": "string", "description": "Marp Markdown content or file path" },
      "title": { "type": "string" }
    },
    "required": ["source"]
  }
}
```

### Tool: `slidekit.course.inspect`

Inspect course structure from a deck.

```json
{
  "name": "slidekit.course.inspect",
  "description": "Inspect course structure for LMS export",
  "inputSchema": {
    "type": "object",
    "properties": {
      "source": { "type": "string", "description": "Marp Markdown file path" },
      "format": { "type": "string", "enum": ["toon", "json"], "default": "toon" }
    },
    "required": ["source"]
  }
}
```

## LMS Integration

### Course Structure

A course maps directly to a Deck with Sections:

| Deck Component | LMS Equivalent |
|----------------|----------------|
| Deck | Course |
| Section | Chapter/Module |
| Slides in Section | Lecture content |
| Section Audio | Lecture video audio track |

### Udemy Export Format

```
course/
  course.json           # Course metadata
  sections/
    01-intro/
      slides.pdf        # Section slides as PDF
      audio.mp3         # Section audio
      video.mp4         # Combined video (slides + audio)
      transcript.txt    # Audio transcript
    02-fundamentals/
      ...
```

### Audio Assignment

Audio can be assigned at section or slide level:

1. **Section-level audio** (recommended for Udemy): One audio track per section
2. **Slide-level audio**: Individual audio per slide, concatenated for video

Audio sources:

- **File**: Pre-recorded audio file
- **TTS from script**: Generate from explicit script text
- **TTS from notes**: Generate from speaker notes

## Security Considerations

- OAuth tokens stored securely (system keychain preferred)
- Service account credentials for MCP server
- No free-form code execution
- Deck IDs validated before mutation
- Automatic backup before destructive operations

## Success Metrics

- Successfully round-trip Marp ↔ Google Slides without data loss
- Generate Reveal.js presentations matching Marp visual intent
- Export course to Udemy format with synchronized audio/video
- CLI commands complete in < 5s for typical decks
- MCP tools work reliably with Claude Code and other AI assistants

## Open Questions

1. **ID stability**: How to maintain stable slide IDs across backends?
2. **Theme mapping**: How to map Marp themes to Google Slides templates?
3. **Transition mapping**: How to handle Reveal.js transitions in other backends?
4. **Audio timing**: How to handle slide timing when audio duration varies?
5. **Gamma/Canva APIs**: Are these APIs publicly available?

## References

- [Google Slides API](https://developers.google.com/slides/api)
- [Marp](https://marp.app/)
- [Reveal.js](https://revealjs.com/)
- [Udemy Course Creation](https://www.udemy.com/instructor/resources/)
- [TOON Format](https://github.com/grokify/toon) (token-optimized object notation)
