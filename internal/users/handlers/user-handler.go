package handlers

import (
	"errors"
	"net/http"
	"yadwy-backend/internal/common"
	"yadwy-backend/internal/users/application"
	"yadwy-backend/internal/users/db"
	"yadwy-backend/internal/users/domain/modles"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// @Summary Register a new user
// @Description Register a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body application.CreateUserReq true "User registration details"
// @Success 201 {object} application.CreateUserRes
// @Failure 400 {object} common.ErrorResponse
// @Failure 409 {object} common.ErrorResponse
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req, _ := common.DecodeAndValidate[application.CreateUserReq](r)

	res, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err = common.Encode(w, http.StatusCreated, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body application.LoginUserReq true "User login credentials"
// @Success 200 {object} application.LoginUserRes
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	req, err := common.DecodeAndValidate[application.LoginUserReq](r)
	if err != nil {
		handleError(w, err)
		return
	}

	res, err := h.service.LoginUser(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	if err = common.Encode(w, http.StatusOK, res); err != nil {
		handleError(w, err)
		return
	}
}

// @Summary Get private user information
// @Description Get authenticated user's information
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} common.UserClaims
// @Failure 401 {object} common.ErrorResponse
// @Router /users/private [get]
func (h *UserHandler) privateHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := common.GetLoggedInUser(r)
	err := common.Encode(w, http.StatusOK, claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoadUserRoutes(b *sqlx.DB, r chi.Router, jwt *common.JWTGenerator) {
	userRepo := db.NewUserRepo(b)
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

func handleError(w http.ResponseWriter, err error) {
	var appErr *common.Error
	if errors.As(err, &appErr) {
		// Handle custom application errors
		switch appErr.Code() {
		case modles.UserNotFoundError:
			common.SendError(w, http.StatusNotFound, string(appErr.Code()), appErr.Error())
		case modles.EmailAlreadyExistsError:
			common.SendError(w, http.StatusConflict, string(appErr.Code()), appErr.Error())
		case modles.InvalidUserCredentialsError:
			common.SendError(w, http.StatusUnauthorized, string(appErr.Code()), appErr.Error())
		case modles.InvalidUserRoleError:
			common.SendError(w, http.StatusBadRequest, string(appErr.Code()), appErr.Error())
		default:
			common.SendError(w, http.StatusInternalServerError, string(appErr.Code()), appErr.Error())
		}
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		common.SendError(w, http.StatusBadRequest, "validation_error", common.FormatValidationError(err))
		return
	}

	// Handle unknown errors
	common.SendError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
}
