// Package marp implements the Marp Markdown backend for slidekit.
package marp

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/slidekit/model"
)

// Reader parses Marp Markdown files into the canonical slide model.
type Reader struct{}

// NewReader creates a new Marp reader.
func NewReader() *Reader {
	return &Reader{}
}

// ReadFile reads a Marp Markdown file and returns a Deck.
func (r *Reader) ReadFile(path string) (*model.Deck, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}
	return r.Parse(string(data))
}

// Read implements the Backend interface for reading from a Ref.
func (r *Reader) Read(ctx context.Context, ref model.Ref) (*model.Deck, error) {
	if ref.Path == "" {
		return nil, fmt.Errorf("marp backend requires a file path")
	}
	return r.ReadFile(ref.Path)
}

// Parse parses Marp Markdown content into a Deck.
func (r *Reader) Parse(content string) (*model.Deck, error) {
	// Parse frontmatter
	frontmatter, body := parseFrontmatter(content)

	// Split into raw slides
	rawSlides := splitSlides(body)

	// Parse each raw slide
	var slides []parsedSlide
	for _, raw := range rawSlides {
		slides = append(slides, parseRawSlide(raw))
	}

	// Group into sections based on section-divider slides
	deck := buildDeck(frontmatter, slides)

	return deck, nil
}

// parsedSlide holds intermediate parsed state for a single slide.
type parsedSlide struct {
	directives map[string]string
	notes      []noteBlock
	content    string // Remaining markdown/HTML content
	rawContent string // Original raw content
}

// noteBlock represents a speaker note.
type noteBlock struct {
	text string
}

// Frontmatter holds parsed YAML frontmatter.
type Frontmatter struct {
	Marp     bool
	Theme    string
	Paginate bool
	Style    string
	Custom   map[string]string
}

// parseFrontmatter extracts YAML frontmatter from content.
func parseFrontmatter(content string) (Frontmatter, string) {
	fm := Frontmatter{Custom: make(map[string]string)}

	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return fm, content
	}

	// Find closing ---
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return fm, content
	}

	fmContent := strings.TrimSpace(rest[:idx])
	body := rest[idx+4:] // Skip past \n---

	// Parse YAML-like frontmatter (simple key: value parsing)
	lines := strings.Split(fmContent, "\n")
	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			i++
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			i++
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "marp":
			fm.Marp = value == "true"
		case "theme":
			fm.Theme = value
		case "paginate":
			fm.Paginate = value == "true"
		case "style":
			// Handle multiline style block (indicated by |)
			if value == "|" {
				var styleLines []string
				i++
				for i < len(lines) {
					sl := lines[i]
					// Multiline block continues while indented
					if len(sl) > 0 && (sl[0] == ' ' || sl[0] == '\t') {
						styleLines = append(styleLines, sl)
					} else {
						break
					}
					i++
				}
				fm.Style = strings.Join(styleLines, "\n")
				continue
			}
			fm.Style = value
		default:
			fm.Custom[key] = value
		}
		i++
	}

	return fm, body
}

// splitSlides splits the body into raw slide strings on --- separators.
// It respects code blocks (``` fences) and does not split within them.
func splitSlides(body string) []string {
	lines := strings.Split(body, "\n")
	var slides []string
	var current []string
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track code block state
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
		}

		// Check for slide separator (only when not in code block)
		if !inCodeBlock && trimmed == "---" {
			slideContent := strings.Join(current, "\n")
			if strings.TrimSpace(slideContent) != "" {
				slides = append(slides, slideContent)
			}
			current = nil
			continue
		}

		current = append(current, line)
	}

	// Don't forget the last slide
	if len(current) > 0 {
		slideContent := strings.Join(current, "\n")
		if strings.TrimSpace(slideContent) != "" {
			slides = append(slides, slideContent)
		}
	}

	return slides
}

var (
	// Matches directive comments like <!-- _class: section-divider -->
	reDirective = regexp.MustCompile(`<!--\s*(_\w+)\s*:\s*(.+?)\s*-->`)

	// Matches speaker note comments (multi-line)
	reNoteBlock = regexp.MustCompile(`(?s)<!--\s*\n(.*?)\n\s*-->`)

	// Matches [PAUSE:ms] markers
	rePause = regexp.MustCompile(`\[PAUSE:(\d+)\]`)

	// Matches <break time="ms"/> SSML markers
	reBreak = regexp.MustCompile(`<break\s+time="(\d+)ms"\s*/?>`)
)

