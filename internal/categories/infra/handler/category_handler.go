package handler

import (
	"encoding/json"
	"net/http"
	"yadwy-backend/internal/categories/usercase"
)

type CreateCategoryHandler struct {
	createUseCase *usercase.CreateCategoryUseCase
}

func NewCategoryHandler(createUseCase *usercase.CreateCategoryUseCase) *CreateCategoryHandler {
	return &CreateCategoryHandler{
		createUseCase: createUseCase,
	}
}

type createCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *CreateCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "there is noo cat", http.StatusBadRequest)
		return
	}

	category, err := h.createUseCase.Execute(req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(category)
	if err != nil {
		return
	}
}
