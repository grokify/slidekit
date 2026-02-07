package model

import (
	"testing"
	"time"
)

// Deck tests

func TestDeckSlideCount(t *testing.T) {
	deck := &Deck{
		Sections: []Section{
			{Slides: []Slide{{}, {}, {}}},
			{Slides: []Slide{{}, {}}},
		},
	}
	if got := deck.SlideCount(); got != 5 {
		t.Errorf("SlideCount() = %d, want 5", got)
	}
}

func TestDeckSlideCountEmpty(t *testing.T) {
	deck := &Deck{}
	if got := deck.SlideCount(); got != 0 {
		t.Errorf("SlideCount() = %d, want 0", got)
	}
}

func TestDeckAllSlides(t *testing.T) {
	deck := &Deck{
		Sections: []Section{
			{Slides: []Slide{{ID: "s1"}, {ID: "s2"}}},
			{Slides: []Slide{{ID: "s3"}}},
		},
	}
	slides := deck.AllSlides()
	if len(slides) != 3 {
		t.Errorf("AllSlides() returned %d slides, want 3", len(slides))
	}
	if slides[0].ID != "s1" || slides[1].ID != "s2" || slides[2].ID != "s3" {
		t.Error("AllSlides() returned slides in wrong order")
	}
}

func TestDeckFindSlide(t *testing.T) {
	deck := &Deck{
		Sections: []Section{
			{Slides: []Slide{{ID: "s1"}, {ID: "s2"}}},
			{Slides: []Slide{{ID: "s3"}}},
		},
	}

	// Find existing slide
	slide := deck.FindSlide("s2")
	if slide == nil {
		t.Fatal("FindSlide(s2) returned nil")
	}
	if slide.ID != "s2" {
		t.Errorf("FindSlide(s2) returned slide with ID %s", slide.ID)
	}

	// Find non-existing slide
	if deck.FindSlide("nonexistent") != nil {
		t.Error("FindSlide(nonexistent) should return nil")
	}
}

func TestDeckFindSection(t *testing.T) {
	deck := &Deck{
		Sections: []Section{
			{ID: "sec1", Title: "Section 1"},
			{ID: "sec2", Title: "Section 2"},
		},
	}

	// Find existing section
	section := deck.FindSection("sec2")
	if section == nil {
		t.Fatal("FindSection(sec2) returned nil")
	}
	if section.Title != "Section 2" {
		t.Errorf("FindSection(sec2) returned section with title %s", section.Title)
	}

	// Find non-existing section
	if deck.FindSection("nonexistent") != nil {
		t.Error("FindSection(nonexistent) should return nil")
	}
}

func TestDeckTotalDuration(t *testing.T) {
	deck := &Deck{
		Sections: []Section{
			{Audio: &Audio{Duration: 5 * time.Minute}},
			{Audio: &Audio{Duration: 3 * time.Minute}},
		},
	}
	if got := deck.TotalDuration(); got != 8*time.Minute {
		t.Errorf("TotalDuration() = %v, want 8m", got)
	}
}

// Section tests

func TestSectionSlideCount(t *testing.T) {
	section := &Section{Slides: []Slide{{}, {}, {}}}
	if got := section.SlideCount(); got != 3 {
		t.Errorf("SlideCount() = %d, want 3", got)
	}
}

func TestSectionFindSlide(t *testing.T) {
	section := &Section{
		Slides: []Slide{{ID: "s1"}, {ID: "s2"}},
	}

	slide := section.FindSlide("s1")
	if slide == nil || slide.ID != "s1" {
		t.Error("FindSlide(s1) failed")
	}

	if section.FindSlide("nonexistent") != nil {
		t.Error("FindSlide(nonexistent) should return nil")
	}
}

func TestSectionTotalDuration(t *testing.T) {
	// Section-level audio takes precedence
	section := &Section{
		Audio: &Audio{Duration: 10 * time.Minute},
		Slides: []Slide{
			{Audio: &Audio{Duration: 2 * time.Minute}},
			{Audio: &Audio{Duration: 3 * time.Minute}},
		},
	}
	if got := section.TotalDuration(); got != 10*time.Minute {
		t.Errorf("TotalDuration() with section audio = %v, want 10m", got)
	}

	// Sum slide-level audio when no section audio
	section2 := &Section{
		Slides: []Slide{
			{Audio: &Audio{Duration: 2 * time.Minute}},
			{Audio: &Audio{Duration: 3 * time.Minute}},
		},
	}
	if got := section2.TotalDuration(); got != 5*time.Minute {
		t.Errorf("TotalDuration() with slide audio = %v, want 5m", got)
	}
}

