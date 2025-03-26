package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"yadwy-backend/internal/banner/service"
)

// BannerHandler handles HTTP requests for banners
type BannerHandler struct {
	service *service.BannerService
}

// NewBannerHandler creates a new banner handler
func NewBannerHandler(service *service.BannerService) *BannerHandler {
	return &BannerHandler{
		service: service,
	}
}

// CreateBannerRequest represents the request body for creating a banner
type CreateBannerRequest struct {
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	TargetURL string `json:"target_url"`
	Position  string `json:"position"`
	IsActive  bool   `json:"is_active"`
}

// BannerResponse represents the response body for a banner
type BannerResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	TargetURL string `json:"target_url"`
	IsActive  bool   `json:"is_active"`
	Position  string `json:"position"`
}

// CreateBanner handles the creation of a new banner
func (h *BannerHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	var req CreateBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateBanner(req.Title, req.ImageURL, req.TargetURL, req.Position, req.IsActive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// GetBanner handles retrieving a banner by ID
func (h *BannerHandler) GetBanner(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	banner, err := h.service.GetBanner(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if banner == nil {
		http.Error(w, "Banner not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banner)
}

// ListBanners handles listing all banners
func (h *BannerHandler) ListBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.service.ListBanners()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// ListActiveBanners handles listing all active banners
func (h *BannerHandler) ListActiveBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.service.ListActiveBanners()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

// UpdateBanner handles updating a banner
func (h *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req CreateBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateBanner(id, req.Title, req.ImageURL, req.TargetURL, req.Position, req.IsActive)
	if err != nil {
		if err.Error() == "banner not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteBanner handles deleting a banner
func (h *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteBanner(id)
	if err != nil {
		if err.Error() == "banner not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