// parseRawSlide parses directives, notes, and content from a raw slide string.
func parseRawSlide(raw string) parsedSlide {
	ps := parsedSlide{
		directives: make(map[string]string),
		rawContent: raw,
	}

	remaining := raw

	// Extract directives (<!-- _key: value -->)
	for _, match := range reDirective.FindAllStringSubmatch(remaining, -1) {
		key := match[1]
		value := match[2]
		ps.directives[key] = value
	}
	remaining = reDirective.ReplaceAllString(remaining, "")

	// Extract speaker notes (multi-line <!-- ... --> that are not directives)
	for _, match := range reNoteBlock.FindAllStringSubmatch(remaining, -1) {
		noteText := match[1]
		// Skip if it looks like a script or HTML tag
		if strings.Contains(noteText, "<script") {
			continue
		}
		nb := noteBlock{text: cleanNoteText(noteText)}
		ps.notes = append(ps.notes, nb)
	}
	remaining = reNoteBlock.ReplaceAllString(remaining, "")

	ps.content = strings.TrimSpace(remaining)
	return ps
}

// cleanNoteText removes pause/break markers and cleans whitespace from note text.
func cleanNoteText(text string) string {
	// Remove PAUSE markers
	text = rePause.ReplaceAllString(text, "")
	// Remove SSML break markers
	text = reBreak.ReplaceAllString(text, "")
	// Clean up whitespace
	lines := strings.Split(text, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return strings.Join(cleaned, " ")
}

// buildDeck constructs a Deck from frontmatter and parsed slides.
func buildDeck(fm Frontmatter, slides []parsedSlide) *model.Deck {
	deck := &model.Deck{
		Meta: model.Meta{
			Custom: make(map[string]string),
		},
	}

	// Apply frontmatter
	if fm.Theme != "" {
		deck.Theme = &model.Theme{Name: fm.Theme}
		if fm.Style != "" {
			deck.Theme.SetCustom("style", fm.Style)
		}
	}
	for k, v := range fm.Custom {
		deck.Meta.Custom[k] = v
	}

	// Group slides into sections based on section-divider slides
	type sectionGroup struct {
		title  string
		slides []parsedSlide
	}

	var groups []sectionGroup
	currentGroup := sectionGroup{title: "default"}

	for _, ps := range slides {
		class := ps.directives["_class"]

		if class == "section-divider" || class == "lead" {
			// This slide starts a new section
			// Save current group if it has slides
			if len(currentGroup.slides) > 0 {
				groups = append(groups, currentGroup)
			}
			// Extract section title from the slide content
			title := extractSectionTitle(ps.content)
			currentGroup = sectionGroup{title: title}
			// The section-divider slide itself becomes the first slide in this section
			currentGroup.slides = append(currentGroup.slides, ps)
		} else {
			currentGroup.slides = append(currentGroup.slides, ps)
		}
	}
	// Don't forget the last group
	if len(currentGroup.slides) > 0 {
		groups = append(groups, currentGroup)
	}

	// Convert groups to sections
	for i, group := range groups {
		section := model.Section{
			ID:    fmt.Sprintf("section-%d", i),
			Title: group.title,
		}

		for j, ps := range group.slides {
			slide := convertToSlide(ps, i, j)
			section.Slides = append(section.Slides, slide)
		}

		deck.Sections = append(deck.Sections, section)
	}

	// Set deck title from first title slide
	if len(deck.Sections) > 0 && len(deck.Sections[0].Slides) > 0 {
		deck.Title = deck.Sections[0].Slides[0].Title
	}

	return deck
}

// extractSectionTitle extracts the title from a section-divider slide.
// Prefers H2 over H1 since section dividers often have "Section N" as H1
// and the actual topic as H2.
func extractSectionTitle(content string) string {
	lines := strings.Split(content, "\n")
	var h1, h2 string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if h1 == "" && strings.HasPrefix(trimmed, "# ") {
			h1 = strings.TrimPrefix(trimmed, "# ")
		} else if h2 == "" && strings.HasPrefix(trimmed, "## ") {
			h2 = strings.TrimPrefix(trimmed, "## ")
		}
	}
	// Prefer H2 if present (more descriptive)
	if h2 != "" {
		return h2
	}
	if h1 != "" {
		return h1
	}
	return "untitled"
}

