package category

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type CreateCategoryHandler struct {
	cs *CategoryService
}

func NewCategoryHandler(createUseCase *CategoryService) *CreateCategoryHandler {
	return &CreateCategoryHandler{
		cs: createUseCase,
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

	category, err := h.cs.Execute(req.Name, req.Description)
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

func LoadCategoryRoutes(b *sqlx.DB) http.Handler {
	ar := chi.NewRouter()
	ar.Use(AdminOnly)
	cr := NewCategoryRepo(b)
	cs := NewCategoryService(cr)
	ch := NewCategoryHandler(cs)
	ar.Post("/", ch.Create)

	return ar
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("AdminOnly")

		next.ServeHTTP(w, r)
	})
}
