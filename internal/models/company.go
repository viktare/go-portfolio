package models

// =============================================================================
// Company
// =============================================================================

type Company struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Website *string `json:"website"`
	Logo    *string `json:"logo"`
}

type CreateCompanyDTO struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
}

type UpdateCompanyDTO struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
}