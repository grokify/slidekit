package model

import "time"

// Audio represents audio attachment for LMS/video generation.
type Audio struct {
	Source   AudioSource   `json:"source"`
	Path     string        `json:"path,omitempty"`     // Local file path
	URL      string        `json:"url,omitempty"`      // Remote URL
	Script   string        `json:"script,omitempty"`   // Text for TTS generation
	Duration time.Duration `json:"duration,omitempty"` // Explicit or computed duration
	Voice    string        `json:"voice,omitempty"`    // TTS voice ID
}

// AudioSource identifies audio origin.
type AudioSource string

const (
	AudioSourceFile  AudioSource = "file"  // Pre-recorded audio file
	AudioSourceURL   AudioSource = "url"   // Remote audio URL
	AudioSourceTTS   AudioSource = "tts"   // Generate via TTS from script
	AudioSourceNotes AudioSource = "notes" // Generate via TTS from speaker notes
)

// AudioSources returns all valid audio source values.
func AudioSources() []AudioSource {
	return []AudioSource{
		AudioSourceFile,
		AudioSourceURL,
		AudioSourceTTS,
		AudioSourceNotes,
	}
}

// IsValid returns true if the audio source is a recognized value.
func (s AudioSource) IsValid() bool {
	switch s {
	case AudioSourceFile, AudioSourceURL, AudioSourceTTS, AudioSourceNotes:
		return true
	}
	return false
}

// NeedsTTS returns true if this audio source requires TTS generation.
func (s AudioSource) NeedsTTS() bool {
	return s == AudioSourceTTS || s == AudioSourceNotes
}

// NewFileAudio creates an audio reference from a local file.
func NewFileAudio(path string, duration time.Duration) *Audio {
	return &Audio{
		Source:   AudioSourceFile,
		Path:     path,
		Duration: duration,
	}
}

// NewURLAudio creates an audio reference from a remote URL.
func NewURLAudio(url string, duration time.Duration) *Audio {
	return &Audio{
		Source:   AudioSourceURL,
		URL:      url,
		Duration: duration,
	}
}

// NewTTSAudio creates an audio reference for TTS generation from a script.
func NewTTSAudio(script, voice string) *Audio {
	return &Audio{
		Source: AudioSourceTTS,
		Script: script,
		Voice:  voice,
	}
}

// NewNotesAudio creates an audio reference for TTS generation from speaker notes.
func NewNotesAudio(voice string) *Audio {
	return &Audio{
		Source: AudioSourceNotes,
		Voice:  voice,
	}
}

// HasContent returns true if the audio has a valid source reference.
func (a *Audio) HasContent() bool {
	if a == nil {
		return false
	}
	switch a.Source {
	case AudioSourceFile:
		return a.Path != ""
	case AudioSourceURL:
		return a.URL != ""
	case AudioSourceTTS:
		return a.Script != ""
	case AudioSourceNotes:
		return true // Will be derived from slide notes
	}
	return false
}