func TestSectionHasAudio(t *testing.T) {
	// Section-level audio
	s1 := &Section{Audio: &Audio{Source: AudioSourceFile}}
	if !s1.HasAudio() {
		t.Error("HasAudio() should be true with section audio")
	}

	// Slide-level audio
	s2 := &Section{
		Slides: []Slide{{Audio: &Audio{Source: AudioSourceTTS}}},
	}
	if !s2.HasAudio() {
		t.Error("HasAudio() should be true with slide audio")
	}

	// No audio
	s3 := &Section{Slides: []Slide{{}}}
	if s3.HasAudio() {
		t.Error("HasAudio() should be false with no audio")
	}
}

// Slide tests

func TestSlideHasTitle(t *testing.T) {
	if (&Slide{Title: "Test"}).HasTitle() != true {
		t.Error("HasTitle() should be true")
	}
	if (&Slide{}).HasTitle() != false {
		t.Error("HasTitle() should be false")
	}
}

func TestSlideHasBody(t *testing.T) {
	if (&Slide{Body: []Block{{}}}).HasBody() != true {
		t.Error("HasBody() should be true")
	}
	if (&Slide{}).HasBody() != false {
		t.Error("HasBody() should be false")
	}
}

func TestSlideHasNotes(t *testing.T) {
	if (&Slide{Notes: []Block{{}}}).HasNotes() != true {
		t.Error("HasNotes() should be true")
	}
	if (&Slide{}).HasNotes() != false {
		t.Error("HasNotes() should be false")
	}
}

func TestSlideNotesText(t *testing.T) {
	slide := &Slide{
		Notes: []Block{
			{Text: "First note"},
			{Text: "Second note"},
		},
	}
	expected := "First note\nSecond note"
	if got := slide.NotesText(); got != expected {
		t.Errorf("NotesText() = %q, want %q", got, expected)
	}
}

func TestSlideBulletCount(t *testing.T) {
	slide := &Slide{
		Body: []Block{
			{Kind: BlockBullet},
			{Kind: BlockParagraph},
			{Kind: BlockBullet},
			{Kind: BlockNumbered},
		},
	}
	if got := slide.BulletCount(); got != 3 {
		t.Errorf("BulletCount() = %d, want 3", got)
	}
}

func TestLayoutIsValid(t *testing.T) {
	for _, l := range Layouts() {
		if !l.IsValid() {
			t.Errorf("Layout %s should be valid", l)
		}
	}
	if Layout("invalid").IsValid() {
		t.Error("Layout 'invalid' should not be valid")
	}
}

// Block tests

func TestBlockKindIsValid(t *testing.T) {
	for _, k := range BlockKinds() {
		if !k.IsValid() {
			t.Errorf("BlockKind %s should be valid", k)
		}
	}
	if BlockKind("invalid").IsValid() {
		t.Error("BlockKind 'invalid' should not be valid")
	}
}

func TestBlockConstructors(t *testing.T) {
	tests := []struct {
		name  string
		block Block
		kind  BlockKind
	}{
		{"NewParagraph", NewParagraph("text"), BlockParagraph},
		{"NewBullet", NewBullet("text", 1), BlockBullet},
		{"NewNumbered", NewNumbered("text", 0), BlockNumbered},
		{"NewCode", NewCode("code", "go"), BlockCode},
		{"NewImage", NewImage("url", "alt"), BlockImage},
		{"NewQuote", NewQuote("quote"), BlockQuote},
		{"NewHeading", NewHeading("heading", 2), BlockHeading},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.block.Kind != tt.kind {
				t.Errorf("%s created block with kind %s, want %s", tt.name, tt.block.Kind, tt.kind)
			}
		})
	}
}

func TestNewBulletLevel(t *testing.T) {
	b := NewBullet("test", 2)
	if b.Level != 2 {
		t.Errorf("NewBullet level = %d, want 2", b.Level)
	}
}

