package util

import "time"

type User struct {
	ID          string     `json:"id" db:"id"`
	Login       string     `json:"login" db:"login"`
	PassWord    string     `json:"password" db:"password"`
	Name        string     `json:"name" db:"name"`
	DateOfBirth *time.Time `json:"date_of_birth" db:"date_of_birth"`
}

type Session struct {
	ID     string `db:"id"`
	UserId string `db:"user_id"`
}

type Document struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Text      string `json:"text" db:"text"`
	OwnerID   string `json:"owner_id" db:"owner_id"`
	SectionID string `json:"section_id" db:"section_id"`
}

type Section struct {
	ID      string `db:"id"`
	Path    string `db:"path"`
	Name    string `db:"name"`
	OwnerID string `db:"owner_id"`
}

type SectionTree struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	OwnerID     string         `json:"owner_id"`
	SubSections []*SectionTree `json:"subsections"`
	Documents   []*Document    `json:"documents"`
}

func (s *Section) ToSectionTree() *SectionTree {
	return &SectionTree{
		ID:          s.ID,
		Name:        s.Name,
		OwnerID:     s.OwnerID,
		SubSections: make([]*SectionTree, 0),
		Documents:   make([]*Document, 0),
	}
}
