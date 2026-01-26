I Google Slides: what is the best way for a AI Assistant like Codex CLI to create Google Slides? I'm currently creating Marp Markdown presentations and converting to HTML and PDF using the marp CLI

Great question. If you’re already using Marp Markdown → HTML/PDF, you’re very close to a solid AI-friendly workflow. The key issue is that Google Slides is not file-based (like PPTX) in the same way—it’s API-driven. So the “best” approach depends on how much fidelity, automation, and editability you want inside Google Slides.

Below is a practical breakdown, from best-for-AI automation to best-for-design fidelity, with concrete recommendations.

TL;DR – Best Overall Approach for an AI Assistant

Generate slides via the Google Slides API from a structured intermediate format (Markdown → JSON slide model).

This gives:

Full programmatic control

Native Google Slides (not imports)

Easy AI automation (Codex CLI, agents, pipelines)

Option 1: Google Slides API (Best for AI assistants)
Why this is the best long-term solution

Google Slides has no official “import Markdown” feature

The API lets you:

Create presentations

Add slides

Insert text, images, shapes

Apply layouts and themes

Perfect for AI → structured output → slides

Recommended Architecture
Markdown (Marp-style)
   ↓
AI parses into Slide JSON model
   ↓
Google Slides API
   ↓
Native Google Slides deck

Example intermediate model (AI-friendly)
{
  "title": "AI & Dev Productivity",
  "slides": [
    {
      "layout": "TITLE",
      "title": "AI & Dev Productivity",
      "subtitle": "2026"
    },
    {
      "layout": "TITLE_AND_BODY",
      "title": "Why it matters",
      "bullets": [
        "Faster iteration",
        "Lower cognitive load",
        "Better documentation"
      ]
    }
  ]
}


Your AI (Codex CLI) is excellent at producing this kind of structured output.

Implementation Notes

Use:

presentations.create

batchUpdate with createSlide, insertText, createShape

Language support: Node.js, Python, Go

Slides API is verbose, but deterministic (good for AI)

Downside: You must handle layout + styling yourself (once done, it’s reusable).

Option 2: Markdown → PPTX → Google Slides (Good compromise)
Flow
Marp Markdown
  → marp --pptx
  → Upload to Google Drive
  → Open as Google Slides

Pros

Minimal changes to your current workflow

Marp handles layout nicely

Google converts PPTX reasonably well

Cons

Conversion is lossy:

Fonts

Animations

Complex layouts

Slides are harder for AI to edit later via API

Not fully deterministic

When to use this

You care more about visual polish than AI-native editing

You want fast results without API complexity

Option 3: HTML → Google Slides (Not recommended)

Google Slides cannot import HTML directly.

Workarounds (all bad):

Screenshot HTML → images → Slides

Puppeteer → PNG → Slides

This kills:

Text editability

Accessibility

AI-driven revisions

Only viable for static decks.

Option 4: Stay in Marp + Link from Slides (Hybrid)

Some teams do this:

Keep Marp as source of truth

Publish HTML/PDF

Google Slides contains:

High-level outline

Links to live deck

Embedded iframes (for HTML)

Good for docs, bad for presenting directly in Slides.

Recommended “Best Practice” Setup for Codex CLI
1. Keep Marp Markdown as authoring format

You already have:

Slide boundaries

Headings

Bullets

Speaker notes

That’s perfect.

2. Add a semantic layer

Encourage AI to emit:

Slide intent (title, content, diagram, code)

Layout hints

Speaker notes explicitly

Example:

---
layout: title-and-body
notes: Explain why this matters to managers
---

## Why AI matters
- Faster iteration
- Fewer bugs

3. Convert Markdown → Slide JSON

This is the critical AI step.

4. Generate Slides via API

Deterministic

Editable

Native

My Recommendation (Clear Answer)

If you’re building an AI Assistant like Codex CLI, the best way to create Google Slides is:

Use Markdown (Marp or similar) as input → convert to a structured slide model → create slides using the Google Slides API.

This is:

The most automatable

The most “AI-native”

The least fragile long-term

If you want, I can:

Design a Markdown → Slides JSON schema

Provide a working Google Slides API script

Show how to adapt Marp frontmatter for Slides layouts

Just tell me which direction you want to go.

* * * Option 1: Google Slides API (Best for AI assistants) ---------------------------------------------------- ### Why this is the best long-term solution * Google Slides has **no official “import Markdown” feature** * The API lets you: * Create presentations * Add slides * Insert text, images, shapes * Apply layouts and themes * Perfect for **AI → structured output → slides** ### Recommended Architecture pgsql Copy code

