// Package format provides serialization formats for slidekit.
// TOON (Token-Optimized Object Notation) is the default format for AI agents.
// JSON is available for human readability.
package format

import (
	"fmt"
	"strings"

	"github.com/grokify/slidekit/model"
)

// TOONEncoder encodes decks and diffs to TOON format.
type TOONEncoder struct {
	indent string
}

// NewTOONEncoder creates a new TOON encoder.
func NewTOONEncoder() *TOONEncoder {
	return &TOONEncoder{indent: "  "}
}

// EncodeDeck encodes a deck to TOON format.
func (e *TOONEncoder) EncodeDeck(d *model.Deck) string {
	var b strings.Builder

	// Header
	b.WriteString("deck ")
	b.WriteString(d.Title)
	b.WriteString("\n")

	// Metadata
	if d.Meta.Author != "" {
		b.WriteString("meta author ")
		b.WriteString(d.Meta.Author)
		b.WriteString("\n")
	}
	if d.Meta.Date != "" {
		b.WriteString("meta date ")
		b.WriteString(d.Meta.Date)
		b.WriteString("\n")
	}
	if d.Meta.Description != "" {
		b.WriteString("meta description ")
		b.WriteString(d.Meta.Description)
		b.WriteString("\n")
	}

	// Sections
	for _, section := range d.Sections {
		b.WriteString("\n")
		e.encodeSection(&b, &section)
	}

	return b.String()
}

func (e *TOONEncoder) encodeSection(b *strings.Builder, s *model.Section) {
	b.WriteString("section ")
	b.WriteString(s.ID)
	if s.Title != "" {
		b.WriteString(" ")
		b.WriteString(s.Title)
	}
	b.WriteString("\n")

	// Section audio
	if s.Audio != nil {
		b.WriteString(e.indent)
		e.encodeAudio(b, s.Audio)
		b.WriteString("\n")
	}

	// Slides
	for _, slide := range s.Slides {
		e.encodeSlide(b, &slide)
	}
}

func (e *TOONEncoder) encodeSlide(b *strings.Builder, s *model.Slide) {
	b.WriteString(e.indent)
	b.WriteString("slide ")
	b.WriteString(s.ID)
	b.WriteString(" ")
	b.WriteString(string(s.Layout))
	b.WriteString("\n")

	indent2 := e.indent + e.indent

	// Title
	if s.Title != "" {
		b.WriteString(indent2)
		b.WriteString("title ")
		b.WriteString(s.Title)
		b.WriteString("\n")
	}

	// Subtitle
	if s.Subtitle != "" {
		b.WriteString(indent2)
		b.WriteString("subtitle ")
		b.WriteString(s.Subtitle)
		b.WriteString("\n")
	}

	// Body
	for _, block := range s.Body {
		b.WriteString(indent2)
		e.encodeBlock(b, &block)
		b.WriteString("\n")
	}

	// Notes
	for _, block := range s.Notes {
		b.WriteString(indent2)
		b.WriteString("note ")
		b.WriteString(block.Text)
		b.WriteString("\n")
	}

	// Audio
	if s.Audio != nil {
		b.WriteString(indent2)
		e.encodeAudio(b, s.Audio)
		b.WriteString("\n")
	}

	// Transition
	if s.Transition != nil {
		b.WriteString(indent2)
		b.WriteString("transition ")
		b.WriteString(*s.Transition)
		b.WriteString("\n")
	}
}

func (e *TOONEncoder) encodeBlock(b *strings.Builder, block *model.Block) {
	switch block.Kind {
	case model.BlockBullet:
		for i := 0; i < block.Level; i++ {
			b.WriteString("  ")
		}
		b.WriteString("bullet ")
		b.WriteString(block.Text)
	case model.BlockNumbered:
		for i := 0; i < block.Level; i++ {
			b.WriteString("  ")
		}
		b.WriteString("numbered ")
		b.WriteString(block.Text)
	case model.BlockParagraph:
		b.WriteString("para ")
		b.WriteString(block.Text)
	case model.BlockCode:
		b.WriteString("code ")
		if block.Lang != "" {
			b.WriteString(block.Lang)
			b.WriteString(" ")
		}
		b.WriteString(block.Text)
	case model.BlockImage:
		b.WriteString("image ")
		b.WriteString(block.URL)
		if block.Alt != "" {
			b.WriteString(" ")
			b.WriteString(block.Alt)
		}
	case model.BlockQuote:
		b.WriteString("quote ")
		b.WriteString(block.Text)
	case model.BlockHeading:
		b.WriteString("heading ")
		if block.Level > 0 {
			fmt.Fprintf(b, "%d ", block.Level)
		}
		b.WriteString(block.Text)
	}
}

func (e *TOONEncoder) encodeAudio(b *strings.Builder, a *model.Audio) {
	b.WriteString("audio ")
	b.WriteString(string(a.Source))
	switch a.Source {
	case model.AudioSourceFile:
		b.WriteString(" ")
		b.WriteString(a.Path)
	case model.AudioSourceURL:
		b.WriteString(" ")
		b.WriteString(a.URL)
	case model.AudioSourceTTS:
		if a.Voice != "" {
			b.WriteString(" voice=")
			b.WriteString(a.Voice)
		}
	case model.AudioSourceNotes:
		if a.Voice != "" {
			b.WriteString(" voice=")
			b.WriteString(a.Voice)
		}
	}
	if a.Duration > 0 {
		fmt.Fprintf(b, " duration=%s", a.Duration)
	}
}

// EncodeDiff encodes a diff to TOON format.
func (e *TOONEncoder) EncodeDiff(d *model.Diff) string {
	var b strings.Builder

	b.WriteString("plan deck ")
	b.WriteString(d.DeckID)
	b.WriteString("\n")

	for _, c := range d.Changes {
		e.encodeChange(&b, &c)
	}

	return b.String()
}

func (e *TOONEncoder) encodeChange(b *strings.Builder, c *model.Change) {
	switch c.Op {
	case model.ChangeAdd:
		b.WriteString("+ ")
	case model.ChangeRemove:
		b.WriteString("- ")
	case model.ChangeUpdate:
		b.WriteString("~ ")
	case model.ChangeMove:
		b.WriteString("> ")
	}

	b.WriteString(c.Path)
	b.WriteString("\n")

	if c.OldValue != nil {
		b.WriteString(e.indent)
		b.WriteString("- ")
		fmt.Fprintf(b, "%v", c.OldValue)
		b.WriteString("\n")
	}
	if c.NewValue != nil {
		b.WriteString(e.indent)
		b.WriteString("+ ")
		fmt.Fprintf(b, "%v", c.NewValue)
		b.WriteString("\n")
	}
}

// Format is an output format type.
type Format string

const (
	FormatTOON Format = "toon"
	FormatJSON Format = "json"
)

// IsValid returns true if the format is recognized.
func (f Format) IsValid() bool {
	return f == FormatTOON || f == FormatJSON
}
