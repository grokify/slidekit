package marp

import (
	"strings"
	"testing"

	"github.com/grokify/slidekit/model"
)

func TestWriteSimpleDeck(t *testing.T) {
	deck := &model.Deck{
		Title: "Test Deck",
		Theme: &model.Theme{Name: "default"},
		Sections: []model.Section{
			{
				ID:    "intro",
				Title: "Introduction",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitle,
						Title:  "Test Deck",
						Notes: []model.Block{
							model.NewParagraph("Welcome to the deck."),
						},
					},
					{
						ID:     "s2",
						Layout: model.LayoutTitleBody,
						Title:  "Key Points",
						Body: []model.Block{
							model.NewBullet("First point", 0),
							model.NewBullet("Second point", 0),
							model.NewBullet("Sub point", 1),
						},
					},
				},
			},
		},
	}

	writer := NewWriter()
	output := writer.Encode(deck)

	// Check frontmatter
	if !strings.Contains(output, "marp: true") {
		t.Error("expected 'marp: true' in output")
	}
	if !strings.Contains(output, "theme: default") {
		t.Error("expected 'theme: default' in output")
	}

	// Check slide separators
	if !strings.Contains(output, "---") {
		t.Error("expected slide separator in output")
	}

	// Check title
	if !strings.Contains(output, "# Test Deck") {
		t.Error("expected '# Test Deck' in output")
	}

	// Check bullets
	if !strings.Contains(output, "- First point") {
		t.Error("expected '- First point' in output")
	}
	if !strings.Contains(output, "    - Sub point") {
		t.Error("expected indented '    - Sub point' in output")
	}

	// Check speaker notes
	if !strings.Contains(output, "<!--\nWelcome to the deck.\n-->") {
		t.Error("expected speaker notes in output")
	}

	// Check lead class for title slide
	if !strings.Contains(output, "<!-- _class: lead -->") {
		t.Error("expected lead class directive for title slide")
	}
}

func TestWriteCodeBlock(t *testing.T) {
	deck := &model.Deck{
		Title: "Code Demo",
		Theme: &model.Theme{Name: "default"},
		Sections: []model.Section{
			{
				ID:    "s1",
				Title: "Code",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitleBody,
						Title:  "Code Example",
						Body: []model.Block{
							model.NewCode("fmt.Println(\"hello\")", "go"),
						},
					},
				},
			},
		},
	}

	writer := NewWriter()
	output := writer.Encode(deck)

	if !strings.Contains(output, "```go") {
		t.Error("expected code block with 'go' language")
	}
	if !strings.Contains(output, "fmt.Println") {
		t.Error("expected code content")
	}
}

func TestWriteSectionDivider(t *testing.T) {
	deck := &model.Deck{
		Title: "Sections",
		Theme: &model.Theme{Name: "default"},
		Sections: []model.Section{
			{
				ID:    "s1",
				Title: "Section 1",
				Slides: []model.Slide{
					{
						ID:     "div1",
						Layout: model.LayoutSection,
						Title:  "Section 1",
					},
				},
			},
		},
	}

	writer := NewWriter()
	output := writer.Encode(deck)

	if !strings.Contains(output, "<!-- _class: section-divider -->") {
		t.Error("expected section-divider class in output")
	}
}

func TestWriteImage(t *testing.T) {
	deck := &model.Deck{
		Title: "Images",
		Theme: &model.Theme{Name: "default"},
		Sections: []model.Section{
			{
				ID:    "s1",
				Title: "Images",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitleBody,
						Title:  "Architecture",
						Body: []model.Block{
							model.NewImage("./arch.png", "Architecture diagram"),
						},
					},
				},
			},
		},
	}

	writer := NewWriter()
	output := writer.Encode(deck)

	if !strings.Contains(output, "![Architecture diagram](./arch.png)") {
		t.Error("expected image markdown in output")
	}
}

func TestWriteQuote(t *testing.T) {
	deck := &model.Deck{
		Title: "Quotes",
		Theme: &model.Theme{Name: "default"},
		Sections: []model.Section{
			{
				ID:    "s1",
				Title: "Quotes",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitleBody,
						Title:  "Inspiration",
						Body: []model.Block{
							model.NewQuote("To be or not to be"),
						},
					},
				},
			},
		},
	}

	writer := NewWriter()
	output := writer.Encode(deck)

	if !strings.Contains(output, "> To be or not to be") {
		t.Error("expected blockquote in output")
	}
}

func TestRoundTrip(t *testing.T) {
	// Parse the sample, write it out, parse again, compare structure
	reader := NewReader()
	writer := NewWriter()

	deck1, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("first parse error: %v", err)
	}

	encoded := writer.Encode(deck1)

	deck2, err := reader.Parse(encoded)
	if err != nil {
		t.Fatalf("second parse error: %v", err)
	}

	// Compare structural properties
	if deck1.Title != deck2.Title {
		t.Errorf("title mismatch: %q vs %q", deck1.Title, deck2.Title)
	}
	if len(deck1.Sections) != len(deck2.Sections) {
		t.Errorf("section count mismatch: %d vs %d", len(deck1.Sections), len(deck2.Sections))
	}
	if deck1.SlideCount() != deck2.SlideCount() {
		t.Errorf("slide count mismatch: %d vs %d", deck1.SlideCount(), deck2.SlideCount())
	}
}

func TestWriteNumberedList(t *testing.T) {
	deck := &model.Deck{
		Title: "Lists",
		Theme: &model.Theme{Name: "default"},
		Sections: []model.Section{
			{
				ID:    "s1",
				Title: "Lists",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitleBody,
						Title:  "Steps",
						Body: []model.Block{
							model.NewNumbered("First step", 0),
							model.NewNumbered("Second step", 0),
						},
					},
				},
			},
		},
	}

	writer := NewWriter()
	output := writer.Encode(deck)

	if !strings.Contains(output, "1. First step") {
		t.Error("expected '1. First step' in output")
	}
	if !strings.Contains(output, "1. Second step") {
		t.Error("expected '1. Second step' in output")
	}
}
