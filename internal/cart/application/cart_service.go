package application

import (
	"context"
	"yadwy-backend/internal/cart/domain"
	"yadwy-backend/internal/common"

	"go.uber.org/zap"
)

type CartService struct {
	repo   domain.CartRepository
	logger *zap.Logger
}

func NewCartService(repo domain.CartRepository, logger *zap.Logger) *CartService {
	return &CartService{
		repo:   repo,
		logger: logger,
	}
}

func (s *CartService) GetCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	cart, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get cart", zap.Error(err))
		return nil, common.NewErrorf(domain.FailedToGetCart, "failed to get cart: %v", err)
	}
	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	if quantity <= 0 {
		return common.NewErrorf(domain.InvalidQuantityError, "quantity must be greater than 0")
	}

	err := s.repo.AddItem(ctx, userID, productID, quantity)
	if err != nil {
		s.logger.Error("Failed to add item to cart",
			zap.Int64("userID", userID),
			zap.Int64("productID", productID),
			zap.Error(err))
		return common.NewErrorf(domain.FailedToAddItem, "failed to add item to cart: %v", err)
	}
	return nil
}

func (s *CartService) UpdateItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	if quantity <= 0 {
		return common.NewErrorf(domain.InvalidQuantityError, "quantity must be greater than 0")
	}

	err := s.repo.UpdateItem(ctx, userID, productID, quantity)
	if err != nil {
		s.logger.Error("Failed to update cart item",
			zap.Int64("userID", userID),
			zap.Int64("productID", productID),
			zap.Error(err))
		return common.NewErrorf(domain.FailedToUpdateItem, "failed to update cart item: %v", err)
	}
	return nil
}

func (s *CartService) RemoveItem(ctx context.Context, userID int64, productID int64) error {
	err := s.repo.RemoveItem(ctx, userID, productID)
	if err != nil {
		s.logger.Error("Failed to remove item from cart",
			zap.Int64("userID", userID),
			zap.Int64("productID", productID),
			zap.Error(err))
		return common.NewErrorf(domain.FailedToRemoveItem, "failed to remove item from cart: %v", err)
	}
	return nil
}

func (s *CartService) ClearCart(ctx context.Context, userID int64) error {
	err := s.repo.ClearCart(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to clear cart",
			zap.Int64("userID", userID),
			zap.Error(err))
		return common.NewErrorf(domain.FailedToClearCart, "failed to clear cart: %v", err)
	}
	return nil
}
