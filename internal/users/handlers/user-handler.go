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

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req, _ := common.DecodeAndValidate[application.CreateUserReq](r)

	id, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = common.Encode(w, http.StatusCreated, map[string]int{"id": id})
	if err != nil {
		return
	}
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	req, err := common.DecodeAndValidate[application.LoginUserReq](r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.service.LoginUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = common.Encode(w, http.StatusOK, res)
	if err != nil {
		return
	}
}

func (h *UserHandler) privateHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := common.GetLoggedInUser(r)
	err := common.Encode(w, http.StatusOK, claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoadUserRoutes(b *sqlx.DB, r chi.Router, key string) {
	userRepo := db.NewUserRepo(b)
	jwt := common.NewJWTGenerator(key)
	userSvc := application.NewUserService(userRepo, jwt)
	userHandler := NewUserHandler(userSvc)

	// Public routes group
	r.Post("/register", userHandler.RegisterUser)
	r.Post("/login", userHandler.LoginUser)

	//Protected routes group
	r.Group(func(r chi.Router) {
		r.Use(common.GetAuthMiddlewareFunc(jwt))
		r.Get("/private", userHandler.privateHandler)
	})
}
