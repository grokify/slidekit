package ops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/grokify/slidekit/format"
	"github.com/grokify/slidekit/model"
)

// ErrSlideNotFound is returned when a slide cannot be found.
var ErrSlideNotFound = errors.New("slide not found")

// SlideInfo provides summary information about a slide.
type SlideInfo struct {
	ID        string `json:"id"`
	SectionID string `json:"section_id"`
	Title     string `json:"title"`
	Layout    string `json:"layout"`
}

// ListSlidesResult contains the result of a ListSlides operation.
type ListSlidesResult struct {
	Slides []SlideInfo
	Output string
}

// ListSlides returns all slides in a deck.
func ListSlides(ctx context.Context, ref model.Ref, f format.Format) (*ListSlidesResult, error) {
	result, err := ReadDeck(ctx, ref, ReadOptions{})
	if err != nil {
		return nil, err
	}

	slides := make([]SlideInfo, 0)
	for _, section := range result.Deck.Sections {
		for _, slide := range section.Slides {
			slides = append(slides, SlideInfo{
				ID:        slide.ID,
				SectionID: section.ID,
				Title:     slide.Title,
				Layout:    string(slide.Layout),
			})
		}
	}

	output, err := encodeSlideList(slides, f)
	if err != nil {
		return nil, err
	}

	return &ListSlidesResult{
		Slides: slides,
		Output: output,
	}, nil
}

// ListSlidesFromPath is a convenience function that detects the backend.
func ListSlidesFromPath(ctx context.Context, path string, f format.Format) (*ListSlidesResult, error) {
	backendName := DetectBackend(path)
	ref := model.Ref{
		Backend: backendName,
		Path:    path,
	}
	return ListSlides(ctx, ref, f)
}

// GetSlideResult contains the result of a GetSlide operation.
type GetSlideResult struct {
	Slide  *model.Slide
	Output string
}

// GetSlide returns a single slide by ID.
func GetSlide(ctx context.Context, ref model.Ref, slideID string, f format.Format) (*GetSlideResult, error) {
	result, err := ReadDeck(ctx, ref, ReadOptions{})
	if err != nil {
		return nil, err
	}

	slide := result.Deck.FindSlide(slideID)
	if slide == nil {
		return nil, fmt.Errorf("%w: %s", ErrSlideNotFound, slideID)
	}

	output, err := encodeSlide(slide, f)
	if err != nil {
		return nil, err
	}

	return &GetSlideResult{
		Slide:  slide,
		Output: output,
	}, nil
}

// GetSlideFromPath is a convenience function that detects the backend.
func GetSlideFromPath(ctx context.Context, path, slideID string, f format.Format) (*GetSlideResult, error) {
	backendName := DetectBackend(path)
	ref := model.Ref{
		Backend: backendName,
		Path:    path,
	}
	return GetSlide(ctx, ref, slideID, f)
}

// UpdateSlideOptions configures the UpdateSlide operation.
type UpdateSlideOptions struct {
	Confirm bool
}

// UpdateSlideResult contains the result of an UpdateSlide operation.
type UpdateSlideResult struct {
	Updated bool
	Message string
}

// UpdateSlide updates a single slide.
func UpdateSlide(ctx context.Context, ref model.Ref, slideID string, updates *model.Slide, opts UpdateSlideOptions) (*UpdateSlideResult, error) {
	if !opts.Confirm {
		return &UpdateSlideResult{
			Updated: false,
			Message: "Set confirm=true to apply changes",
		}, ErrConfirmRequired
	}

	// Read current deck
	readResult, err := ReadDeck(ctx, ref, ReadOptions{})
	if err != nil {
		return nil, err
	}

	// Find and update the slide
	currentSlide := readResult.Deck.FindSlide(slideID)
	if currentSlide == nil {
		return nil, fmt.Errorf("%w: %s", ErrSlideNotFound, slideID)
	}

	// Create a diff for the update
	diff := model.NewDiff(readResult.Deck.ID)

	if updates.Title != "" && updates.Title != currentSlide.Title {
		diff.AddChange(model.NewUpdateChange(
			fmt.Sprintf("slides/%s/title", slideID),
			currentSlide.Title, updates.Title))
		currentSlide.Title = updates.Title
	}

	if updates.Subtitle != "" && updates.Subtitle != currentSlide.Subtitle {
		diff.AddChange(model.NewUpdateChange(
			fmt.Sprintf("slides/%s/subtitle", slideID),
			currentSlide.Subtitle, updates.Subtitle))
		currentSlide.Subtitle = updates.Subtitle
	}

	if len(updates.Body) > 0 {
		diff.AddChange(model.NewUpdateChange(
			fmt.Sprintf("slides/%s/body", slideID),
			currentSlide.Body, updates.Body))
		currentSlide.Body = updates.Body
	}

	if len(updates.Notes) > 0 {
		diff.AddChange(model.NewUpdateChange(
			fmt.Sprintf("slides/%s/notes", slideID),
			currentSlide.Notes, updates.Notes))
		currentSlide.Notes = updates.Notes
	}

	if diff.IsEmpty() {
		return &UpdateSlideResult{
			Updated: false,
			Message: "No changes to apply",
		}, nil
	}

	// Apply the changes
	backend, err := DefaultRegistry.Get(ref.Backend)
	if err != nil {
		return nil, err
	}

	if err := backend.Apply(ctx, ref, diff); err != nil {
		return nil, err
	}

	return &UpdateSlideResult{
		Updated: true,
		Message: fmt.Sprintf("Slide %s updated successfully", slideID),
	}, nil
}

// encodeSlideList serializes a slide list to the requested format.
func encodeSlideList(slides []SlideInfo, f format.Format) (string, error) {
	switch f {
	case format.FormatJSON:
		data, err := json.MarshalIndent(slides, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case format.FormatTOON, "":
		return encodeSlidesToTOON(slides), nil
	default:
		return encodeSlidesToTOON(slides), nil
	}
}

// encodeSlidesToTOON encodes a slide list to TOON format.
func encodeSlidesToTOON(slides []SlideInfo) string {
	var result string
	for _, s := range slides {
		result += fmt.Sprintf("slide %s %s", s.ID, s.Layout)
		if s.Title != "" {
			result += fmt.Sprintf(" %s", s.Title)
		}
		result += "\n"
	}
	return result
}

// encodeSlide serializes a single slide to the requested format.
func encodeSlide(slide *model.Slide, f format.Format) (string, error) {
	switch f {
	case format.FormatJSON:
		data, err := json.MarshalIndent(slide, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case format.FormatTOON, "":
		encoder := format.NewTOONEncoder()
		return encoder.EncodeSlide(slide), nil
	default:
		encoder := format.NewTOONEncoder()
		return encoder.EncodeSlide(slide), nil
	}
}
