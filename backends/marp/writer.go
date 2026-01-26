package marp

import (
	"fmt"
	"os"
	"strings"

	"github.com/grokify/slidekit/model"
)

// Writer converts a Deck back to Marp Markdown format.
type Writer struct{}

// NewWriter creates a new Marp writer.
func NewWriter() *Writer {
	return &Writer{}
}

// WriteFile writes a deck to a Marp Markdown file.
func (w *Writer) WriteFile(deck *model.Deck, path string) error {
	content := w.Encode(deck)
	return os.WriteFile(path, []byte(content), 0644)
}

// Encode converts a Deck to Marp Markdown string.
func (w *Writer) Encode(deck *model.Deck) string {
	var b strings.Builder

	// Write frontmatter
	w.writeFrontmatter(&b, deck)

	// Write slides
	first := true
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			if !first {
				b.WriteString("\n---\n\n")
			}
			first = false
			w.writeSlide(&b, &slide)
		}
	}

	return b.String()
}

func (w *Writer) writeFrontmatter(b *strings.Builder, deck *model.Deck) {
	b.WriteString("---\n")
	b.WriteString("marp: true\n")

	if deck.Theme != nil && deck.Theme.Name != "" {
		fmt.Fprintf(b, "theme: %s\n", deck.Theme.Name)
	}

	b.WriteString("paginate: true\n")

	// Write custom style if present
	if deck.Theme != nil {
		style := deck.Theme.GetCustom("style", "")
		if style != "" {
			b.WriteString("style: |\n")
			for _, line := range strings.Split(style, "\n") {
				b.WriteString(line)
				b.WriteString("\n")
			}
		}
	}

	b.WriteString("---\n\n")
}

func (w *Writer) writeSlide(b *strings.Builder, slide *model.Slide) {
	// Write directives
	switch slide.Layout {
	case model.LayoutSection:
		b.WriteString("<!-- _class: section-divider -->\n")
		b.WriteString("<!-- _paginate: false -->\n\n")
	case model.LayoutTitle:
		b.WriteString("<!-- _class: lead -->\n")
		b.WriteString("<!-- _paginate: false -->\n\n")
	}

	// Write speaker notes before content
	if len(slide.Notes) > 0 {
		b.WriteString("<!--\n")
		for _, note := range slide.Notes {
			b.WriteString(note.Text)
			b.WriteString("\n")
		}
		b.WriteString("-->\n\n")
	}

	// Write title
	if slide.Title != "" {
		fmt.Fprintf(b, "# %s\n", slide.Title)
	}

	// Write subtitle
	if slide.Subtitle != "" {
		fmt.Fprintf(b, "## %s\n", slide.Subtitle)
	}

	// Write body blocks
	if slide.Title != "" || slide.Subtitle != "" {
		b.WriteString("\n")
	}
	for _, block := range slide.Body {
		w.writeBlock(b, &block)
	}
}

func (w *Writer) writeBlock(b *strings.Builder, block *model.Block) {
	switch block.Kind {
	case model.BlockBullet:
		indent := strings.Repeat("    ", block.Level)
		fmt.Fprintf(b, "%s- %s\n", indent, block.Text)
	case model.BlockNumbered:
		indent := strings.Repeat("    ", block.Level)
		fmt.Fprintf(b, "%s1. %s\n", indent, block.Text)
	case model.BlockParagraph:
		if strings.HasPrefix(strings.TrimSpace(block.Text), "<") {
			// HTML block - write as-is
			b.WriteString(block.Text)
			b.WriteString("\n")
		} else if strings.HasPrefix(strings.TrimSpace(block.Text), "|") {
			// Table - write as-is
			b.WriteString(block.Text)
			b.WriteString("\n")
		} else {
			b.WriteString("\n")
			b.WriteString(block.Text)
			b.WriteString("\n")
		}
	case model.BlockCode:
		fmt.Fprintf(b, "\n```%s\n%s\n```\n", block.Lang, block.Text)
	case model.BlockImage:
		fmt.Fprintf(b, "![%s](%s)\n", block.Alt, block.URL)
	case model.BlockQuote:
		fmt.Fprintf(b, "> %s\n", block.Text)
	case model.BlockHeading:
		prefix := strings.Repeat("#", block.Level+1) // Level 1 = ##, Level 2 = ###
		fmt.Fprintf(b, "%s %s\n", prefix, block.Text)
	}
}