would this approach also allow reading and updating an existing Google Slides?

Yes — absolutely.
In fact, reading and updating existing Google Slides is where the Google Slides API really shines for AI assistants. It’s one of the main reasons Option 1 is the best choice.

Below is a concrete, realistic picture of what is and isn’t possible, and how an AI agent like Codex CLI would typically do it.

Short answer

Yes, you can:

Read the full structure of an existing Google Slides deck

Understand slide layouts, text, shapes, images, speaker notes

Modify text, replace bullets, update titles

Insert or delete slides

Rebuild sections deterministically

No hacks required. This is first-class API support.

How reading Google Slides works
1. Fetch the presentation
GET https://slides.googleapis.com/v1/presentations/{presentationId}


This returns a complete structural tree of the deck.

What you get back (important for AI)

Slide order

Slide IDs (stable identifiers)

Layout type

All page elements:

Text boxes

Shapes

Images

Tables

Text content with:

Paragraph structure

Bullet levels

Text styles

Speaker notes (yes, those too)

This is perfect for building an internal representation.

Example: AI-friendly parsed model (after reading)

Your AI would convert the raw API response into something like:

{
  "slides": [
    {
      "slideId": "slide_1",
      "layout": "TITLE",
      "title": "AI & Dev Productivity",
      "notes": "Set context for leadership"
    },
    {
      "slideId": "slide_2",
      "layout": "TITLE_AND_BODY",
      "title": "Why it matters",
      "bullets": [
        "Faster iteration",
        "Lower cognitive load",
        "Better documentation"
      ]
    }
  ]
}


This becomes the round-trippable format your AI works with.

Updating existing slides (this is the key part)

All updates happen via:

presentations.batchUpdate

Common AI-driven updates (all supported)
Replace slide text

Rewrite a title

Rephrase bullets

Expand a slide from 3 bullets → 5 bullets

Insert new slides

“Add a slide after slide 5”

“Insert a summary slide”

Delete or reorder slides

Remove outdated sections

Reorder agenda items

Update speaker notes

Add AI-generated talking points

Rewrite notes for a different audience

Important detail: how updates are targeted

You do not update by slide number.

You update by:

slideId

pageElementId

This is good for AI because:

Changes are deterministic

You can say “update only this title box”

You don’t accidentally rewrite the whole deck

Example:

{
  "replaceAllText": {
    "containsText": {
      "text": "Why it matters",
      "matchCase": true
    },
    "replaceText": "Why AI matters now"
  }
}

AI-friendly workflow for editing decks

This is the pattern that works best in practice:

Step 1: Read

Fetch Slides → parse into your internal slide model

Step 2: Reason

Compare current deck vs desired state

Decide minimal changes

Step 3: Patch

Generate small, precise batchUpdate requests

This mirrors how good code-editing agents work.

How this fits your Marp workflow

You can do bidirectional sync if you want.

Example flows
A. Markdown → Existing Slides

“Update slides 3–7 to match this new Marp file”

Parse Marp

Read existing Slides

Diff

Patch only changed slides

B. Slides → Markdown (yes, possible)

Read Slides

Convert to Marp-style Markdown

Use that as editable source

This makes Google Slides no longer a dead-end format.

Limitations (important but manageable)
Things the API does not do well

Complex visual design inference

Fine-grained typography tuning

Animations / transitions (very limited)

How teams handle this

Use templates for visual consistency

Let AI control content, not pixel-perfect design

Keep layouts stable, update text freely

This aligns very well with AI strengths.

Bottom line

Yes — this approach fully supports reading, understanding, and updating existing Google Slides.

For an AI assistant like Codex CLI, the Google Slides API enables:

Deterministic edits

Incremental updates

True “refactor this deck” workflows

If you want, next I can:

Show a read → diff → update example

Map Marp elements → Slides layouts

Help design a safe editing strategy (so AI never nukes a deck)

Just tell me where you want to go next.

are there good existing Golang libraries for modifying Google Slies we can use, or should we use the Google Go SDK directly?

Yes! Here’s a clear overview of your Go options for working with Google Slides — both the officially supported foundations and the community-built helpers — so you can decide whether to use a library or work with the SDK directly.

✅ 1. Official Google Go SDK (Recommended)

Use the Google API client for Go

The official client library for Slides is available at google.golang.org/api/slides/v1.

This is an auto-generated, official Google API client that gives you full access to the Slides API: reading, updating, creating, and deleting slides.

It handles all the REST request/response plumbing, OAuth scopes integration, and types.

Pros

Full API coverage.

