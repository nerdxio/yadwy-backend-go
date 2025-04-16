package infra

import (
	"context"
	"database/sql"
	"time"
	"yadwy-backend/internal/cart/domain"
	"yadwy-backend/internal/common"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type CartRepositoryImpl struct {
	db     *sqlx.DB
	logger *zap.Logger
}

type cartDbo struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type cartItemDbo struct {
	ID        int64   `db:"id"`
	CartID    int64   `db:"cart_id"`
	ProductID int64   `db:"product_id"`
	Quantity  int     `db:"quantity"`
	Price     float64 `db:"price"`
}

func NewCartRepository(db *sqlx.DB, logger *zap.Logger) domain.CartRepository {
	return &CartRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *CartRepositoryImpl) CreateCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	query := `INSERT INTO carts (user_id) VALUES ($1) RETURNING id, user_id, created_at, updated_at`
	var cart cartDbo
	err := r.db.QueryRowxContext(ctx, query, userID).StructScan(&cart)
	if err != nil {
		return nil, common.NewErrorf(domain.FailedToCreateCart, "failed to create cart: %v", err)
	}

	return &domain.Cart{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     []domain.CartItem{},
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}, nil
}

func (r *CartRepositoryImpl) GetCart(ctx context.Context, userID int64) (*domain.Cart, error) {
	var cart cartDbo
	err := r.db.GetContext(ctx, &cart,
		"SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1", userID)
	if err == sql.ErrNoRows {
		return r.CreateCart(ctx, userID)
	}
	if err != nil {
		return nil, common.NewErrorf(domain.FailedToGetCart, "failed to get cart: %v", err)
	}

	var items []cartItemDbo
	err = r.db.SelectContext(ctx, &items,
		"SELECT * FROM cart_items WHERE cart_id = $1", cart.ID)
	if err != nil {
		return nil, common.NewErrorf(domain.FailedToGetCart, "failed to get cart items: %v", err)
	}

	domainItems := make([]domain.CartItem, len(items))
	for i, item := range items {
		domainItems[i] = domain.CartItem{
			ID:        item.ID,
			CartID:    item.CartID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &domain.Cart{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     domainItems,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}, nil
}

func (r *CartRepositoryImpl) AddItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return common.NewErrorf(domain.FailedToAddItem, "failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	var cartID int64
	err = tx.GetContext(ctx, &cartID, "SELECT id FROM carts WHERE user_id = $1", userID)
	if err == sql.ErrNoRows {
		var cart cartDbo
		err = tx.QueryRowxContext(ctx,
			"INSERT INTO carts (user_id) VALUES ($1) RETURNING id", userID).Scan(&cart.ID)
		if err != nil {
			return common.NewErrorf(domain.FailedToCreateCart, "failed to create cart: %v", err)
		}
		cartID = cart.ID
	} else if err != nil {
		return common.NewErrorf(domain.FailedToGetCart, "failed to get cart: %v", err)
	}

	var price float64
	err = tx.GetContext(ctx, &price, "SELECT price FROM products WHERE id = $1", productID)
	if err == sql.ErrNoRows {
		return common.NewErrorf(domain.ProductNotFoundError, "product not found")
	}
	if err != nil {
		return common.NewErrorf(domain.FailedToAddItem, "failed to get product price: %v", err)
	}

	var existingItem cartItemDbo
	err = tx.GetContext(ctx, &existingItem,
		"SELECT * FROM cart_items WHERE cart_id = $1 AND product_id = $2", cartID, productID)
	if err == sql.ErrNoRows {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO cart_items (cart_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
			cartID, productID, quantity, price)
	} else if err == nil {
		_, err = tx.ExecContext(ctx,
			"UPDATE cart_items SET quantity = quantity + $1 WHERE cart_id = $2 AND product_id = $3",
			quantity, cartID, productID)
	}
	if err != nil {
		return common.NewErrorf(domain.FailedToAddItem, "failed to add item to cart: %v", err)
	}

	return tx.Commit()
}

func (r *CartRepositoryImpl) UpdateItem(ctx context.Context, userID int64, productID int64, quantity int) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE cart_items ci
		SET quantity = $1
		FROM carts c
		WHERE c.id = ci.cart_id
		AND c.user_id = $2
		AND ci.product_id = $3`,
		quantity, userID, productID)

	if err != nil {
		return common.NewErrorf(domain.FailedToUpdateItem, "failed to update item: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return common.NewErrorf(domain.FailedToUpdateItem, "failed to get affected rows: %v", err)
	}
	if rows == 0 {
		return common.NewErrorf(domain.CartItemNotFoundError, "cart item not found")
	}

	return nil
}

func (r *CartRepositoryImpl) RemoveItem(ctx context.Context, userID int64, productID int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM cart_items ci
		USING carts c
		WHERE c.id = ci.cart_id
		AND c.user_id = $1
		AND ci.product_id = $2`,
		userID, productID)

	if err != nil {
		return common.NewErrorf(domain.FailedToRemoveItem, "failed to remove item: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return common.NewErrorf(domain.FailedToRemoveItem, "failed to get affected rows: %v", err)
	}
	if rows == 0 {
		return common.NewErrorf(domain.CartItemNotFoundError, "cart item not found")
	}

	return nil
}

func (r *CartRepositoryImpl) ClearCart(ctx context.Context, userID int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return common.NewErrorf(domain.FailedToClearCart, "failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE FROM cart_items ci
		USING carts c
		WHERE c.id = ci.cart_id
		AND c.user_id = $1`,
		userID)

	if err != nil {
		return common.NewErrorf(domain.FailedToClearCart, "failed to clear cart items: %v", err)
	}

	return tx.Commit()
}
