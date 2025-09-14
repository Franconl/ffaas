package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Franconl/ffaas/internal/core"
	"github.com/Franconl/ffaas/internal/repo"
)

type AdminHandler struct {
	repo repo.Flags
}

func NewAdminHandler(r repo.Flags) *AdminHandler {
	return &AdminHandler{repo: r}
}

func (h AdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Key == "" {
		writeError(w, http.StatusBadRequest, "key is necesary")
		return
	}

	flag := core.FeatureFlag{
		Key:         req.Key,
		Description: req.Description,
		Percentage:  req.Percentage,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := h.repo.Create(&flag)
	if err != nil {
		if err.Error() == "flag key already exist" {
			writeError(w, http.StatusConflict, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
	}

	writeJSON(w, http.StatusCreated, FlagResponse{
		ID:          flag.ID,
		Key:         flag.Key,
		Description: flag.Description,
		Enabled:     flag.Enabled,
		Percentage:  flag.Percentage,
		CreatedAt:   flag.CreatedAt,
		UpdatedAt:   flag.UpdatedAt,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
