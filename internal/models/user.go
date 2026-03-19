package models

import "time"

// =============================================================================
// User
// =============================================================================

type User struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Surname     string        `json:"surname"`
	Bio         string        `json:"bio"`
	BirthDate   *time.Time    `json:"birthDate,omitempty"`
	Location    string        `json:"location"`
	Contacts    []UserContact `json:"contacts"`
	Experiences []Experience  `json:"experiences"`
	Projects    []Project     `json:"projects"`
	ResumeFile  *File         `json:"resumeFile,omitempty"`
}

// Age computes age from BirthDate so it never goes stale
func (u User) Age() int {
	if u.BirthDate == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - u.BirthDate.Year()
	if now.YearDay() < u.BirthDate.YearDay() {
		age--
	}
	return age
}

// UpdateUser — no ID (comes from auth/session), no nested relations
type UpdateUser struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birthDate"` // "2000-01-15" — parse server-side
	Location  string `json:"location"`
}

// =============================================================================
// UserContact
// =============================================================================

type UserContact struct {
	ID   string `json:"id"`
	Name string `json:"name"` // "LinkedIn"
	Code string `json:"code"` // "linkedin"
	Link string `json:"link"` // "https://linkedin.com/in/..."
}

type CreateUserContact struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Link string `json:"link"`
}

type UpdateUserContact struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Link string `json:"link"`
}