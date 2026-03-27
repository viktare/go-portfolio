package models

// =============================================================================
// Project
// =============================================================================


type Project struct {
    ID           string       `json:"id"`
    Name         string       `json:"name"`
    Code         string       `json:"code"`
    Status       string       `json:"status"`
    Description  *string      `json:"description"`
    Link         *string      `json:"link"`
    Technologies []Technology `json:"technologies"`
}

type CreateProject struct {
    Name          string   `json:"name"`
    Code          string   `json:"code"`
    Status        string   `json:"status"`
    Description   *string  `json:"description"`
    Link          *string  `json:"link"`
    TechnologyIDs []string `json:"technologyIds"`
}

type UpdateProject struct {
    Name          string   `json:"name"`
    Status        string   `json:"status"`
    Description   *string  `json:"description"`
    Link          *string  `json:"link"`
    TechnologyIDs []string `json:"technologyIds"`
}