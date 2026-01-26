package marp

import (
	"testing"

	"github.com/grokify/slidekit/model"
)

const sampleMarp = `---
marp: true
theme: agentplexus
paginate: true
style: |
  .columns {
    display: flex;
    gap: 40px;
  }
---

<!-- _class: lead -->
<!-- _paginate: false -->

<!--
Welcome to the presentation.
[PAUSE:1000]
This is the introduction.
-->

# My Presentation
## A Subtitle Here

**Built with AI**

---

# The Problem

Building AI applications requires:

- **Multiple LLM providers** for redundancy
- **Different APIs** with incompatible formats
- Code duplication for error handling

---

<!-- _class: section-divider -->
<!-- _paginate: false -->

<!--
Section 2: Architecture. <break time="600ms"/>
Let's explore the design.
-->

# Section 2
## Architecture

Understanding the system design

---

<!--
The architecture uses a modular approach.
[PAUSE:1500]
Each component is independent.
-->

# Modular Design

` + "```go" + `
func main() {
    fmt.Println("Hello")
}
` + "```" + `

---

# Two Column Layout

<div class="columns">
<div class="column-left">

**Left Side**
- Item A
- Item B

</div>
<div class="column-right">

**Right Side**
- Item C
- Item D

</div>
</div>

---

# Data Table

| Feature | Status |
|---------|--------|
| Auth | Done |
| API | WIP |

---

# With Image

![Architecture diagram](./images/arch.png)

> This is a blockquote

---

<!-- _class: section-divider -->

# Section 3
## Conclusion

Final thoughts

---

# Thank You

1. Check the repo
2. Star on GitHub
3. Submit PRs
`

func TestParseFrontmatter(t *testing.T) {
	fm, body := parseFrontmatter(sampleMarp)

	if !fm.Marp {
		t.Error("expected marp: true")
	}
	if fm.Theme != "agentplexus" {
		t.Errorf("expected theme 'agentplexus', got %q", fm.Theme)
	}
	if !fm.Paginate {
		t.Error("expected paginate: true")
	}
	if fm.Style == "" {
		t.Error("expected non-empty style")
	}
	if !contains(fm.Style, ".columns") {
		t.Error("expected style to contain .columns")
	}
	if body == "" {
		t.Error("expected non-empty body after frontmatter")
	}
}

func TestSplitSlides(t *testing.T) {
	_, body := parseFrontmatter(sampleMarp)
	slides := splitSlides(body)

	if len(slides) < 8 {
		t.Errorf("expected at least 8 slides, got %d", len(slides))
	}
}

func TestSplitSlidesRespectsCodeBlocks(t *testing.T) {
	content := "# Slide 1\n\n```\nfoo\n---\nbar\n```\n"
	slides := splitSlides(content)

	if len(slides) != 1 {
		t.Errorf("expected 1 slide (--- inside code block), got %d", len(slides))
	}
}

func TestParseDirectives(t *testing.T) {
	raw := "<!-- _class: section-divider -->\n<!-- _paginate: false -->\n\n# Title"
	ps := parseRawSlide(raw)

	if ps.directives["_class"] != "section-divider" {
		t.Errorf("expected _class 'section-divider', got %q", ps.directives["_class"])
	}
	if ps.directives["_paginate"] != "false" {
		t.Errorf("expected _paginate 'false', got %q", ps.directives["_paginate"])
	}
}

func TestParseSpeakerNotes(t *testing.T) {
	raw := `<!--
Welcome to the presentation.
[PAUSE:1000]
This is the introduction.
-->

# Title`

	ps := parseRawSlide(raw)

	if len(ps.notes) == 0 {
		t.Fatal("expected at least one note block")
	}
	if !contains(ps.notes[0].text, "Welcome to the presentation") {
		t.Errorf("expected note to contain 'Welcome', got %q", ps.notes[0].text)
	}
	// PAUSE markers should be stripped from text
	if contains(ps.notes[0].text, "[PAUSE") {
		t.Error("expected PAUSE markers to be stripped from note text")
	}
}

func TestParseSpeakerNotesWithBreak(t *testing.T) {
	raw := `<!--
Section intro. <break time="600ms"/>
More text here.
-->

# Title`

	ps := parseRawSlide(raw)

	if len(ps.notes) == 0 {
		t.Fatal("expected at least one note block")
	}
	if contains(ps.notes[0].text, "<break") {
		t.Error("expected <break> markers to be stripped from note text")
	}
	if !contains(ps.notes[0].text, "Section intro") {
		t.Errorf("expected note to contain 'Section intro', got %q", ps.notes[0].text)
	}
}

func TestParseFullDeck(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Check deck title
	if deck.Title != "My Presentation" {
		t.Errorf("expected title 'My Presentation', got %q", deck.Title)
	}

	// Check theme
	if deck.Theme == nil {
		t.Fatal("expected non-nil theme")
	}
	if deck.Theme.Name != "agentplexus" {
		t.Errorf("expected theme 'agentplexus', got %q", deck.Theme.Name)
	}

	// Check sections
	if len(deck.Sections) < 3 {
		t.Errorf("expected at least 3 sections, got %d", len(deck.Sections))
	}
}

