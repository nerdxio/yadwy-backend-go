package category

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"yadwy-backend/internal/common"
)

type Handler struct {
	s      *CategoryService
	logger *zap.Logger
}

func NewCategoryHandler(createUseCase *CategoryService, logger *zap.Logger) *Handler {
	return &Handler{
		s:      createUseCase,
		logger: logger,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Create category request", zap.String("method", "Create"))
	_ = r.ParseMultipartForm(128 * 1024 * 1024) // 128 MB limit
	_, file, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("Failed to parse form file", zap.Error(err))
		http.Error(w, "Failed to parse form file", http.StatusBadRequest)
		return
	}

	//req, err := common.DecodeAndValidate[CreateCategoryRequest](r)
	//if err != nil {
	//	h.logger.Error("Failed to decode and validate request", zap.Error(err))
	//	http.Error(w, "Invalid request", http.StatusBadRequest)
	//	return
	//}

	name := r.FormValue("name")
	description := r.FormValue("description")
	err = h.s.Execute(name, description, file)
	if err != nil {
		h.logger.Warn("Failed to create category", zap.Error(err))
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoadCategoryRoutes(b *sqlx.DB, logger *zap.Logger) http.Handler {
	ar := chi.NewRouter()
	ar.Use(AdminOnly)
	cr := NewCategoryRepo(b, logger)
	files, _ := common.NewFileService("/home/nerd/images", "http://localhost:3000/images")
	cs := NewCategoryService(cr, files, logger)
	ch := NewCategoryHandler(cs, logger)
	ar.Post("/", ch.Create)
	return ar
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("AdminOnly")

		next.ServeHTTP(w, r)
	})
}
