package reporter

import "time"

// Question describes a single possible question
type Question struct {
	ID           string `json:"uniqueIdentifier,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	QuestionType *int   `json:"questionType,omitempty"`
	Placeholder  string `json:"placeholderString,omitempty"`
}

// Day contains all snapshots, possible questions (schema version 2 only) and metadata about a specific day
// Reporter writes one JSON file per day
type Day struct {
	Snapshots     []Snapshot `json:"snapshots,omitempty"`
	Questions     []Question `json:"questions,omitempty"`
	Date          time.Time  `json:"-,omitempty"` // Only filled when data wasn't loaded from string
	FileInfo      File       `json:"-,omitempty"` // Only filled when data wasn't loaded from string
	SchemaVersion int        `json:"-"`
}

// GetEarliestSnapshot returns the first snapshot for a given day
func (d *Day) GetEarliestSnapshot() Snapshot {
	return d.Snapshots[len(d.Snapshots)]
}

// GetLatestSnapshot returns the latest snapshot for a given day
func (d *Day) GetLatestSnapshot() Snapshot {
	return d.Snapshots[len(d.Snapshots)-1]
}
