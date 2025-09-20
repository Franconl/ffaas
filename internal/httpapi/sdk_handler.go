package httpapi

import (
	"net/http"

	"github.com/Franconl/ffaas/internal/repo"
)

type SdkHandler struct {
	repo repo.Flags
}

// NewSdkHandler constructor, recibe un repo que cumpla la interfaz Repo.Flags
func NewSdkHandler(repo repo.Flags) *SdkHandler {
	return &SdkHandler{
		repo: repo,
	}
}

func (h *SdkHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
	}

	flags := make([]FlagResponse, 0, len(list))

	for _, val := range list {
		flags = append(flags, FlagResponse{
			Key:        val.Key,
			ID:         val.ID,
			Percentage: val.Percentage,
			Enabled:    val.Enabled,
			CreatedAt:  val.CreatedAt,
			UpdatedAt:  val.UpdatedAt,
		})
	}

	resp := ListFlagsResponse{Items: flags}

	writeJSON(w, http.StatusOK, resp)
}

func (h *SdkHandler) Eval(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	userID := r.URL.Query().Get("userId")

	if key == "" {
		writeError(w, http.StatusBadRequest, "key is required")
		return
	}
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId is required")
		return
	}

	f, err := h.repo.GetByKey(key)
	if err != nil {
		writeError(w, http.StatusNotFound, "flag not found")
		return
	}

	enabled := f.Eval(userID)

	resp := struct {
		Key     string `json:"key"`
		UserID  string `json:"user_id"`
		Enabled bool   `json:"enabled"`
	}{
		Key:     f.Key,
		UserID:  userID,
		Enabled: enabled,
	}

	writeJSON(w, http.StatusOK, resp)
}
