package application

import (
	"context"
	"testing"
	"time"
	"yadwy-backend/internal/cart/domain"
	"yadwy-backend/internal/cart/domain/mock"
	"yadwy-backend/internal/common"

	"go.uber.org/zap"
)

func TestCartService_GetCart(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name    string
		userID  int64
		mock    func() *mock.CartRepository
		want    *domain.Cart
		wantErr bool
	}{
		{
			name:   "should return cart successfully",
			userID: 1,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
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
			want: &domain.Cart{
				ID:        1,
				UserID:    1,
				Items:     []domain.CartItem{},
				CreatedAt: now,
				UpdatedAt: now,
			},
			wantErr: false,
		},
		{
			name:   "should return error when repository fails",
			userID: 1,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					GetCartFunc: func(ctx context.Context, userID int64) (*domain.Cart, error) {
						return nil, common.NewErrorf(domain.FailedToGetCart, "database error")
					},
				}
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			service := NewCartService(repo, logger)

			got, err := service.GetCart(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.GetCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.want.ID {
				t.Errorf("CartService.GetCart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCartService_AddItem(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    int64
		productID int64
		quantity  int
		mock      func() *mock.CartRepository
		wantErr   bool
	}{
		{
			name:      "should add item successfully",
			userID:    1,
			productID: 1,
			quantity:  2,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					AddItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:      "should return error when quantity is invalid",
			userID:    1,
			productID: 1,
			quantity:  0,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{}
			},
			wantErr: true,
		},
		{
			name:      "should return error when repository fails",
			userID:    1,
			productID: 1,
			quantity:  2,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					AddItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return common.NewErrorf(domain.FailedToAddItem, "database error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			service := NewCartService(repo, logger)

			err := service.AddItem(ctx, tt.userID, tt.productID, tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.AddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCartService_UpdateItem(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    int64
		productID int64
		quantity  int
		mock      func() *mock.CartRepository
		wantErr   bool
	}{
		{
			name:      "should update item successfully",
			userID:    1,
			productID: 1,
			quantity:  3,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					UpdateItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:      "should return error when quantity is invalid",
			userID:    1,
			productID: 1,
			quantity:  0,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{}
			},
			wantErr: true,
		},
		{
			name:      "should return error when repository fails",
			userID:    1,
			productID: 1,
			quantity:  3,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					UpdateItemFunc: func(ctx context.Context, userID int64, productID int64, quantity int) error {
						return common.NewErrorf(domain.FailedToUpdateItem, "database error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			service := NewCartService(repo, logger)

			err := service.UpdateItem(ctx, tt.userID, tt.productID, tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCartService_RemoveItem(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    int64
		productID int64
		mock      func() *mock.CartRepository
		wantErr   bool
	}{
		{
			name:      "should remove item successfully",
			userID:    1,
			productID: 1,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					RemoveItemFunc: func(ctx context.Context, userID int64, productID int64) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:      "should return error when repository fails",
			userID:    1,
			productID: 1,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					RemoveItemFunc: func(ctx context.Context, userID int64, productID int64) error {
						return common.NewErrorf(domain.FailedToRemoveItem, "database error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			service := NewCartService(repo, logger)

			err := service.RemoveItem(ctx, tt.userID, tt.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCartService_ClearCart(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  int64
		mock    func() *mock.CartRepository
		wantErr bool
	}{
		{
			name:   "should clear cart successfully",
			userID: 1,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					ClearCartFunc: func(ctx context.Context, userID int64) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "should return error when repository fails",
			userID: 1,
			mock: func() *mock.CartRepository {
				return &mock.CartRepository{
					ClearCartFunc: func(ctx context.Context, userID int64) error {
						return common.NewErrorf(domain.FailedToClearCart, "database error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mock()
			service := NewCartService(repo, logger)

			err := service.ClearCart(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.ClearCart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
