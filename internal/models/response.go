package models

// MetaInfo represents metadata for responses (e.g., pagination)
type MetaInfo struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

type PaginatedResponse struct {
	Items []interface{} `json:"items"`
	Meta  MetaInfo      `json:"meta"`
}
