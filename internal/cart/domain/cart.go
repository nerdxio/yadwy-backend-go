package domain

import (
	"time"
)

type CartItem struct {
	ID        int64   `json:"id"`
	CartID    int64   `json:"-"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Cart struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *Cart) GetTotalPrice() float64 {
	var total float64
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
