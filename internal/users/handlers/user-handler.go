package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
	"yadwy-backend/internal/common"
	"yadwy-backend/internal/users/application"
	"yadwy-backend/internal/users/db"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req, _ := common.Decode[CreateUserRequest](r)

	id, err := h.service.CreateUser(req.Name, req.Email, req.Password, "ADMIN")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = common.Encode(w, http.StatusCreated, map[string]int{"id": id})
	if err != nil {
		return
	}
}

func LoadUserRoutes(b *sqlx.DB, r chi.Router) {
	userRepo := db.NewUserRepo(b)
	userSvc := application.NewUserService(userRepo)
	userHandler := NewUserHandler(userSvc)

	r.Post("/", userHandler.RegisterUser)
}
