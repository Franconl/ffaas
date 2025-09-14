package httpapi

import "time"

// --- Requests ---

type CreateFlagRequest struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Percentage  int    `json:"percentage"`
}

type UpdateFlagRequest struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Percentage  int    `json:"percentage"`
}

// --- Responses ---

type FlagResponse struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Percentage  int       `json:"percentage"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Para listas (ej: GET /api/flags)
type ListFlagsResponse struct {
	Items []FlagResponse `json:"items"`
}

// Para SDK /sdk/eval
type EvalResponse struct {
	Key     string `json:"key"`
	UserID  string `json:"user_id"`
	Enabled bool   `json:"enabled"`
}
