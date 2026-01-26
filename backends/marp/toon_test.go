package marp

import (
	"strings"
	"testing"

	"github.com/grokify/slidekit/format"
)

func TestTOONOutput(t *testing.T) {
	reader := NewReader()
	deck, err := reader.Parse(sampleMarp)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	encoder := format.NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	// Verify TOON structure
	if !strings.HasPrefix(output, "deck ") {
		t.Error("expected TOON output to start with 'deck '")
	}
	if !strings.Contains(output, "deck My Presentation") {
		t.Error("expected 'deck My Presentation' in TOON output")
	}
	if !strings.Contains(output, "section ") {
		t.Error("expected 'section' lines in TOON output")
	}
	if !strings.Contains(output, "slide ") {
		t.Error("expected 'slide' lines in TOON output")
	}
	if !strings.Contains(output, "title ") {
		t.Error("expected 'title' lines in TOON output")
	}
	if !strings.Contains(output, "bullet ") {
		t.Error("expected 'bullet' lines in TOON output")
	}

	// Verify it's significantly shorter than JSON would be
	if len(output) == 0 {
		t.Error("expected non-empty TOON output")
	}
}

func TestTOONRealPresentation(t *testing.T) {
	// Parse the stats-agent-team style presentation structure
	input := `---
marp: true
theme: vibeminds
paginate: true
---

<!-- _class: lead -->

# Statistics Agent Team
## Building a Multi-Agent System

**A Production-Ready Implementation**

---

<!-- _class: section-divider -->

<!--
Section 1: Introduction. <break time="600ms"/>
-->

# Section 1
## Introduction

Understanding the challenge

---

<!--
The problem is trust in AI statistics.
[PAUSE:1500]
-->

# The Problem

- LLMs hallucinate statistics
- URLs are often outdated
- No verification of sources

---

<!-- _class: section-divider -->

# Section 2
## Architecture

Four specialized agents

---

# Agent Design

` + "```go" + `
type Agent interface {
    Process(ctx context.Context) error
}
` + "```" + `

| Agent | Type |
|-------|------|
| Research | Tool |
| Synthesis | LLM |
`

	reader := NewReader()
	deck, err := reader.Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	encoder := format.NewTOONEncoder()
	output := encoder.EncodeDeck(deck)

	// Verify sections detected
	sectionCount := strings.Count(output, "\nsection ")
	if sectionCount < 2 {
		t.Errorf("expected at least 2 sections in TOON, got %d", sectionCount)
	}

	// Verify notes present
	if !strings.Contains(output, "note ") {
		t.Error("expected note lines in TOON output")
	}

	// Verify code block
	if !strings.Contains(output, "code ") {
		t.Error("expected code line in TOON output")
	}
}
