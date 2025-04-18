package infra

import (
	"errors"
	"net/http"
	"yadwy-backend/internal/category/application"
	"yadwy-backend/internal/category/domain"
	"yadwy-backend/internal/common"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Handler struct {
	s      *application.CategoryService
	logger *zap.Logger
}

func NewCategoryHandler(service *application.CategoryService, logger *zap.Logger) *Handler {
	return &Handler{
		s:      service,
		logger: logger,
	}
}

// @Summary Create a new category
// @Description Create a new category with image upload
// @Tags categories
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param image formData file true "Category image"
// @Success 201 "Category created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid input"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 403 {object} common.ErrorResponse "Forbidden - Admin only"
// @Router /category [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(128 * 1024 * 1024); err != nil {
		handleError(w, common.NewErrorf(domain.FailedToParseImage, "%v", err))
		return
	}

	_, file, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("Failed to parse form file", zap.Error(err))
		handleError(w, common.NewErrorf(domain.FailedToParseImage, "%v", err))
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	err = h.s.CreateCategory(name, description, file)

	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Get all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Success 200 {array} application.CategoryRes
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Router /category [get]
func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.s.GetAllCategories(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	catRes := make([]application.CategoryRes, len(cats))
	for i, cat := range cats {
		catRes[i] = application.MapToCategoryRes(cat)
	}

	if err = common.Encode(w, http.StatusOK, catRes); err != nil {
		handleError(w, err)
		return
	}
}

func LoadCategoryRoutes(b *sqlx.DB, logger *zap.Logger, jwt *common.JWTGenerator) http.Handler {
	ar := chi.NewRouter()
	cr := NewCategoryRepo(b, logger)
	files, _ := common.NewFileService("/home/nerd/images", "http://localhost:3000/images")
	cs := application.NewCategoryService(cr, files, logger)
	ch := NewCategoryHandler(cs, logger)

	ar.Use(common.GetAuthMiddlewareFunc(jwt))
	ar.Get("/", ch.GetAllCategories)

	// Admin routes
	ar.With(common.GetAdminMiddlewareFun(jwt)).Post("/", ch.Create)
	return ar
}

func handleError(w http.ResponseWriter, err error) {
	var appErr *common.Error
	if errors.As(err, &appErr) {
		switch appErr.Code() {
		case domain.FailedToParseImage:
			common.SendError(w, http.StatusInternalServerError, string(appErr.Code()), appErr.Error())
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

	common.SendError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
}
