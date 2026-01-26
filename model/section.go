package model

import "time"

// Section groups related slides (maps to LMS sections/chapters).
type Section struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Slides []Slide `json:"slides"`
	Audio  *Audio  `json:"audio,omitempty"` // Section-level audio for LMS
}

// SlideCount returns the number of slides in this section.
func (s *Section) SlideCount() int {
	return len(s.Slides)
}

// FindSlide finds a slide by ID within this section.
func (s *Section) FindSlide(id string) *Slide {
	for i := range s.Slides {
		if s.Slides[i].ID == id {
			return &s.Slides[i]
		}
	}
	return nil
}

// TotalDuration returns the total audio duration of the section.
// If section-level audio is set, returns its duration.
// Otherwise, sums slide-level audio durations.
func (s *Section) TotalDuration() time.Duration {
	if s.Audio != nil && s.Audio.Duration > 0 {
		return s.Audio.Duration
	}
	var total time.Duration
	for _, slide := range s.Slides {
		if slide.Audio != nil && slide.Audio.Duration > 0 {
			total += slide.Audio.Duration
		}
	}
	return total
}

// HasAudio returns true if the section has audio (either section-level or slide-level).
func (s *Section) HasAudio() bool {
	if s.Audio != nil {
		return true
	}
	for _, slide := range s.Slides {
		if slide.Audio != nil {
			return true
		}
	}
	return false
}
