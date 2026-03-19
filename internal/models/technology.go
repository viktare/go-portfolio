package models

// =============================================================================
// TechnologyField
// =============================================================================

type TechnologyField struct {
	ID   string `json:"id"`
	Name string `json:"name"` // "Backend"
	Code string `json:"code"` // "backend"
}

// =============================================================================
// Technology
// =============================================================================

type Technology struct {
	ID    int             `json:"id"`
	Name  string          `json:"name"`  // "Php - Laravel"
	Code  string          `json:"code"`  // "laravel" (unique)
	Field TechnologyField `json:"field"`
}

type CreateTechnology struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	FieldID string `json:"fieldId"`
}

type UpdateTechnology struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	FieldID string `json:"fieldId"`
}