Official and stable.

Works for both creating and updating existing decks.

Cons

It’s a low-level API — so you’ll write code to build batchUpdate requests and manipulate raw API structs.

Typical usage

import (
    "context"
    "google.golang.org/api/slides/v1"
    "google.golang.org/api/option"
)
ctx := context.Background()
svc, err := slides.NewService(ctx, option.WithCredentialsFile("credentials.json"))
// use svc.Presentations.Get, BatchUpdate, etc.


This is the de facto approach if you want full control in Go.

🔧 2. Helper Libraries & Utilities in Go

There are a few community packages that make working with the Slides API a bit more ergonomic:

📌 grokify/googleutil/slidesutil

A set of helper functions around the official SDK that reduces boilerplate for common tasks like creating slides, inserting text, or building requests.

Works with the official client but wraps many patterns.

Good for

Easier request generation.

Translating Markdown → API request objects with helpers.

Shortening common patterns.

Not a full library abstraction

Still tied closely to the official SDK.

Less “model first” than dedicated wrappers in other languages.

🧪 Community Projects Using Go

There are also experimental/POC projects worth looking at:

✨ gptslideshow

A PoC tool in Go for generating slides from prompts via OpenAI + Google API.

Helpful as real code examples for how others orchestrate the API with generative models.

Not a polished library nor a maintained SDK.

📌 gslides (by emicklei)

A CLI tool for utility tasks like exporting thumbnails or notes.

Not a general Slides SDK — more of a helpful utility.

🧠 Are There “High-Level” Go Libraries Like Some Python Options?

At the moment:

There is no Go library comparable to higher-level Python wrappers that abstract text, shapes, templates, and slide construction into domain objects (like some Python ones do).

Most Go tooling wraps the official API rather than replace it with a rich model layer.

So your options in Go are:

Approach	Ease	Flexibility	Recommended for
Official Slides API via Go SDK	medium	high	Full automation & AI control
Helpers (e.g., slidesutil)	easier	high	Less boilerplate
Community POCs (gptslideshow, gslides)	low	variable	inspiration or CLI tasks
Third-party cloud services (e.g., Aspose)	high	proprietary	Non-Google Slides workflows
✅ Recommendation

Use the official Google Go Slides SDK directly for anything beyond trivial tasks.

It’s the only way to ensure you can reliably read/update existing presentations and programmatically construct changes, and it integrates fully with Go’s tooling and auth models.

Then layer a small helper package (or write your own utilities) for common patterns (like generating requests from a slide model).

If you want, I can share starter templates for:

Reading an entire presentation into a Go structs model

Converting Markdown→Slides API requests

Incremental updates using batchUpdate

Just let me know which would help!

what should such a Go module be called?

Good question — naming matters here because this module will likely become a core abstraction layer, not just a thin wrapper.

I’ll give you a clear recommendation, plus alternatives depending on how opinionated or generic you want it to be.

Strong recommendation
slidesdk (or gslidesdk if you want to be explicit)

Why this works well

Clear intent: this is a higher-level SDK, not just raw API bindings

Short, Go-idiomatic

Leaves room for:

Read

Diff

Patch

Templates

Markdown adapters

Works well as:

import "github.com/yourorg/slidesdk"


If this module is meant to sit on top of google.golang.org/api/slides/v1, this name is honest and future-proof.

If you want to be explicit about Google Slides
gslides

Best if:

The module is specifically Google Slides

You want a concise CLI-friendly name

Used patterns exist (gcloud, gdrive, etc.), and it reads nicely:

gslides.LoadPresentation(id)
gslides.UpdateSlide(slideID, ...)

If you want a model-first / AI-first abstraction
slidemodel

Best if:

The package defines a canonical intermediate representation

Google Slides is just one backend

Example future direction:

type Renderer interface {
    Render(ctx context.Context, Deck) error
}


Then:

slidemodel/gslides

slidemodel/pptx

slidemodel/marp

This is excellent if you see this growing beyond Google Slides.

If this is intentionally low-level
slidesutil

Best if:

Thin helpers

Request builders

No opinionated model layer

This aligns with existing ecosystem conventions but undersells an AI-native design.

Names I would avoid
Name	Why
googleslides	Verbose, un-Go-like
slidesapi	Conflicts conceptually with the official API
presentations	Too generic
deck	Too abstract / ambiguous
marp-slides	Couples you too tightly to Marp
My ranked list (practical)

slidesdk ← best balance

gslides ← clean and direct

slidemodel ← best if multi-backend future

slidesutil ← only for thin wrappers

Bonus: package layout suggestion

