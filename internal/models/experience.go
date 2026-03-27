package models

import "time"

type Experience struct {
    ID           string       `json:"id"`
    Role         string       `json:"role"`
    Description  string       `json:"description"`
    Company      Company      `json:"company"`
    StartDate    time.Time    `json:"startDate"`
    EndDate      *time.Time   `json:"endDate,omitempty"`
    Technologies []Technology `json:"technologies"`
    Order        int          `json:"order"`
}

func (e Experience) IsCurrent() bool {
    return e.EndDate == nil
}

type CreateExperience struct {
    Role          string   `json:"role"`
    Description   string   `json:"description"`
    CompanyID     string   `json:"companyId"`
    StartDate     string   `json:"startDate"`
    EndDate       string   `json:"endDate"`
    TechnologyIDs []string `json:"technologyIds"`
    Order         int      `json:"order"`
}

type UpdateExperience struct {
    Role          string   `json:"role"`
    Description   string   `json:"description"`
    StartDate     string   `json:"startDate"`
    EndDate       string   `json:"endDate"`
    TechnologyIDs []string `json:"technologyIds"`
    Order         int      `json:"order"`
}