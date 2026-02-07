package format

import (
	"strings"
	"testing"
	"time"

	"github.com/grokify/slidekit/model"
)

func TestTOONEncoderEncodeDeck(t *testing.T) {
	deck := &model.Deck{
		Title: "Test Presentation",
		Meta: model.Meta{
			Author: "John Doe",
			Date:   "2026-01-26",
		},
		Sections: []model.Section{
			{
				ID:    "intro",
				Title: "Introduction",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitle,
						Title:  "Welcome",
					},
				},
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	// Check header
	if !strings.HasPrefix(output, "deck Test Presentation\n") {
		t.Error("output should start with deck title")
	}

	// Check metadata
	if !strings.Contains(output, "meta author John Doe") {
		t.Error("output should contain author")
	}
	if !strings.Contains(output, "meta date 2026-01-26") {
		t.Error("output should contain date")
	}

	// Check section
	if !strings.Contains(output, "section intro Introduction") {
		t.Error("output should contain section")
	}

	// Check slide
	if !strings.Contains(output, "slide s1 title") {
		t.Error("output should contain slide")
	}
	if !strings.Contains(output, "title Welcome") {
		t.Error("output should contain slide title")
	}
}

func TestTOONEncoderEncodeBlocks(t *testing.T) {
	deck := &model.Deck{
		Title: "Test",
		Sections: []model.Section{
			{
				ID: "sec1",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitleBody,
						Body: []model.Block{
							model.NewBullet("First item", 0),
							model.NewBullet("Nested item", 1),
							model.NewNumbered("Numbered item", 0),
							model.NewParagraph("A paragraph"),
							model.NewCode("fmt.Println()", "go"),
							model.NewImage("https://example.com/img.png", "Example"),
							model.NewQuote("A quote"),
							model.NewHeading("Subheading", 2),
						},
						Notes: []model.Block{
							{Text: "Speaker note"},
						},
					},
				},
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	expectations := []string{
		"bullet First item",
		"  bullet Nested item", // nested with extra indent
		"numbered Numbered item",
		"para A paragraph",
		"code go fmt.Println()",
		"image https://example.com/img.png Example",
		"quote A quote",
		"heading 2 Subheading",
		"note Speaker note",
	}

	for _, exp := range expectations {
		if !strings.Contains(output, exp) {
			t.Errorf("output should contain %q", exp)
		}
	}
}

func TestTOONEncoderEncodeAudio(t *testing.T) {
	transition := "fade"
	deck := &model.Deck{
		Title: "Test",
		Sections: []model.Section{
			{
				ID:    "sec1",
				Audio: model.NewFileAudio("/audio/section.mp3", 5*time.Minute),
				Slides: []model.Slide{
					{
						ID:         "s1",
						Layout:     model.LayoutTitleBody,
						Audio:      model.NewTTSAudio("Hello world", "en-US"),
						Transition: &transition,
					},
					{
						ID:     "s2",
						Layout: model.LayoutTitleBody,
						Audio:  model.NewURLAudio("https://example.com/audio.mp3", 2*time.Minute),
					},
					{
						ID:     "s3",
						Layout: model.LayoutTitleBody,
						Audio:  model.NewNotesAudio("en-GB"),
					},
				},
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	expectations := []string{
		"audio file /audio/section.mp3 duration=5m0s",
		"audio tts voice=en-US",
		"transition fade",
		"audio url https://example.com/audio.mp3 duration=2m0s",
		"audio notes voice=en-GB",
	}

	for _, exp := range expectations {
		if !strings.Contains(output, exp) {
			t.Errorf("output should contain %q\nGot:\n%s", exp, output)
		}
	}
}

func TestTOONEncoderEncodeDiff(t *testing.T) {
	diff := model.NewDiff("deck123")
	diff.AddChange(model.NewAddChange("sections/0/slides/0", "new slide"))
	diff.AddChange(model.NewRemoveChange("sections/1/slides/2", "old slide"))
	diff.AddChange(model.NewUpdateChange("sections/0/slides/0/title", "Old Title", "New Title"))
	diff.AddChange(model.NewMoveChange("sections/0/slides/1", "sections/1/slides/0"))

	encoder := NewTOONEncoder()
	output := encoder.EncodeDiff(diff)

	// Check header
	if !strings.HasPrefix(output, "plan deck deck123\n") {
		t.Error("output should start with plan deck")
	}

	// Check operations
	expectations := []string{
		"+ sections/0/slides/0",
		"- sections/1/slides/2",
		"~ sections/0/slides/0/title",
		"> sections/0/slides/1",
	}

	for _, exp := range expectations {
		if !strings.Contains(output, exp) {
			t.Errorf("output should contain %q", exp)
		}
	}

	// Check values
	if !strings.Contains(output, "+ new slide") {
		t.Error("output should contain new value for add")
	}
	if !strings.Contains(output, "- old slide") {
		t.Error("output should contain old value for remove")
	}
	if !strings.Contains(output, "- Old Title") {
		t.Error("output should contain old value for update")
	}
	if !strings.Contains(output, "+ New Title") {
		t.Error("output should contain new value for update")
	}
}

func TestTOONEncoderEmptyDeck(t *testing.T) {
	deck := &model.Deck{Title: "Empty"}
	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	if output != "deck Empty\n" {
		t.Errorf("empty deck output = %q, want 'deck Empty\\n'", output)
	}
}

func TestTOONEncoderEmptyDiff(t *testing.T) {
	diff := model.NewDiff("deck123")
	encoder := NewTOONEncoder()
	output := encoder.EncodeDiff(diff)

	if output != "plan deck deck123\n" {
		t.Errorf("empty diff output = %q", output)
	}
}

func TestFormatIsValid(t *testing.T) {
	if !FormatTOON.IsValid() {
		t.Error("FormatTOON should be valid")
	}
	if !FormatJSON.IsValid() {
		t.Error("FormatJSON should be valid")
	}
	if Format("invalid").IsValid() {
		t.Error("invalid format should not be valid")
	}
}

func TestTOONEncoderSubtitle(t *testing.T) {
	deck := &model.Deck{
		Title: "Test",
		Sections: []model.Section{
			{
				ID: "sec1",
				Slides: []model.Slide{
					{
						ID:       "s1",
						Layout:   model.LayoutTitle,
						Title:    "Main Title",
						Subtitle: "A Subtitle",
					},
				},
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	if !strings.Contains(output, "subtitle A Subtitle") {
		t.Error("output should contain subtitle")
	}
}

func TestTOONEncoderDescription(t *testing.T) {
	deck := &model.Deck{
		Title: "Test",
		Meta: model.Meta{
			Description: "A test presentation",
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	if !strings.Contains(output, "meta description A test presentation") {
		t.Error("output should contain description")
	}
}

func TestTOONEncoderSectionWithoutTitle(t *testing.T) {
	deck := &model.Deck{
		Title: "Test",
		Sections: []model.Section{
			{
				ID: "sec1",
				// No title
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	if !strings.Contains(output, "section sec1\n") {
		t.Errorf("section without title should just have ID, got:\n%s", output)
	}
}

func TestTOONEncoderHeadingWithoutLevel(t *testing.T) {
	deck := &model.Deck{
		Title: "Test",
		Sections: []model.Section{
			{
				ID: "sec1",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitleBody,
						Body: []model.Block{
							{Kind: model.BlockHeading, Text: "Heading", Level: 0},
						},
					},
				},
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	// Level 0 should not print level prefix
	if !strings.Contains(output, "heading Heading") {
		t.Errorf("heading level 0 should not have level prefix, got:\n%s", output)
	}
}

func TestTOONEncoderImageWithoutAlt(t *testing.T) {
	deck := &model.Deck{
		Title: "Test",
		Sections: []model.Section{
			{
				ID: "sec1",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutImage,
						Body: []model.Block{
							{Kind: model.BlockImage, URL: "https://example.com/img.png"},
						},
					},
				},
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	if !strings.Contains(output, "image https://example.com/img.png\n") {
		t.Errorf("image without alt should just have URL, got:\n%s", output)
	}
}

func TestTOONEncoderCodeWithoutLang(t *testing.T) {
	deck := &model.Deck{
		Title: "Test",
		Sections: []model.Section{
			{
				ID: "sec1",
				Slides: []model.Slide{
					{
						ID:     "s1",
						Layout: model.LayoutTitleBody,
						Body: []model.Block{
							{Kind: model.BlockCode, Text: "some code"},
						},
					},
				},
			},
		},
	}

	encoder := NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	if !strings.Contains(output, "code some code") {
		t.Errorf("code without lang should work, got:\n%s", output)
	}
}
