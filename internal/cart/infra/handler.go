package infra

import (
	"encoding/json"
	"net/http"
	"strconv"
	"yadwy-backend/internal/cart/application"
	"yadwy-backend/internal/cart/domain"
	"yadwy-backend/internal/common"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type CartHandler struct {
	service *application.CartService
	logger  *zap.Logger
}

type addToCartRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required,gt=0"`
}

type updateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

type cartResponse struct {
	*domain.Cart
	Total float64 `json:"total"`
}

func NewCartHandler(service *application.CartService, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	claims, err := common.GetLoggedInUser(r)
	if err != nil {
		h.logger.Error("Failed to get logged in user", zap.Error(err))
		common.SendError(w, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	cart, err := h.service.GetCart(r.Context(), claims.ID)
	if err != nil {
		h.logger.Error("Failed to get cart", zap.Error(err))
		common.SendError(w, http.StatusInternalServerError, string(domain.FailedToGetCart), err.Error())
		return
	}

	response := cartResponse{
		Cart:  cart,
		Total: cart.GetTotalPrice(),
	}

	if err = common.Encode(w, http.StatusOK, response); err != nil {
		h.logger.Error("Failed to encode cart", zap.Error(err))
		common.SendError(w, http.StatusInternalServerError, "encode-error", "failed to encode response")
		return
	}
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req addToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.SendError(w, http.StatusBadRequest, "invalid-request", "invalid request body")
		return
	}

	claims, err := common.GetLoggedInUser(r)
	if err != nil {
		common.SendError(w, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	err = h.service.AddItem(r.Context(), claims.ID, req.ProductID, req.Quantity)
	if err != nil {
		h.logger.Error("Failed to add item to cart", zap.Error(err))
		common.SendError(w, http.StatusInternalServerError, string(domain.FailedToAddItem), err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		common.SendError(w, http.StatusBadRequest, "invalid-product-id", "invalid product ID")
		return
	}

	var req updateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.SendError(w, http.StatusBadRequest, "invalid-request", "invalid request body")
		return
	}

	claims, err := common.GetLoggedInUser(r)
	if err != nil {
		common.SendError(w, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	err = h.service.UpdateItem(r.Context(), claims.ID, productID, req.Quantity)
	if err != nil {
		h.logger.Error("Failed to update cart item", zap.Error(err))
		common.SendError(w, http.StatusInternalServerError, string(domain.FailedToUpdateItem), err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(chi.URLParam(r, "productId"), 10, 64)
	if err != nil {
		common.SendError(w, http.StatusBadRequest, "invalid-product-id", "invalid product ID")
		return
	}

	claims, err := common.GetLoggedInUser(r)
	if err != nil {
		common.SendError(w, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	err = h.service.RemoveItem(r.Context(), claims.ID, productID)
	if err != nil {
		h.logger.Error("Failed to remove item from cart", zap.Error(err))
		common.SendError(w, http.StatusInternalServerError, string(domain.FailedToRemoveItem), err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	claims, err := common.GetLoggedInUser(r)
	if err != nil {
		common.SendError(w, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	err = h.service.ClearCart(r.Context(), claims.ID)
	if err != nil {
		h.logger.Error("Failed to clear cart", zap.Error(err))
		common.SendError(w, http.StatusInternalServerError, string(domain.FailedToClearCart), err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoadCartRoutes(db *sqlx.DB, logger *zap.Logger, jwt *common.JWTGenerator) http.Handler {
	router := chi.NewRouter()
	repo := NewCartRepository(db, logger)
	service := application.NewCartService(repo, logger)
	handler := NewCartHandler(service, logger)

	router.Use(common.GetAuthMiddlewareFunc(jwt))

	router.Get("/", handler.GetCart)
	router.Post("/items", handler.AddToCart)
	router.Put("/items/{productId}", handler.UpdateCartItem)
	router.Delete("/items/{productId}", handler.RemoveFromCart)
	router.Delete("/", handler.ClearCart)

	return router
}