// convertToSlide converts a parsedSlide into a model.Slide.
func convertToSlide(ps parsedSlide, sectionIdx, slideIdx int) model.Slide {
	slide := model.Slide{
		ID: fmt.Sprintf("s%d-%d", sectionIdx, slideIdx),
	}

	// Determine layout from directives
	class := ps.directives["_class"]
	switch class {
	case "section-divider":
		slide.Layout = model.LayoutSection
	case "lead":
		slide.Layout = model.LayoutTitle
	default:
		slide.Layout = model.LayoutTitleBody
	}

	// Parse content into blocks
	parseSlideContent(&slide, ps.content)

	// Add speaker notes
	for _, nb := range ps.notes {
		if nb.text != "" {
			slide.Notes = append(slide.Notes, model.Block{
				Kind: model.BlockParagraph,
				Text: nb.text,
			})
		}
	}

	// Detect layout from HTML structure
	if containsColumns(ps.content) {
		slide.Layout = model.LayoutTitleTwoCol
	}

	return slide
}

// containsColumns checks if the content has multi-column HTML layout.
func containsColumns(content string) bool {
	return strings.Contains(content, `class="columns"`) ||
		strings.Contains(content, "grid-template-columns")
}

// parseSlideContent parses markdown/HTML content into a slide's title, subtitle, and body.
func parseSlideContent(slide *model.Slide, content string) {
	lines := strings.Split(content, "\n")
	var bodyLines []string
	inCodeBlock := false
	codeBlockLang := ""
	var codeBlockLines []string
	inHTMLBlock := false
	var htmlBlockLines []string
	htmlBlockDepth := 0

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Handle code blocks
		if strings.HasPrefix(trimmed, "```") {
			if !inCodeBlock {
				inCodeBlock = true
				codeBlockLang = strings.TrimPrefix(trimmed, "```")
				codeBlockLines = nil
				continue
			}
			// End of code block
			inCodeBlock = false
			slide.Body = append(slide.Body, model.Block{
				Kind: model.BlockCode,
				Text: strings.Join(codeBlockLines, "\n"),
				Lang: codeBlockLang,
			})
			continue
		}
		if inCodeBlock {
			codeBlockLines = append(codeBlockLines, line)
			continue
		}

		// Handle HTML blocks (div, script, etc.)
		if !inHTMLBlock && isHTMLBlockStart(trimmed) {
			inHTMLBlock = true
			htmlBlockDepth = 1
			htmlBlockLines = []string{line}
			// Check if self-closing on same line
			htmlBlockDepth += countHTMLOpens(line) - 1 // -1 for the one we already counted
			htmlBlockDepth -= countHTMLCloses(line)
			if htmlBlockDepth <= 0 {
				inHTMLBlock = false
				slide.Body = append(slide.Body, model.Block{
					Kind: model.BlockParagraph,
					Text: strings.Join(htmlBlockLines, "\n"),
				})
				htmlBlockLines = nil
			}
			continue
		}
		if inHTMLBlock {
			htmlBlockLines = append(htmlBlockLines, line)
			htmlBlockDepth += countHTMLOpens(line)
			htmlBlockDepth -= countHTMLCloses(line)
			if htmlBlockDepth <= 0 {
				inHTMLBlock = false
				slide.Body = append(slide.Body, model.Block{
					Kind: model.BlockParagraph,
					Text: strings.Join(htmlBlockLines, "\n"),
				})
				htmlBlockLines = nil
			}
			continue
		}

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Handle headings
		if strings.HasPrefix(trimmed, "# ") {
			title := strings.TrimPrefix(trimmed, "# ")
			if slide.Title == "" {
				slide.Title = title
			} else {
				slide.Body = append(slide.Body, model.NewHeading(title, 1))
			}
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			subtitle := strings.TrimPrefix(trimmed, "## ")
			if slide.Title == "" {
				slide.Title = subtitle
			} else if slide.Subtitle == "" {
				slide.Subtitle = subtitle
			} else {
				slide.Body = append(slide.Body, model.NewHeading(subtitle, 2))
			}
			continue
		}
		if strings.HasPrefix(trimmed, "### ") {
			slide.Body = append(slide.Body, model.NewHeading(
				strings.TrimPrefix(trimmed, "### "), 3))
			continue
		}

		// Handle bullet points
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			level := countIndentLevel(line)
			text := strings.TrimSpace(trimmed[2:])
			slide.Body = append(slide.Body, model.NewBullet(text, level))
			continue
		}

		// Handle numbered lists
		if isNumberedItem(trimmed) {
			level := countIndentLevel(line)
			text := extractNumberedText(trimmed)
			slide.Body = append(slide.Body, model.NewNumbered(text, level))
			continue
		}

		// Handle blockquotes
		if strings.HasPrefix(trimmed, "> ") {
			text := strings.TrimPrefix(trimmed, "> ")
			slide.Body = append(slide.Body, model.NewQuote(text))
			continue
		}

		// Handle images
		if strings.HasPrefix(trimmed, "![") {
			alt, url := parseImageLine(trimmed)
			if url != "" {
				slide.Body = append(slide.Body, model.NewImage(url, alt))
				continue
			}
		}

		// Handle tables (collect all table lines as a single block)
		if strings.HasPrefix(trimmed, "|") {
			var tableLines []string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "|") {
				tableLines = append(tableLines, lines[i])
				i++
			}
			i-- // Back up one since the loop will increment
			slide.Body = append(slide.Body, model.Block{
				Kind: model.BlockParagraph,
				Text: strings.Join(tableLines, "\n"),
			})
			continue
		}

		// Default: paragraph text
		bodyLines = append(bodyLines, trimmed)
	}

	// Flush remaining body lines as paragraphs
	if len(bodyLines) > 0 {
		text := strings.Join(bodyLines, "\n")
		if text != "" {
			slide.Body = append(slide.Body, model.NewParagraph(text))
		}
	}
}

