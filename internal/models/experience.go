package models

import "time"

// =============================================================================
// Experience
// =============================================================================

type Experience struct {
	ID           int          `json:"id"`
	Role         string       `json:"role"`
	Description  string       `json:"description"`
	Company      Company      `json:"company"`
	StartDate    time.Time    `json:"startDate"`
	EndDate      *time.Time   `json:"endDate,omitempty"` // nil = current position
	Technologies []Technology `json:"technologies"`
	Order        int          `json:"order"`
}

// IsCurrent returns true if this is an ongoing experience
func (e Experience) IsCurrent() bool {
	return e.EndDate == nil
}

// CreateExperience — references company by ID, technologies by ID list
type CreateExperience struct {
	Role          string `json:"role"`
	Description   string `json:"description"`
	CompanyID     string `json:"companyId"`
	StartDate     string `json:"startDate"`     // "2024-06-04" — parse server-side
	EndDate       string `json:"endDate"`        // "" or "2025-09-16"
	TechnologyIDs []int  `json:"technologyIds"`
	Order         int    `json:"order"`
}

type UpdateExperience struct {
	Role          string `json:"role"`
	Description   string `json:"description"`
	CompanyID     string `json:"companyId"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	TechnologyIDs []int  `json:"technologyIds"`
	Order         int    `json:"order"`
}