package model

// Block represents content within a slide.
type Block struct {
	Kind  BlockKind `json:"kind"`
	Text  string    `json:"text,omitempty"`
	Level int       `json:"level,omitempty"` // Bullet nesting level (0 = top level)
	Lang  string    `json:"lang,omitempty"`  // Code language
	URL   string    `json:"url,omitempty"`   // Image/link URL
	Alt   string    `json:"alt,omitempty"`   // Image alt text
}

// BlockKind identifies content type.
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

// BlockKinds returns all valid block kind values.
func BlockKinds() []BlockKind {
	return []BlockKind{
		BlockParagraph,
		BlockBullet,
		BlockNumbered,
		BlockCode,
		BlockImage,
		BlockQuote,
		BlockHeading,
	}
}

// IsValid returns true if the block kind is a recognized value.
func (k BlockKind) IsValid() bool {
	switch k {
	case BlockParagraph, BlockBullet, BlockNumbered, BlockCode,
		BlockImage, BlockQuote, BlockHeading:
		return true
	}
	return false
}

// NewParagraph creates a paragraph block.
func NewParagraph(text string) Block {
	return Block{Kind: BlockParagraph, Text: text}
}

// NewBullet creates a bullet block at the specified nesting level.
func NewBullet(text string, level int) Block {
	return Block{Kind: BlockBullet, Text: text, Level: level}
}

// NewNumbered creates a numbered list item block.
func NewNumbered(text string, level int) Block {
	return Block{Kind: BlockNumbered, Text: text, Level: level}
}

// NewCode creates a code block with the specified language.
func NewCode(code, lang string) Block {
	return Block{Kind: BlockCode, Text: code, Lang: lang}
}

// NewImage creates an image block.
func NewImage(url, alt string) Block {
	return Block{Kind: BlockImage, URL: url, Alt: alt}
}

// NewQuote creates a quote block.
func NewQuote(text string) Block {
	return Block{Kind: BlockQuote, Text: text}
}

// NewHeading creates a heading block.
func NewHeading(text string, level int) Block {
	return Block{Kind: BlockHeading, Text: text, Level: level}
}