var reHTMLOpen = regexp.MustCompile(`<(div|section|script|table|ol|ul)\b`)
var reHTMLClose = regexp.MustCompile(`</(div|section|script|table|ol|ul)>`)

// isHTMLBlockStart checks if a line starts an HTML block.
func isHTMLBlockStart(line string) bool {
	return strings.HasPrefix(line, "<div") ||
		strings.HasPrefix(line, "<section") ||
		strings.HasPrefix(line, "<script") ||
		strings.HasPrefix(line, "<ol") ||
		strings.HasPrefix(line, "<table")
}

// countHTMLOpens counts opening HTML tags in a line.
func countHTMLOpens(line string) int {
	return len(reHTMLOpen.FindAllString(line, -1))
}

// countHTMLCloses counts closing HTML tags in a line.
func countHTMLCloses(line string) int {
	return len(reHTMLClose.FindAllString(line, -1))
}

// countIndentLevel counts the bullet nesting level based on leading whitespace.
func countIndentLevel(line string) int {
	spaces := 0
	for _, ch := range line {
		switch ch {
		case ' ':
			spaces++
		case '\t':
			spaces += 4
		default:
			return spaces / 4
		}
	}
	return spaces / 4 // 4 spaces per level
}

var reNumbered = regexp.MustCompile(`^(\d+)\.\s+(.+)$`)

// isNumberedItem checks if a line is a numbered list item.
func isNumberedItem(line string) bool {
	return reNumbered.MatchString(line)
}

// extractNumberedText extracts the text from a numbered list item.
func extractNumberedText(line string) string {
	match := reNumbered.FindStringSubmatch(line)
	if len(match) >= 3 {
		return match[2]
	}
	return line
}

var reImage = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

// parseImageLine extracts alt text and URL from a markdown image.
func parseImageLine(line string) (alt, url string) {
	match := reImage.FindStringSubmatch(line)
	if len(match) >= 3 {
		return match[1], match[2]
	}
	return "", ""
}
