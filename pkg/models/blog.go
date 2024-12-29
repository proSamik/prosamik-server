package models

type Blog struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	ViewsCount  int    `json:"views_count"`
}

type BlogResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