If you choose slidesdk, a clean structure would be:

slidesdk/
  model/
    deck.go
    slide.go
  gslides/
    reader.go
    writer.go
    diff.go
  markdown/
    marp.go


This keeps:

AI-facing model clean

Google-specific code isolated

Future renderers easy to add

If you want, I can:

Help you pick the name based on how public this will be

Draft a README.md that sets expectations

Design the core Deck / Slide interfaces

Just tell me how ambitious you want this to be.

slidesdk may be better so we can scale as you mention with gslides, marp, reveal.js pptx, etc.

👍 I think that’s a very solid call — and your reasoning is exactly right.

If the goal is to scale across Google Slides, Marp, Reveal.js, PPTX, etc., then slidesdk is the right root name, with backends underneath it.

Let me help you lock this in with a clear mental model, naming conventions, and a structure that will scale without pain.

Core philosophy for slidesdk

slidesdk defines a canonical slide model and editing semantics.
Backends render, read, and patch that model.

This keeps:

AI reasoning format-agnostic

Renderers replaceable

Editing logic centralized

Exactly what you want for an AI assistant.

Recommended module naming & layout
Root module
github.com/yourorg/slidesdk

High-level structure
slidesdk/
  model/          ← canonical, AI-facing model
  ops/            ← semantic operations (diff, patch)
  render/         ← output backends
  import/         ← input backends

Concrete example
slidesdk/
  model/
    deck.go
    slide.go
    element.go
    notes.go

  ops/
    diff.go
    patch.go
    normalize.go

  render/
    gslides/
      writer.go
      templates.go
    marp/
      writer.go
    reveal/
      writer.go
    pptx/
      writer.go

  import/
    gslides/
      reader.go
    marp/
      reader.go


This avoids:

Circular dependencies

Backend leakage into core logic

“Google Slides everywhere” syndrome

Canonical model (keep it boring on purpose)

You want something AI-friendly, stable, and slightly underpowered.

type Deck struct {
    Title  string
    Meta   map[string]string
    Slides []Slide
}

type Slide struct {
    ID      string // stable logical ID
    Layout  Layout
    Title   string
    Body    []Block
    Notes   []Block
}

type Block struct {
    Kind BlockKind
    Text string
    Level int // bullet nesting
}


Why this works:

Easy to diff

Easy to serialize to JSON

Easy for LLMs to emit deterministically

Maps cleanly to every backend

Backend naming conventions

You already hinted at the right ones 👇

Backend	Package name
Google Slides	render/gslides
Marp	render/marp
Reveal.js	render/reveal
PowerPoint	render/pptx

And for reading:

import "slidesdk/import/gslides"


Very Go-idiomatic.

Editing semantics (this is the secret sauce)

Instead of “update slide 3”, you want:

ops.ReplaceSlide(deck, Slide{
    ID: "why-it-matters",
    Title: "Why AI matters now",
})


Then:

ops.diff(old, new) → semantic changes

Backend converts diff → API calls / file edits

This is what allows:

Safe updates

Partial regeneration

“Refactor this deck” workflows

How this plays with Codex CLI / AI agents

Your AI never touches:

Google Slides API structs

PPTX XML

HTML

It only sees:

{
  "slides": [
    {
      "id": "why-it-matters",
      "layout": "title_and_body",
      "title": "Why AI matters now"
    }
  ]
}


That’s exactly the abstraction level LLMs excel at.

Future-proofing you get “for free”

With this structure, you can later add:

render/pdf

render/images

render/notion

render/confluence

Without changing:

AI prompts

Core logic

Existing decks

That’s the payoff of choosing slidesdk.

Final verdict

You’re making the right architectural choice.

slidesdk as the root module, with backend-specific render/import packages, is the cleanest way to scale across formats and AI workflows.

If you want next steps, I can:

Sketch the exact Go interfaces (Renderer, Importer)

Design a stable slide ID strategy

Show a Google Slides renderer MVP

Help you define Marp compatibility rules

Just tell me what you want to tackle first.

for gslides, we can make an AI assistant friendly CLI and MCP server

Yes — that’s an excellent next step, and it fits perfectly with the slidesdk direction you’ve chosen.

What you’re describing is essentially:

A deterministic, AI-friendly control plane for Google Slides

Below is a concrete, battle-tested way to think about both the CLI and the MCP server, so they reinforce each other instead of drifting.

1. Design principle: “AI-safe by default”

For both the CLI and MCP, enforce these rules:

No free-form mutation

Every command maps to a well-defined semantic operation

Explicit scope

