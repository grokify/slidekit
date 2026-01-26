package model

// Slide represents a single slide.
type Slide struct {
	ID         string  `json:"id"`
	Layout     Layout  `json:"layout"`
	Title      string  `json:"title,omitempty"`
	Subtitle   string  `json:"subtitle,omitempty"`
	Body       []Block `json:"body,omitempty"`
	Notes      []Block `json:"notes,omitempty"`      // Speaker notes
	Audio      *Audio  `json:"audio,omitempty"`      // Slide-level audio
	Transition *string `json:"transition,omitempty"` // Reveal.js transitions
	Background *string `json:"background,omitempty"`
}

// Layout identifies slide layout type.
type Layout string

const (
	LayoutTitle       Layout = "title"         // Title slide
	LayoutTitleBody   Layout = "title_body"    // Title + body content
	LayoutTitleTwoCol Layout = "title_two_col" // Title + two columns
	LayoutSection     Layout = "section"       // Section divider
	LayoutBlank       Layout = "blank"         // No predefined structure
	LayoutImage       Layout = "image"         // Full-bleed image
	LayoutComparison  Layout = "comparison"    // Side-by-side comparison
)

// Layouts returns all valid layout values.
func Layouts() []Layout {
	return []Layout{
		LayoutTitle,
		LayoutTitleBody,
		LayoutTitleTwoCol,
		LayoutSection,
		LayoutBlank,
		LayoutImage,
		LayoutComparison,
	}
}

// IsValid returns true if the layout is a recognized value.
func (l Layout) IsValid() bool {
	switch l {
	case LayoutTitle, LayoutTitleBody, LayoutTitleTwoCol, LayoutSection,
		LayoutBlank, LayoutImage, LayoutComparison:
		return true
	}
	return false
}

// HasTitle returns true if the slide has a title.
func (s *Slide) HasTitle() bool {
	return s.Title != ""
}

// HasBody returns true if the slide has body content.
func (s *Slide) HasBody() bool {
	return len(s.Body) > 0
}

// HasNotes returns true if the slide has speaker notes.
func (s *Slide) HasNotes() bool {
	return len(s.Notes) > 0
}

// NotesText returns speaker notes as plain text.
func (s *Slide) NotesText() string {
	var text string
	for i, block := range s.Notes {
		if i > 0 {
			text += "\n"
		}
		text += block.Text
	}
	return text
}

// BulletCount returns the number of bullet points in the body.
func (s *Slide) BulletCount() int {
	count := 0
	for _, block := range s.Body {
		if block.Kind == BlockBullet || block.Kind == BlockNumbered {
			count++
		}
	}
	return count
}
