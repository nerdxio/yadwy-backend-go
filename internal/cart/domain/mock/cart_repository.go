package mock

import (
	"context"
	"yadwy-backend/internal/cart/domain"
)

// CartRepository is a simple mock implementation of domain.CartRepository
type CartRepository struct {
	CreateCartFunc func(ctx context.Context, userID int64) (*domain.Cart, error)
	GetCartFunc    func(ctx context.Context, userID int64) (*domain.Cart, error)
	AddItemFunc    func(ctx context.Context, userID int64, productID int64, quantity int) error
	UpdateItemFunc func(ctx context.Context, userID int64, productID int64, quantity int) error
	RemoveItemFunc func(ctx context.Context, userID int64, productID int64) error
	ClearCartFunc  func(ctx context.Context, userID int64) error
}

func (m *CartRepository) CreateCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	if m.CreateCartFunc != nil {
		return m.CreateCartFunc(ctx, userID)
	}
	return nil, nil
}

func (m *CartRepository) GetCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	if m.GetCartFunc != nil {
		return m.GetCartFunc(ctx, userID)
	}
	return nil, nil
}

func (m *CartRepository) AddItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	if m.AddItemFunc != nil {
		return m.AddItemFunc(ctx, userID, productID, quantity)
	}
	return nil
}

func (m *CartRepository) UpdateItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	if m.UpdateItemFunc != nil {
		return m.UpdateItemFunc(ctx, userID, productID, quantity)
	}
	return nil
}

func (m *CartRepository) RemoveItem(ctx context.Context, userID int64, productID int64) error {
	if m.RemoveItemFunc != nil {
		return m.RemoveItemFunc(ctx, userID, productID)
	}
	return nil
}

func (m *CartRepository) ClearCart(ctx context.Context, userID int64) error {
	if m.ClearCartFunc != nil {
		return m.ClearCartFunc(ctx, userID)
	}
	return nil
}