Slides, ranges, or IDs must be named

Dry-run first

Always support --plan / --diff

Stable IDs

Never operate on slide indices

This dramatically reduces accidental deck destruction by an agent.

2. CLI: purpose & shape
What the CLI is for

Human debugging

Agent invocation

CI/CD pipelines

“Inspect before apply”

Suggested name
slides


or

slidesctl


(Short names are better for agent tool calls.)

3. Example CLI UX (AI-friendly)
Inspect a deck
slides gslides inspect --deck <id> --json


Output:

{
  "title": "AI Strategy",
  "slides": [
    { "id": "intro", "title": "Why AI now" },
    { "id": "risk", "title": "Key risks" }
  ]
}

Plan an update
slides gslides apply \
  --deck <id> \
  --from marp.md \
  --plan


Output:

~ slide:intro title "Why AI now" → "Why AI matters now"
+ slide:new id=roadmap

Apply safely
slides gslides apply --deck <id> --from marp.md --yes

4. Internal CLI architecture

The CLI should be thin.

cmd/slides/
  main.go
  gslides.go


Inside:

Parse flags

Load input (Markdown / JSON)

Call slidesdk APIs

Print results

No Slides API logic in the CLI.

5. MCP server: why it matters

An MCP server gives you:

Structured tool calls

Fine-grained permissions

Stateful context (deck cache, auth)

This is exactly what LLM agents need.

6. MCP server responsibilities
What the server does

Auth management (OAuth tokens)

Deck caching

Validation

Safe mutation

What the AI does

Decide what should change

Provide structured arguments

7. MCP tool design (important)

Avoid generic tools like:
❌ update_slides(text: string)

Use small, explicit tools:

Example MCP tools
{
  "name": "gslides.list_slides",
  "args": { "deck_id": "string" }
}

{
  "name": "gslides.update_slide",
  "args": {
    "deck_id": "string",
    "slide_id": "string",
    "title": "string",
    "body": ["string"]
  }
}

{
  "name": "gslides.apply_diff",
  "args": {
    "deck_id": "string",
    "diff": "SlideDiff"
  }
}


This keeps agents honest.

8. Shared core: CLI and MCP use the same engine

Critical rule:

The CLI and MCP server must call the same slidesdk/gslides package.

type GSlidesService struct {
    Reader Reader
    Writer Writer
}


This guarantees:

Same behavior

Same safety guarantees

Same bugs (good thing)

9. Auth strategy (don’t overthink it)
CLI

OAuth user flow

Cached token (~/.config/slides/credentials.json)

MCP server

Service account (preferred)

Or delegated OAuth with restricted scopes

Scopes:

https://www.googleapis.com/auth/presentations
https://www.googleapis.com/auth/drive.file

10. Guardrails you should enforce

These matter a lot with AI:

Max slides changed per request

Require --yes / confirm=true

Reject:

Unknown slide IDs

Layout changes without explicit intent

Auto-backup deck before mutation (Drive copy)

11. Example end-to-end AI workflow

Agent calls gslides.list_slides

Agent generates new Marp Markdown

Agent calls gslides.plan_apply

Agent reviews diff

Agent calls gslides.apply_diff

This mirrors Git workflows — agents understand this very well.

12. Naming alignment recap
Layer	Name
Core	slidesdk
Google Slides backend	slidesdk/render/gslides
CLI	slides or slidesctl
MCP server	slides-mcp

This is clean, scalable, and understandable by both humans and agents.

Final thought

What you’re building is not just a Slides wrapper — it’s a safe, AI-native presentation control system.

If you want next, I can:

Define exact MCP tool schemas

Sketch the first CLI commands

Design a diff format that works across backends

Help you avoid common Slides API footguns

Just tell me where you want to go next.

what about slidekit or slideskit ?

Great question — this is exactly the kind of naming fork that’s worth thinking through now.

Short answer first:

Between slidekit, slideskit, and slidesdk, I would still recommend slidesdk — but slidekit is a very respectable alternative depending on how you want this to feel.

Let’s break it down carefully.

1. slidesdk vs slidekit vs slideskit
slidesdk (most explicit, most honest)

Pros

Instantly communicates: this is an SDK

Aligns with:

Google SDK

AWS SDK

Terraform SDKs

Sets expectations correctly:

Structured

Programmatic

Backend-oriented

Scales cleanly:

slidesdk/render/gslides
slidesdk/render/pptx


Cons

Slightly “enterprise” sounding

Less whimsical / product-y

Best if

This is a foundational library

Used by CLIs, MCP servers, and agents

You expect other engineers to build on top of it