func TestNewCodeLang(t *testing.T) {
	b := NewCode("fmt.Println()", "go")
	if b.Lang != "go" {
		t.Errorf("NewCode lang = %s, want go", b.Lang)
	}
	if b.Text != "fmt.Println()" {
		t.Errorf("NewCode text = %s, want fmt.Println()", b.Text)
	}
}

func TestNewImageURLAlt(t *testing.T) {
	b := NewImage("https://example.com/img.png", "Example image")
	if b.URL != "https://example.com/img.png" {
		t.Errorf("NewImage URL = %s", b.URL)
	}
	if b.Alt != "Example image" {
		t.Errorf("NewImage Alt = %s", b.Alt)
	}
}

// Audio tests

func TestAudioSourceIsValid(t *testing.T) {
	for _, s := range AudioSources() {
		if !s.IsValid() {
			t.Errorf("AudioSource %s should be valid", s)
		}
	}
	if AudioSource("invalid").IsValid() {
		t.Error("AudioSource 'invalid' should not be valid")
	}
}

func TestAudioSourceNeedsTTS(t *testing.T) {
	if !AudioSourceTTS.NeedsTTS() {
		t.Error("AudioSourceTTS.NeedsTTS() should be true")
	}
	if !AudioSourceNotes.NeedsTTS() {
		t.Error("AudioSourceNotes.NeedsTTS() should be true")
	}
	if AudioSourceFile.NeedsTTS() {
		t.Error("AudioSourceFile.NeedsTTS() should be false")
	}
	if AudioSourceURL.NeedsTTS() {
		t.Error("AudioSourceURL.NeedsTTS() should be false")
	}
}

func TestAudioFactoryMethods(t *testing.T) {
	file := NewFileAudio("/path/to/audio.mp3", 5*time.Minute)
	if file.Source != AudioSourceFile || file.Path != "/path/to/audio.mp3" {
		t.Error("NewFileAudio failed")
	}

	url := NewURLAudio("https://example.com/audio.mp3", 3*time.Minute)
	if url.Source != AudioSourceURL || url.URL != "https://example.com/audio.mp3" {
		t.Error("NewURLAudio failed")
	}

	tts := NewTTSAudio("Hello world", "en-US")
	if tts.Source != AudioSourceTTS || tts.Script != "Hello world" || tts.Voice != "en-US" {
		t.Error("NewTTSAudio failed")
	}

	notes := NewNotesAudio("en-GB")
	if notes.Source != AudioSourceNotes || notes.Voice != "en-GB" {
		t.Error("NewNotesAudio failed")
	}
}

func TestAudioHasContent(t *testing.T) {
	// Nil audio
	var nilAudio *Audio
	if nilAudio.HasContent() {
		t.Error("nil audio should not have content")
	}

	// File audio with path
	if !(&Audio{Source: AudioSourceFile, Path: "/path"}).HasContent() {
		t.Error("file audio with path should have content")
	}
	if (&Audio{Source: AudioSourceFile}).HasContent() {
		t.Error("file audio without path should not have content")
	}

	// URL audio
	if !(&Audio{Source: AudioSourceURL, URL: "https://example.com"}).HasContent() {
		t.Error("URL audio should have content")
	}

	// TTS audio with script
	if !(&Audio{Source: AudioSourceTTS, Script: "Hello"}).HasContent() {
		t.Error("TTS audio with script should have content")
	}

	// Notes audio always has content
	if !(&Audio{Source: AudioSourceNotes}).HasContent() {
		t.Error("Notes audio should have content")
	}
}

// Diff tests

func TestDiffIsEmpty(t *testing.T) {
	d := NewDiff("deck1")
	if !d.IsEmpty() {
		t.Error("NewDiff should be empty")
	}

	d.AddChange(NewAddChange("path", "value"))
	if d.IsEmpty() {
		t.Error("Diff with changes should not be empty")
	}
}

func TestDiffChangeCount(t *testing.T) {
	d := NewDiff("deck1")
	d.AddChange(NewAddChange("p1", "v1"))
	d.AddChange(NewRemoveChange("p2", "v2"))
	if got := d.ChangeCount(); got != 2 {
		t.Errorf("ChangeCount() = %d, want 2", got)
	}
}

