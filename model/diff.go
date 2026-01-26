package model

// Diff represents changes between two deck states.
type Diff struct {
	DeckID  string   `json:"deck_id"`
	Changes []Change `json:"changes"`
}

// Change represents a single modification.
type Change struct {
	Op        ChangeOp `json:"op"`
	Path      string   `json:"path"`                 // e.g., "sections/0/slides/2/title"
	SlideID   string   `json:"slide_id,omitempty"`   // Target slide ID
	SectionID string   `json:"section_id,omitempty"` // Target section ID
	OldValue  any      `json:"old_value,omitempty"`
	NewValue  any      `json:"new_value,omitempty"`
}

// ChangeOp identifies the type of change.
type ChangeOp string

const (
	ChangeAdd    ChangeOp = "add"
	ChangeRemove ChangeOp = "remove"
	ChangeUpdate ChangeOp = "update"
	ChangeMove   ChangeOp = "move"
)

// ChangeOps returns all valid change operation values.
func ChangeOps() []ChangeOp {
	return []ChangeOp{
		ChangeAdd,
		ChangeRemove,
		ChangeUpdate,
		ChangeMove,
	}
}

// IsValid returns true if the change op is a recognized value.
func (op ChangeOp) IsValid() bool {
	switch op {
	case ChangeAdd, ChangeRemove, ChangeUpdate, ChangeMove:
		return true
	}
	return false
}

// IsEmpty returns true if the diff has no changes.
func (d *Diff) IsEmpty() bool {
	return len(d.Changes) == 0
}

// ChangeCount returns the number of changes.
func (d *Diff) ChangeCount() int {
	return len(d.Changes)
}

// CountByOp returns the count of changes by operation type.
func (d *Diff) CountByOp() map[ChangeOp]int {
	counts := make(map[ChangeOp]int)
	for _, c := range d.Changes {
		counts[c.Op]++
	}
	return counts
}

// AddChanges returns only add operations.
func (d *Diff) AddChanges() []Change {
	return d.filterByOp(ChangeAdd)
}

// RemoveChanges returns only remove operations.
func (d *Diff) RemoveChanges() []Change {
	return d.filterByOp(ChangeRemove)
}

// UpdateChanges returns only update operations.
func (d *Diff) UpdateChanges() []Change {
	return d.filterByOp(ChangeUpdate)
}

// MoveChanges returns only move operations.
func (d *Diff) MoveChanges() []Change {
	return d.filterByOp(ChangeMove)
}

func (d *Diff) filterByOp(op ChangeOp) []Change {
	var result []Change
	for _, c := range d.Changes {
		if c.Op == op {
			result = append(result, c)
		}
	}
	return result
}

// NewDiff creates a new empty diff for the given deck ID.
func NewDiff(deckID string) *Diff {
	return &Diff{
		DeckID:  deckID,
		Changes: []Change{},
	}
}

// AddChange adds a change to the diff.
func (d *Diff) AddChange(c Change) {
	d.Changes = append(d.Changes, c)
}

// NewAddChange creates an add operation.
func NewAddChange(path string, value any) Change {
	return Change{
		Op:       ChangeAdd,
		Path:     path,
		NewValue: value,
	}
}

// NewRemoveChange creates a remove operation.
func NewRemoveChange(path string, value any) Change {
	return Change{
		Op:       ChangeRemove,
		Path:     path,
		OldValue: value,
	}
}

// NewUpdateChange creates an update operation.
func NewUpdateChange(path string, oldValue, newValue any) Change {
	return Change{
		Op:       ChangeUpdate,
		Path:     path,
		OldValue: oldValue,
		NewValue: newValue,
	}
}

// NewMoveChange creates a move operation.
func NewMoveChange(fromPath, toPath string) Change {
	return Change{
		Op:       ChangeMove,
		Path:     fromPath,
		NewValue: toPath,
	}
}
