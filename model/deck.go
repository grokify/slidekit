package model

import "time"

// Deck represents a complete presentation.
type Deck struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Meta     Meta      `json:"meta"`
	Sections []Section `json:"sections"`
	Theme    *Theme    `json:"theme,omitempty"`
}

// Meta contains presentation metadata.
type Meta struct {
	Author      string            `json:"author,omitempty"`
	Date        string            `json:"date,omitempty"`
	Description string            `json:"description,omitempty"`
	Keywords    []string          `json:"keywords,omitempty"`
	Custom      map[string]string `json:"custom,omitempty"`
}

// SlideCount returns the total number of slides across all sections.
func (d *Deck) SlideCount() int {
	count := 0
	for _, s := range d.Sections {
		count += len(s.Slides)
	}
	return count
}

// AllSlides returns a flat list of all slides in order.
func (d *Deck) AllSlides() []Slide {
	var slides []Slide
	for _, s := range d.Sections {
		slides = append(slides, s.Slides...)
	}
	return slides
}

// FindSlide finds a slide by ID across all sections.
func (d *Deck) FindSlide(id string) *Slide {
	for i := range d.Sections {
		for j := range d.Sections[i].Slides {
			if d.Sections[i].Slides[j].ID == id {
				return &d.Sections[i].Slides[j]
			}
		}
	}
	return nil
}

// FindSection finds a section by ID.
func (d *Deck) FindSection(id string) *Section {
	for i := range d.Sections {
		if d.Sections[i].ID == id {
			return &d.Sections[i]
		}
	}
	return nil
}

// TotalDuration returns the total audio duration of the deck.
func (d *Deck) TotalDuration() time.Duration {
	var total time.Duration
	for _, s := range d.Sections {
		total += s.TotalDuration()
	}
	return total
}