func TestSectionDetection(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// First section should be from the lead slide
	if len(deck.Sections) < 1 {
		t.Fatal("expected at least 1 section")
	}

	// Check that section-divider slides create new sections
	foundArchSection := false
	foundConclusionSection := false
	for _, s := range deck.Sections {
		if contains(s.Title, "Architecture") {
			foundArchSection = true
		}
		if contains(s.Title, "Conclusion") {
			foundConclusionSection = true
		}
	}
	if !foundArchSection {
		t.Error("expected to find Architecture section")
	}
	if !foundConclusionSection {
		t.Error("expected to find Conclusion section")
	}
}

func TestSlideLayouts(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// First slide should be title layout (from lead class)
	firstSlide := deck.Sections[0].Slides[0]
	if firstSlide.Layout != model.LayoutTitle {
		t.Errorf("expected first slide layout 'title', got %q", firstSlide.Layout)
	}

	// Find a section-divider slide
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			if slide.Layout == model.LayoutSection {
				return // Found one, test passes
			}
		}
	}
	t.Error("expected to find at least one section layout slide")
}

func TestCodeBlockParsing(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Find the slide with a code block
	found := false
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			for _, block := range slide.Body {
				if block.Kind == model.BlockCode {
					found = true
					if block.Lang != "go" {
						t.Errorf("expected code lang 'go', got %q", block.Lang)
					}
					if !contains(block.Text, "fmt.Println") {
						t.Error("expected code to contain fmt.Println")
					}
				}
			}
		}
	}
	if !found {
		t.Error("expected to find a code block")
	}
}

func TestBulletParsing(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Find "The Problem" slide
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			if slide.Title == "The Problem" {
				bulletCount := 0
				for _, block := range slide.Body {
					if block.Kind == model.BlockBullet {
						bulletCount++
					}
				}
				if bulletCount < 3 {
					t.Errorf("expected at least 3 bullets in 'The Problem', got %d", bulletCount)
				}
				return
			}
		}
	}
	t.Error("expected to find 'The Problem' slide")
}

func TestNumberedListParsing(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Find "Thank You" slide with numbered list
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			if slide.Title == "Thank You" {
				numCount := 0
				for _, block := range slide.Body {
					if block.Kind == model.BlockNumbered {
						numCount++
					}
				}
				if numCount < 3 {
					t.Errorf("expected at least 3 numbered items in 'Thank You', got %d", numCount)
				}
				return
			}
		}
	}
	t.Error("expected to find 'Thank You' slide")
}

func TestImageParsing(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Find image block
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			for _, block := range slide.Body {
				if block.Kind == model.BlockImage {
					if block.URL != "./images/arch.png" {
						t.Errorf("expected image URL './images/arch.png', got %q", block.URL)
					}
					if block.Alt != "Architecture diagram" {
						t.Errorf("expected alt 'Architecture diagram', got %q", block.Alt)
					}
					return
				}
			}
		}
	}
	t.Error("expected to find an image block")
}

func TestQuoteParsing(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			for _, block := range slide.Body {
				if block.Kind == model.BlockQuote {
					if !contains(block.Text, "blockquote") {
						t.Errorf("expected quote text to contain 'blockquote', got %q", block.Text)
					}
					return
				}
			}
		}
	}
	t.Error("expected to find a quote block")
}

func TestTwoColumnDetection(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Find "Two Column Layout" slide
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			if slide.Title == "Two Column Layout" {
				if slide.Layout != model.LayoutTitleTwoCol {
					t.Errorf("expected layout 'title_two_col', got %q", slide.Layout)
				}
				return
			}
		}
	}
	t.Error("expected to find 'Two Column Layout' slide")
}

func TestSpeakerNotesOnSlides(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Count slides with notes
	notesCount := 0
	for _, section := range deck.Sections {
		for _, slide := range section.Slides {
			if len(slide.Notes) > 0 {
				notesCount++
			}
		}
	}
	if notesCount < 3 {
		t.Errorf("expected at least 3 slides with notes, got %d", notesCount)
	}
}

func TestSlideCount(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	count := deck.SlideCount()
	if count < 8 {
		t.Errorf("expected at least 8 total slides, got %d", count)
	}
}

func TestCleanNoteText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello world.\n[PAUSE:1000]\nNext sentence.",
			expected: "Hello world. Next sentence.",
		},
		{
			input:    "Intro. <break time=\"600ms\"/>\nMore text.",
			expected: "Intro. More text.",
		},
		{
			input:    "Simple note.",
			expected: "Simple note.",
		},
	}

	for _, tt := range tests {
		result := cleanNoteText(tt.input)
		if result != tt.expected {
			t.Errorf("cleanNoteText(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestEmptyInput(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse("")
	if err != nil {
		t.Fatalf("Parse error on empty input: %v", err)
	}
	if deck.SlideCount() != 0 {
		t.Errorf("expected 0 slides for empty input, got %d", deck.SlideCount())
	}
}

func TestFrontmatterOnly(t *testing.T) {
	input := `---
marp: true
theme: default
---`
	reader := NewReader()
	deck, err := reader.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if deck.Theme == nil || deck.Theme.Name != "default" {
		t.Error("expected theme 'default'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
