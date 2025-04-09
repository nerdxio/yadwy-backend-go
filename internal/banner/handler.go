package banner

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"yadwy-backend/internal/common"
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
