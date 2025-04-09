package infra

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"yadwy-backend/internal/category/application"
	"yadwy-backend/internal/category/domain"
	"yadwy-backend/internal/common"
)

type Handler struct {
	s      *application.CategoryService
	logger *zap.Logger
}

func NewCategoryHandler(createUseCase *application.CategoryService, logger *zap.Logger) *Handler {
	return &Handler{
		s:      createUseCase,
		logger: logger,
	}
}

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