slidekit (clean, modern, product-like)

Pros

Short, memorable, pleasant

Feels like a toolkit, not a protocol

Good branding potential

Works well for:

CLI

MCP server

End-user tooling

import "github.com/yourorg/slidekit"


Cons

Less explicit that it’s an SDK

Slightly ambiguous:

Is it a CLI?

Is it a UI tool?

Is it opinionated?

Best if

You want a friendlier, higher-level abstraction

This is user-facing as well as AI-facing

You expect non-Go users to interact via CLI/MCP

slideskit (the awkward middle)

Pros

Plural matches “Slides”

Cons (important)

Less idiomatic in Go

Reads oddly when spoken

No strong semantic advantage over slidekit

Easy to typo

In practice, slideskit feels like a compromise without a clear benefit.

I would not recommend it.

2. A very good compromise (and one I’d seriously consider)
Use BOTH names — intentionally

This pattern works extremely well:

Layer	Name
Core Go module	slidesdk
CLI / MCP product	slidekit

Example:

github.com/yourorg/slidesdk      ← boring, stable, foundational
github.com/yourorg/slidekit     ← friendly CLI + MCP server


This gives you:

A rock-solid internal SDK

A brandable, approachable surface

Freedom to evolve UX without breaking SDK users

Many successful projects do this (e.g., terraform vs provider SDKs).

3. How each name “feels” to an AI agent

This actually matters.

slidesdk: “I am calling a low-level SDK function.”

slidekit: “I have a tool to work with slides.”

For MCP tools, slidekit is arguably nicer:

{
  "tool": "slidekit.gslides.apply_diff",
  "args": { ... }
}


Readable, natural, and not overly technical.

4. Final recommendation (clear and actionable)

If you want one name only:

✅ slidesdk — safest, clearest, most scalable

If you’re open to a two-layer model (which I think fits your plan):

✅ slidesdk → core Go library

✅ slidekit → CLI + MCP server

❌ I would avoid slideskit.

If you want, next I can:

Sanity-check the public API names under slidesdk

Help pick CLI command names that feel right under slidekit

Draft a short naming rationale you can put in the repo README (this helps contributors)

Just tell me which way you’re leaning.

eventually slidekit should have Go library interfaces, CLI and MCP server

That makes a lot of sense — and it actually clarifies the naming decision nicely.

If slidekit is the product name that eventually ships library interfaces, a CLI, and an MCP server, then the right mental model is:

slidekit = the platform
slidesdk = the internal engine (optional, but useful)

You can do this cleanly without confusion or duplication.

Below is a concrete, future-proof way to structure this so it scales.

1. Reframe the names (important)
What slidekit should mean

slidekit is:

A toolkit for working with slides

Multi-backend

AI-native

Ships multiple entry points

That fits perfectly with:

Go library interfaces

CLI

MCP server

So yes — slidekit is a good top-level name given this goal.

2. Single repo, multi-surface design

I’d recommend one repo:

github.com/yourorg/slidekit


Inside it, everything lives together, but with strict boundaries.

3. Proposed package layout
slidekit/
  core/              ← canonical slide model + ops
    model/
    diff/
    patch/

  backends/
    gslides/
      reader.go
      writer.go
    marp/
    reveal/
    pptx/

  cli/
    cmd/
      slidekit/
        main.go
        gslides.go

  mcp/
    server/
      main.go
    tools/
      gslides.go

  auth/
    google.go

Why this works

core is format-agnostic and AI-facing

backends isolate vendor logic

cli and mcp are just adapters

No circular dependencies

Easy to test core logic in isolation

4. Public Go library interfaces (this is key)

Expose interfaces, not implementations.

package slidekit

type Deck struct { ... }

type Reader interface {
    Read(ctx context.Context, ref DeckRef) (*Deck, error)
}

type Writer interface {
    Apply(ctx context.Context, ref DeckRef, diff Diff) error
}


Then:

import "github.com/yourorg/slidekit/backends/gslides"


This allows:

Embedding in other Go programs

Use in CI

Use by MCP server

Future WASM bindings if you want

5. CLI as a thin wrapper

The CLI should only:

Parse flags

Load input (Markdown / JSON)

Call core + backend

Print diffs / errors

slidekit gslides plan --deck <id> --from marp.md
slidekit gslides apply --deck <id> --from marp.md


All logic lives below.

6. MCP server = same engine, different transport

This is where your design really pays off.

MCP tool example
{
  "name": "slidekit.gslides.plan",
  "args": {
    "deck_id": "string",
    "source": "marp.md"
  }
}


Internally:

