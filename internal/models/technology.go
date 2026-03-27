package models

// =============================================================================
// TechnologyField
// =============================================================================

type TechnologyField struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Code string `json:"code"`
}

// =============================================================================
// Technology
// =============================================================================

type Technology struct {
    ID    string          `json:"id"`
    Name  string          `json:"name"`
    Code  string          `json:"code"`
    Field TechnologyField `json:"field"`
}

type CreateTechnology struct {
    Name    string `json:"name"`
    Code    string `json:"code"`
    FieldID string `json:"fieldId"`
}

type UpdateTechnology struct {
    Name    string `json:"name"`
    FieldID string `json:"fieldId"`
}