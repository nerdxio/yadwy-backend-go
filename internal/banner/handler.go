package banner

import (
	"errors"
	"net/http"
	"strconv"
	"yadwy-backend/internal/common"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Handler struct {
	svc    *Service
	logger *zap.Logger
}

func NewHandler(svc *Service, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

// @Summary Get all banners
// @Description Get a list of all banners
// @Tags banners
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Banner
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Banners not found"
// @Router /banners [get]
func (h *Handler) GetBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.svc.GetBanners(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	if err = common.Encode(w, http.StatusOK, banners); err != nil {
		handleError(w, err)
		return
	}
}

// @Summary Create a new banner
// @Description Create a new banner with image upload (Admin only)
// @Tags banners
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Banner name"
// @Param index formData integer true "Banner display index"
// @Param image formData file true "Banner image"
// @Success 201 {object} Banner
// @Failure 400 {object} common.ErrorResponse "Invalid input"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 403 {object} common.ErrorResponse "Forbidden - Admin only"
// @Router /banners [post]
func (h *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(128 * 1024 * 1024); err != nil {
		handleError(w, common.NewErrorf(FailedToParseImage, "%v", err))
		return
	}

	_, file, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("Failed to parse form file", zap.Error(err))
		handleError(w, common.NewErrorf(FailedToParseImage, "%v", err))
		return
	}

	name := r.FormValue("name")
	indexStr := r.FormValue("index")
	index, err := strconv.Atoi(indexStr)

	banner, err := h.svc.CreateBanner(r.Context(), name, file, index)
	if err != nil {
		handleError(w, err)
		return
	}

	if err = common.Encode(w, http.StatusCreated, banner); err != nil {
		handleError(w, err)
		return
	}
}

func handleError(w http.ResponseWriter, err error) {
	var appErr *common.Error
	if errors.As(err, &appErr) {
		switch appErr.Code() {
		case FailedToGetAllBanners:
			common.SendError(w, http.StatusNotFound, string(appErr.Code()), appErr.Error())
		case FailedToParseImage:
			common.SendError(w, http.StatusBadRequest, string(appErr.Code()), appErr.Error())
		default:
			common.SendError(w, http.StatusInternalServerError, string(appErr.Code()), appErr.Error())
		}
		return
	}

	common.SendError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
}

func LoadBannerRoutes(b *sqlx.DB, logger *zap.Logger, jwt *common.JWTGenerator) http.Handler {
	ar := chi.NewRouter()
	br := NewRepo(b, logger)
	files, _ := common.NewFileService("/home/nerd/images", "http://localhost:3000/images")
	svc := NewService(br, logger, files)
	handler := NewHandler(svc, logger)

	ar.Use(common.GetAuthMiddlewareFunc(jwt))

	ar.Get("/", handler.GetBanners)
	ar.With(common.GetAdminMiddlewareFun(jwt)).Post("/", handler.CreateBanner)
	return ar
}
