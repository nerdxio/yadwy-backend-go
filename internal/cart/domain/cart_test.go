package domain

import (
	"testing"
	"time"
)

func TestCart_GetTotalPrice(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		cart Cart
		want float64
	}{
		{
			name: "empty cart should return 0",
			cart: Cart{
				ID:        1,
				UserID:    1,
				Items:     []CartItem{},
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: 0,
		},
		{
			name: "cart with single item should return item price * quantity",
			cart: Cart{
				ID:     1,
				UserID: 1,
				Items: []CartItem{
					{
						ID:        1,
						CartID:    1,
						ProductID: 1,
						Quantity:  2,
						Price:     10.5,
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: 21.0,
		},
		{
			name: "cart with multiple items should return sum of (price * quantity)",
			cart: Cart{
				ID:     1,
				UserID: 1,
				Items: []CartItem{
					{
						ID:        1,
						CartID:    1,
						ProductID: 1,
						Quantity:  2,
						Price:     10.5,
					},
					{
						ID:        2,
						CartID:    1,
						ProductID: 2,
						Quantity:  1,
						Price:     15.75,
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: 36.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cart.GetTotalPrice()
			if got != tt.want {
				t.Errorf("Cart.GetTotalPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
