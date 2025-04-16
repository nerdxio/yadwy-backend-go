package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yadwy-backend/internal/cart/application"
	"yadwy-backend/internal/cart/application/mocks"
	"yadwy-backend/internal/cart/domain"
	"yadwy-backend/internal/common"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func setupTestHandler(repo *mocks.CartRepositoryMock) *CartHandler {
	logger := zap.NewNop()
	service := application.NewCartService(repo, logger)
	return NewCartHandler(service, logger)
}

func setupTestContext(t *testing.T) context.Context {
	return context.WithValue(context.Background(), common.AuthKey{}, &common.UserClaims{
		ID:    1,
		Email: "test@example.com",
		Role:  "CUSTOMER",
	})
}

func TestCartHandler_GetCart(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		mock           func() *mocks.CartRepositoryMock
		expectedStatus int
		expectedBody   *domain.Cart
	}{
		{
			name: "should return cart successfully",
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					GetCartFunc: func(ctx context.Context, userID int64) (*domain.Cart, error) {
						return &domain.Cart{
							ID:        1,
							UserID:    userID,
							Items:     []domain.CartItem{},
							CreatedAt: now,
							UpdatedAt: now,
						}, nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody: &domain.Cart{
				ID:        1,
				UserID:    1,
				Items:     []domain.CartItem{},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "should return error when service fails",
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					GetCartFunc: func(ctx context.Context, userID int64) (*domain.Cart, error) {
						return nil, common.NewErrorf(domain.FailedToGetCart, "database error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			handler := setupTestHandler(repo)

			req := httptest.NewRequest(http.MethodGet, "/cart", nil)
			req = req.WithContext(setupTestContext(t))
			rr := httptest.NewRecorder()

			handler.GetCart(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedBody != nil {
				var got domain.Cart
				err := json.NewDecoder(rr.Body).Decode(&got)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if got.ID != tt.expectedBody.ID {
					t.Errorf("handler returned wrong body: got %v want %v", got, tt.expectedBody)
				}
			}
		})
	}
}

func TestCartHandler_AddToCart(t *testing.T) {
	tests := []struct {
		name           string
		request        *addToCartRequest
		mock           func() *mocks.CartRepositoryMock
		expectedStatus int
	}{
		{
			name: "should add item successfully",
			request: &addToCartRequest{
				ProductID: 1,
				Quantity:  2,
			},
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					AddItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return nil
					},
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "should return error when service fails",
			request: &addToCartRequest{
				ProductID: 1,
				Quantity:  2,
			},
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					AddItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return common.NewErrorf(domain.FailedToAddItem, "database error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			handler := setupTestHandler(repo)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/cart/items", bytes.NewBuffer(body))
			req = req.WithContext(setupTestContext(t))
			rr := httptest.NewRecorder()

			handler.AddToCart(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestCartHandler_UpdateCartItem(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		request        *updateCartItemRequest
		mock           func() *mocks.CartRepositoryMock
		expectedStatus int
	}{
		{
			name:      "should update item successfully",
			productID: "1",
			request: &updateCartItemRequest{
				Quantity: 3,
			},
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					UpdateItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "should return error when service fails",
			productID: "1",
			request: &updateCartItemRequest{
				Quantity: 3,
			},
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					UpdateItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return common.NewErrorf(domain.FailedToUpdateItem, "database error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "should return error with invalid product ID",
			productID: "invalid",
			request: &updateCartItemRequest{
				Quantity: 3,
			},
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			handler := setupTestHandler(repo)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPut, "/cart/items/"+tt.productID, bytes.NewBuffer(body))
			req = req.WithContext(setupTestContext(t))

			// Setup Chi URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("productId", tt.productID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			handler.UpdateCartItem(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestCartHandler_RemoveFromCart(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mock           func() *mocks.CartRepositoryMock
		expectedStatus int
	}{
		{
			name:      "should remove item successfully",
			productID: "1",
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					RemoveItemFunc: func(ctx context.Context, userID int64, productID int64) error {
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "should return error when service fails",
			productID: "1",
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					RemoveItemFunc: func(ctx context.Context, userID int64, productID int64) error {
						return common.NewErrorf(domain.FailedToRemoveItem, "database error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "should return error with invalid product ID",
			productID: "invalid",
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			handler := setupTestHandler(repo)

			req := httptest.NewRequest(http.MethodDelete, "/cart/items/"+tt.productID, nil)
			req = req.WithContext(setupTestContext(t))

			// Setup Chi URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("productId", tt.productID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			handler.RemoveFromCart(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestCartHandler_ClearCart(t *testing.T) {
	tests := []struct {
		name           string
		mock           func() *mocks.CartRepositoryMock
		expectedStatus int
	}{
		{
			name: "should clear cart successfully",
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					ClearCartFunc: func(ctx context.Context, userID int64) error {
						return nil
					},
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "should return error when service fails",
			mock: func() *mocks.CartRepositoryMock {
				return &mocks.CartRepositoryMock{
					ClearCartFunc: func(ctx context.Context, userID int64) error {
						return common.NewErrorf(domain.FailedToClearCart, "database error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			handler := setupTestHandler(repo)

			req := httptest.NewRequest(http.MethodDelete, "/cart", nil)
			req = req.WithContext(setupTestContext(t))
			rr := httptest.NewRecorder()

			handler.ClearCart(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