func TestDiffCountByOp(t *testing.T) {
	d := NewDiff("deck1")
	d.AddChange(NewAddChange("p1", "v1"))
	d.AddChange(NewAddChange("p2", "v2"))
	d.AddChange(NewRemoveChange("p3", "v3"))
	d.AddChange(NewUpdateChange("p4", "old", "new"))

	counts := d.CountByOp()
	if counts[ChangeAdd] != 2 {
		t.Errorf("add count = %d, want 2", counts[ChangeAdd])
	}
	if counts[ChangeRemove] != 1 {
		t.Errorf("remove count = %d, want 1", counts[ChangeRemove])
	}
	if counts[ChangeUpdate] != 1 {
		t.Errorf("update count = %d, want 1", counts[ChangeUpdate])
	}
}

func TestDiffFilterMethods(t *testing.T) {
	d := NewDiff("deck1")
	d.AddChange(NewAddChange("p1", "v1"))
	d.AddChange(NewRemoveChange("p2", "v2"))
	d.AddChange(NewUpdateChange("p3", "old", "new"))
	d.AddChange(NewMoveChange("p4", "p5"))

	if len(d.AddChanges()) != 1 {
		t.Error("AddChanges() wrong count")
	}
	if len(d.RemoveChanges()) != 1 {
		t.Error("RemoveChanges() wrong count")
	}
	if len(d.UpdateChanges()) != 1 {
		t.Error("UpdateChanges() wrong count")
	}
	if len(d.MoveChanges()) != 1 {
		t.Error("MoveChanges() wrong count")
	}
}

func TestChangeOpIsValid(t *testing.T) {
	for _, op := range ChangeOps() {
		if !op.IsValid() {
			t.Errorf("ChangeOp %s should be valid", op)
		}
	}
	if ChangeOp("invalid").IsValid() {
		t.Error("ChangeOp 'invalid' should not be valid")
	}
}

func TestNewChangeHelpers(t *testing.T) {
	add := NewAddChange("path", "value")
	if add.Op != ChangeAdd || add.Path != "path" || add.NewValue != "value" {
		t.Error("NewAddChange failed")
	}

	remove := NewRemoveChange("path", "value")
	if remove.Op != ChangeRemove || remove.OldValue != "value" {
		t.Error("NewRemoveChange failed")
	}

	update := NewUpdateChange("path", "old", "new")
	if update.Op != ChangeUpdate || update.OldValue != "old" || update.NewValue != "new" {
		t.Error("NewUpdateChange failed")
	}

	move := NewMoveChange("from", "to")
	if move.Op != ChangeMove || move.Path != "from" || move.NewValue != "to" {
		t.Error("NewMoveChange failed")
	}
}

// Theme tests

func TestDefaultTheme(t *testing.T) {
	theme := DefaultTheme()
	if theme.Name != "default" {
		t.Errorf("DefaultTheme name = %s, want default", theme.Name)
	}
	if theme.Primary == "" || theme.Background == "" {
		t.Error("DefaultTheme should have colors set")
	}
}

func TestDarkTheme(t *testing.T) {
	theme := DarkTheme()
	if theme.Name != "dark" {
		t.Errorf("DarkTheme name = %s, want dark", theme.Name)
	}
	if theme.Background != "#121212" {
		t.Errorf("DarkTheme background = %s, want #121212", theme.Background)
	}
}

func TestThemeGetCustom(t *testing.T) {
	// Nil theme
	var nilTheme *Theme
	if nilTheme.GetCustom("key", "default") != "default" {
		t.Error("nil theme should return default")
	}

	// Theme without custom map
	theme := &Theme{}
	if theme.GetCustom("key", "default") != "default" {
		t.Error("empty theme should return default")
	}

	// Theme with custom value
	theme.SetCustom("key", "value")
	if theme.GetCustom("key", "default") != "value" {
		t.Error("theme should return custom value")
	}

	// Theme with different key
	if theme.GetCustom("other", "default") != "default" {
		t.Error("theme should return default for missing key")
	}
}

func TestThemeSetCustom(t *testing.T) {
	theme := &Theme{}
	theme.SetCustom("key1", "value1")
	theme.SetCustom("key2", "value2")

	if theme.Custom["key1"] != "value1" || theme.Custom["key2"] != "value2" {
		t.Error("SetCustom failed")
	}
}
