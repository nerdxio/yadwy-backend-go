package api

import (
	"encoding/json"
	"net/http"

	"yadwy-backend/internal/category/service"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	service *service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

type createCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Create handles the creation of a new category
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	category, err := h.service.CreateCategory(req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}