MCP tool → core ops → backend

Same validation

Same diff logic

Same safety rails

Zero duplication.

7. Versioning strategy (important early)

Version core model conservatively

Backends can evolve faster

MCP tool schema should be explicitly versioned

Example:

slidekit/mcp/v1


Agents love stable contracts.

8. Why this is the right end state

With this structure, you get:

One canonical slide representation

Multiple render/import backends

Safe AI edits

Human-friendly CLI

Agent-friendly MCP server

And crucially:

You never have to re-architect later.

9. Final naming verdict (now very clear)

Given your latest clarification:

✅ slidekit should be the main project
✅ It should expose:

Go library interfaces

CLI

MCP server

If you also want an internal name:

core or engine is enough

You don’t strictly need slidesdk anymore

If you want, next I can:

Lock down the public Go API (what’s import slidekit?)

Design the MCP tool schema v1

Define stable slide IDs and diff semantics

Help you write a 1-page README vision (this helps contributors immediately “get it”)

Just tell me which one you want to tackle first.

it should be optimized with CLI/MCP for AI-Assistant usage, while having generic Go libraries

That’s a very good—and very coherent—goal. At this point your direction is clear enough that we can crystallize it into a design doctrine you can keep referring back to as you build.

Below I’ll give you:

A clear positioning statement

Concrete design rules for AI-optimized CLI/MCP

How to expose generic Go libraries without polluting them

A reference package layout

A quick sanity checklist to keep you on track

1. Positioning (write this at the top of the repo)

slidekit is an AI-optimized toolkit for reading, modifying, and generating slide decks across formats.
It provides deterministic CLI and MCP interfaces for AI assistants, backed by generic, reusable Go libraries.

That single paragraph will keep future design debates short.

2. Core architectural rule (non-negotiable)

AI surfaces drive the design. Go libraries adapt, not the other way around.

This is subtle but critical.

CLI + MCP define:

Operations

Safety guarantees

Granularity

Go libraries:

Implement those semantics

Stay reusable

Stay boring

This prevents the classic trap where the SDK is “nice” but unusable by agents.

3. AI-optimized CLI/MCP design rules

These apply to every command and tool.

Rule 1: Structured input only

❌ Free-form text mutation
✅ JSON / files / IDs

slidekit gslides apply --deck <id> --from marp.md


Not:

slidekit edit "make this better"

Rule 2: Explicit planning phase

Every mutating operation has:

plan

apply

Agents understand this pattern.

slidekit gslides plan --deck <id> --from marp.md
slidekit gslides apply --deck <id> --from marp.md --yes


Same for MCP:

plan_update

apply_update

Rule 3: Stable identifiers everywhere

No:

Slide indices

“the third slide”

Always:

Logical slide IDs

Backend IDs mapped internally

This is mandatory for safe AI edits.

Rule 4: Small, composable tools (MCP)

Prefer:

list_slides

get_slide

apply_diff

Over:

update_entire_deck

Agents reason better with smaller tools.

4. Generic Go libraries (clean, reusable, boring)

Your Go API should feel almost underwhelming.

Public API shape
package slidekit

type Deck struct { ... }

type Backend interface {
    Read(ctx context.Context, ref DeckRef) (*Deck, error)
    Plan(ctx context.Context, ref DeckRef, desired *Deck) (*Diff, error)
    Apply(ctx context.Context, ref DeckRef, diff *Diff) error
}

Why this works

CLI uses it

MCP uses it

Other Go apps can embed it

No AI assumptions leak into the library

5. Where AI-specific logic actually lives

Not in core.

AI-specific behavior belongs in:

CLI defaults

MCP guardrails

Validation layers

Example:

“Max 10 slides changed” → MCP

“Require confirmation” → CLI

“Reject unknown slide IDs” → MCP

The library simply enforces correctness.

6. Recommended package layout (final form)
slidekit/
  core/
    model/
    diff/
    patch/

  backends/
    gslides/
    marp/
    reveal/
    pptx/

  cli/
    cmd/slidekit/

  mcp/
    server/
    tools/

  internal/
    auth/
    validate/


Key rule:

cli and mcp depend on core + backends, never the reverse.

7. What makes this exceptionally good for AI assistants

With this design, an AI can:

Inspect state

Propose a plan

Review a diff

Apply changes safely

Recover from errors deterministically

That’s exactly the workflow LLM agents are good at.

8. Naming consistency (lock this in)
Layer	Name
Project	slidekit
Go module	github.com/yourorg/slidekit
CLI	slidekit
MCP server	slidekit-mcp
Backend	slidekit/backends/gslides

