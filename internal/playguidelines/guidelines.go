package playguidelines

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed data/guidelines.json
var guidelinesJSON []byte

type Guideline struct {
	Section         string      `json:"section"`
	Title           string      `json:"title"`
	Content         string      `json:"content"`
	Verification    string      `json:"verification"` // automated, hybrid, manual
	Sources         []string    `json:"sources,omitempty"`
	AutomatedChecks []string    `json:"automated_checks,omitempty"`
	ManualChecks    []string    `json:"manual_checks,omitempty"`
	Subsections     []Guideline `json:"subsections,omitempty"`
}

type DB struct {
	Guidelines []Guideline `json:"guidelines"`
	index      map[string]*Guideline
}

func Load() (*DB, error) {
	var db DB
	if err := json.Unmarshal(guidelinesJSON, &db); err != nil {
		return nil, err
	}
	db.buildIndex()
	return &db, nil
}

func (db *DB) buildIndex() {
	db.index = make(map[string]*Guideline)
	var walk func(gs []Guideline)
	walk = func(gs []Guideline) {
		for i := range gs {
			db.index[gs[i].Section] = &gs[i]
			walk(gs[i].Subsections)
		}
	}
	walk(db.Guidelines)
}

func (db *DB) Get(section string) (*Guideline, bool) {
	g, ok := db.index[section]
	return g, ok
}

func (db *DB) TopLevel() []Guideline {
	return db.Guidelines
}

func (db *DB) Flatten() []Guideline {
	var out []Guideline
	var walk func(gs []Guideline)
	walk = func(gs []Guideline) {
		for _, g := range gs {
			out = append(out, g)
			walk(g.Subsections)
		}
	}
	walk(db.Guidelines)
	return out
}

func (db *DB) Search(query string) []Guideline {
	query = strings.ToLower(query)
	var out []Guideline
	for _, g := range db.Flatten() {
		if strings.Contains(strings.ToLower(g.Section), query) ||
			strings.Contains(strings.ToLower(g.Title), query) ||
			strings.Contains(strings.ToLower(g.Content), query) {
			out = append(out, g)
		}
	}
	return out
}
