package models

// =============================================================================
// File
// =============================================================================

type File struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}