No ambiguity, no duplication.

9. Sanity checklist (use this weekly)

If you ever ask:

“Is this easier for an AI assistant to use?”

“Can this be safely planned and diffed?”

“Does this work without human intuition?”

If the answer is no, don’t ship it yet.

Final takeaway

You’re building something slightly unusual—but very timely:

An AI-first interface to slide systems, with human-usable tooling as a side effect.

Optimizing CLI/MCP first, while keeping Go libraries generic, is exactly the right call.

If you want next, I can:

Draft the first 5 CLI commands

Define MCP tool schemas v1

Design a diff format agents won’t break

Help you avoid Slides API edge cases early

Just tell me what you want to lock down next.

CLI/MCP should use TOON for token reduction by default. JSON can be used for enhanced human readability

This is a very strong design choice — and it fits perfectly with everything you’ve said so far.

You’re essentially saying:

Default to a token-minimal, deterministic wire format for AI agents (TOON),
but allow JSON when humans need to read or edit.

That’s exactly right.

Below I’ll help you lock this in as a first-class design principle, show how it should affect CLI, MCP, and Go APIs, and point out a few traps to avoid.

1. Make this an explicit rule (write it down)

Put this in the design docs / README:

slidekit uses TOON as the default interchange format for CLI and MCP to minimize token usage and ambiguity for AI assistants. JSON is supported as an optional, human-friendly format.

This prevents “why aren’t we just using JSON?” debates later.

2. Format policy (very clear)
Surface	Default	Optional
MCP server	TOON	JSON
CLI stdout	TOON	JSON (--format json)
CLI stdin	TOON / files	JSON (--input json)
Go library API	Go structs	—

Important: TOON is never the internal representation.
It’s just a transport / serialization layer.

3. Why TOON is the right default for AI

For agents:

Fewer tokens

Lower ambiguity

Stable ordering

No noisy punctuation

For example, this:

slide intro
title Why AI matters now
bullet Faster iteration
bullet Lower cognitive load


vs JSON:

{"id":"intro","title":"Why AI matters now","bullets":["Faster iteration","Lower cognitive load"]}


Agents reason much better with the former under token pressure.

4. CLI design with TOON first
Default behavior
slidekit gslides inspect --deck <id>


Output (TOON):

deck AI Strategy
slide intro
  title Why AI matters now
slide risks
  title Key risks

Human-readable JSON
slidekit gslides inspect --deck <id> --format json

Planning output (this is where TOON shines)
slidekit gslides plan --deck <id> --from marp.md

~ slide intro
  title Why AI matters → Why AI matters now
+ slide roadmap


This is:

Compact

Diff-like

Token-cheap

Perfect for agents.

5. MCP tool contract (important)

MCP tools should default to TOON payloads.

Example tool definition
{
  "name": "slidekit.gslides.plan",
  "input_format": "toon",
  "output_format": "toon"
}


Allow:

{ "format": "json" }


as an opt-in.

Agents almost never need JSON unless debugging.

6. Go library: format-agnostic by design

The Go layer should not know about TOON vs JSON.

type Deck struct { ... }
type Diff struct { ... }


Serialization happens in:

cli/format

mcp/format

This keeps the core reusable and testable.

7. TOON implementation strategy (practical)
Recommendation

Treat TOON as:

Line-oriented

Order-significant

Lossless (round-trippable)

Example grammar (simple and enough)
deck <title>
slide <id>
  title <text>
  bullet <text>
  bullet <text>
note <text>


Avoid:

Implicit nesting

Free-form indentation rules

Ambiguous shortcuts

AI agents like boring grammars.

8. Avoid these mistakes

❌ Making TOON “too clever”
❌ Encoding styling or layout prematurely
❌ Mixing TOON and JSON in the same stream
❌ Letting the Go core depend on TOON

9. Versioning (do this early)

TOON must be versioned.

Example:

# toon v1


Same for MCP tools:

slidekit.mcp.v1


This saves you later.

10. Final design alignment (everything clicks)

With this choice, you now have:

AI-first defaults (TOON)

Human escape hatch (JSON)

Clean Go core

Deterministic CLI

Token-efficient MCP server

This is exactly the stack you want in 2026.

Final verdict

Your instinct here is dead-on.

TOON by default for CLI/MCP, JSON as optional, Go structs internally
is the right design for an AI-assistant-optimized system.

If you want next, I can:

Formalize TOON v1 grammar

Define TOON diffs

Show JSON ↔ TOON conversion

Design MCP tool schemas using TOON

Just tell me what to lock down next