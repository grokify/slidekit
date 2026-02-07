# Release Notes - v0.2.0

**Release Date**: 2026-02-07

This release adds a full CLI and MCP (Model Context Protocol) server, enabling AI assistants like Claude to read, plan, and modify presentations programmatically.

## Highlights

- **Cobra CLI** - Full command-line interface with read, plan, apply, create, and serve commands
- **MCP Server** - Model Context Protocol server for AI assistant integration
- **Shared ops layer** - Reusable operations package powering both CLI and MCP

## New Features

### CLI Commands

The new `slidekit` CLI provides the following commands:

| Command | Description |
|---------|-------------|
| `slidekit read <file>` | Read presentation, output TOON (default) or JSON |
| `slidekit plan <file> --desired <file>` | Show diff between current and desired state |
| `slidekit apply <file> --diff <file> --confirm` | Apply changes (requires confirmation) |
| `slidekit create <file>` | Create new presentation from stdin JSON |
| `slidekit serve` | Start MCP server on stdio |

Example usage:

```bash
# Read a presentation
slidekit read presentation.md

# Read with JSON output
slidekit read presentation.md --format json

# Plan changes
slidekit plan deck.md --desired updated.json

# Apply with confirmation
slidekit apply deck.md --diff changes.json --confirm
```

### MCP Server

The MCP server enables AI assistants to interact with presentations through the Model Context Protocol.

**Configuration for Claude Code:**

```json
{
  "mcpServers": {
    "slidekit": {
      "command": "/path/to/slidekit",
      "args": ["serve"]
    }
  }
}
```

**Available Tools:**

| Tool | Description |
|------|-------------|
| `read_deck` | Read presentation in TOON/JSON format |
| `list_slides` | List slide IDs and titles |
| `get_slide` | Get single slide by ID |
| `plan_changes` | Compute diff between states |
| `apply_changes` | Apply diff (requires confirm=true) |
| `create_deck` | Create new presentation |
| `update_slide` | Update single slide (requires confirm=true) |

### Shared Operations Layer

The new `ops` package provides shared logic for both CLI and MCP:

- `ReadDeck` - Read with format options (TOON/JSON)
- `PlanChanges` - Compute diffs between states
- `ApplyChanges` - Apply changes with confirmation requirement
- `CreateDeck` - Create new presentations
- `ListSlides`, `GetSlide`, `UpdateSlide` - Slide-level operations

### Safety Features

All mutation operations require explicit confirmation:

- CLI: Use `--confirm` flag
- MCP: Set `confirm=true` in tool input

This prevents accidental modifications and enables a review workflow.

## Dependencies

Added:

- `github.com/spf13/cobra v1.10.0` - CLI framework
- `github.com/modelcontextprotocol/go-sdk v1.2.0` - MCP SDK

## Testing

- 27 new tests in `ops` package with real file I/O
- 11 new tests in `mcp/tools` package
- All tests use temp directories for isolation

## Requirements

- Go 1.23 or later

## Installation

```bash
go install github.com/grokify/slidekit/cli/cmd/slidekit@latest
```

## What's Next

See the [PRD](PRD.md) for the full roadmap:

- **Phase 2**: Google Slides API integration
- **Phase 3**: Reveal.js HTML generation
- **Phase 4**: LMS/Video integration

## Contributors

- John Wang
- Claude Opus 4.5

## License

MIT License
