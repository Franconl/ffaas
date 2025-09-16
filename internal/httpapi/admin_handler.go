package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Franconl/ffaas/internal/core"
	"github.com/Franconl/ffaas/internal/repo"
	"github.com/go-chi/chi/v5"
)

// AdminHandler agrupa los handlers http para la administracion de feature flags
// Depende de un repositorio que implemente la interface repo.Flags
type AdminHandler struct {
	repo repo.Flags
}

// NewAdminHandler crea un nuevo admin handler usando el repositorio pasado por parametro
func NewAdminHandler(r repo.Flags) *AdminHandler {
	return &AdminHandler{repo: r}
}

// Create maneja POST /flags
func (h *AdminHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateFlagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, repo.ErrInvalidBody.Error())
		return
	}

	if req.Key == "" {
		writeError(w, http.StatusBadRequest, repo.ErrKeyRequired.Error())
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
		if errors.Is(err, repo.ErrKeyAlreadyUsed) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
	}

	// se mapea la flag al DTO de respuesta

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

// List maneja GET /flags
func (h *AdminHandler) List(w http.ResponseWriter, r *http.Request) {
	val, err := h.repo.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	flags := make([]FlagResponse, 0, len(val))

	for _, f := range val {
		flags = append(flags, FlagResponse{
			ID:          f.ID,
			Key:         f.Key,
			Description: f.Description,
			Enabled:     f.Enabled,
			Percentage:  f.Percentage,
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		})
	}

	resp := ListFlagsResponse{Items: flags}

	writeJSON(w, http.StatusOK, resp)
}

// GetByID maneja GET /flags/{id}
func (h *AdminHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	resp := FlagResponse{
		ID:          val.ID,
		Key:         val.Key,
		Description: val.Description,
		Enabled:     val.Enabled,
		Percentage:  val.Percentage,
		CreatedAt:   val.CreatedAt,
		UpdatedAt:   val.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetByKey maneja GET /flags/key/{key}
func (h *AdminHandler) GetByKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	val, err := h.repo.GetByKey(key)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	resp := FlagResponse{
		ID:          val.ID,
		Key:         val.Key,
		Description: val.Description,
		Enabled:     val.Enabled,
		Percentage:  val.Percentage,
		CreatedAt:   val.CreatedAt,
		UpdatedAt:   val.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, resp)
}

// DeleteByID maneja DELETE /flags/{id}
func (h *AdminHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := h.repo.GetByID(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := h.repo.DeleteByID(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusNoContent, "")
}

func (h *AdminHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateFlagRequest

	id := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	flag, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	flag.Description = req.Description
	flag.Enabled = req.Enabled
	flag.Percentage = req.Percentage

	if errUpdate := h.repo.Update(flag); errUpdate != nil {
		writeError(w, http.StatusInternalServerError, errUpdate.Error())
		return
	}

	writeJSON(w, http.StatusOK, FlagResponse{
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
