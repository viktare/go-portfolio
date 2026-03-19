package models

// =============================================================================
// Company
// =============================================================================

type Company struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Website string `json:"website,omitempty"`
	Logo    string `json:"logo,omitempty"`
}

type CreateCompany struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
}

type UpdateCompany struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
}