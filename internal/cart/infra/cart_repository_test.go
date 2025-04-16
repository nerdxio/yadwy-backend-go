package infra

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"yadwy-backend/internal/cart/domain"
	"yadwy-backend/internal/common"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type CartRepositoryTestSuite struct {
	suite.Suite
	db   *sqlx.DB
	mock sqlmock.Sqlmock
	repo domain.CartRepository
}

func (s *CartRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	logger := zap.NewNop()
	s.db = sqlxDB
	s.mock = mock
	s.repo = NewCartRepository(sqlxDB, logger)
}

func (s *CartRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}

func TestCartRepository(t *testing.T) {
	suite.Run(t, new(CartRepositoryTestSuite))
}

func (s *CartRepositoryTestSuite) TestGetCart() {
	ctx := context.Background()
	now := time.Now()
	cartID := int64(1)
	userID := int64(1)

	s.Run("should get cart with items successfully", func() {
		// Prepare cart query mock
		cartRows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
			AddRow(cartID, userID, now, now)
		s.mock.ExpectQuery("SELECT (.+) FROM carts").
			WithArgs(userID).
			WillReturnRows(cartRows)

		// Prepare cart items query mock
		itemRows := sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity", "price"}).
			AddRow(1, cartID, 1, 2, 10.50).
			AddRow(2, cartID, 2, 1, 15.75)
		s.mock.ExpectQuery("SELECT (.+) FROM cart_items").
			WithArgs(cartID).
			WillReturnRows(itemRows)

		cart, err := s.repo.GetCart(ctx, userID)
		s.Require().NoError(err)
		s.Equal(cartID, cart.ID)
		s.Equal(userID, cart.UserID)
		s.Len(cart.Items, 2)
		s.Equal(10.50, cart.Items[0].Price)
		s.Equal(15.75, cart.Items[1].Price)
	})

	s.Run("should create new cart when not found", func() {
		s.mock.ExpectQuery("SELECT (.+) FROM carts").
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		s.mock.ExpectQuery("INSERT INTO carts").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
				AddRow(cartID, userID, now, now))

		cart, err := s.repo.GetCart(ctx, userID)
		s.Require().NoError(err)
		s.Equal(cartID, cart.ID)
		s.Equal(userID, cart.UserID)
		s.Empty(cart.Items)
	})
}

func (s *CartRepositoryTestSuite) TestAddItem() {
	ctx := context.Background()
	userID := int64(1)
	cartID := int64(1)
	productID := int64(1)
	quantity := 2
	price := 10.50

	s.Run("should add new item to cart successfully", func() {
		// Mock transaction
		s.mock.ExpectBegin()

		// Mock getting cart ID
		s.mock.ExpectQuery("SELECT id FROM carts").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(cartID))

		// Mock getting product price
		s.mock.ExpectQuery("SELECT price FROM products").
			WithArgs(productID).
			WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))

		// Mock checking if item exists
		s.mock.ExpectQuery("SELECT (.+) FROM cart_items").
			WithArgs(cartID, productID).
			WillReturnError(sql.ErrNoRows)

		// Mock inserting new item
		s.mock.ExpectExec("INSERT INTO cart_items").
			WithArgs(cartID, productID, quantity, price).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock commit
		s.mock.ExpectCommit()

		err := s.repo.AddItem(ctx, userID, productID, quantity)
		s.Require().NoError(err)
	})

	s.Run("should update quantity when item exists", func() {
		// Mock transaction
		s.mock.ExpectBegin()

		// Mock getting cart ID
		s.mock.ExpectQuery("SELECT id FROM carts").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(cartID))

		// Mock getting product price
		s.mock.ExpectQuery("SELECT price FROM products").
			WithArgs(productID).
			WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(price))

		// Mock checking if item exists
		s.mock.ExpectQuery("SELECT (.+) FROM cart_items").
			WithArgs(cartID, productID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product_id", "quantity", "price"}).
				AddRow(1, cartID, productID, 1, price))

		// Mock updating existing item
		s.mock.ExpectExec("UPDATE cart_items SET quantity").
			WithArgs(quantity, cartID, productID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock commit
		s.mock.ExpectCommit()

		err := s.repo.AddItem(ctx, userID, productID, quantity)
		s.Require().NoError(err)
	})
}

func (s *CartRepositoryTestSuite) TestUpdateItem() {
	ctx := context.Background()
	userID := int64(1)
	productID := int64(1)
	quantity := 3

	s.Run("should update item quantity successfully", func() {
		s.mock.ExpectExec("UPDATE cart_items").
			WithArgs(quantity, userID, productID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.repo.UpdateItem(ctx, userID, productID, quantity)
		s.Require().NoError(err)
	})

	s.Run("should return error when item not found", func() {
		s.mock.ExpectExec("UPDATE cart_items").
			WithArgs(quantity, userID, productID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.repo.UpdateItem(ctx, userID, productID, quantity)
		s.Require().Error(err)
		if e, ok := err.(*common.Error); ok {
			s.Equal(domain.CartItemNotFoundError, e.Code())
		} else {
			s.Fail("Expected *common.Error type")
		}
	})
}

func (s *CartRepositoryTestSuite) TestRemoveItem() {
	ctx := context.Background()
	userID := int64(1)
	productID := int64(1)

	s.Run("should remove item successfully", func() {
		s.mock.ExpectExec("DELETE FROM cart_items").
			WithArgs(userID, productID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.repo.RemoveItem(ctx, userID, productID)
		s.Require().NoError(err)
	})

	s.Run("should return error when item not found", func() {
		s.mock.ExpectExec("DELETE FROM cart_items").
			WithArgs(userID, productID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.repo.RemoveItem(ctx, userID, productID)
		s.Require().Error(err)
		if e, ok := err.(*common.Error); ok {
			s.Equal(domain.CartItemNotFoundError, e.Code())
		} else {
			s.Fail("Expected *common.Error type")
		}
	})
}

func (s *CartRepositoryTestSuite) TestClearCart() {
	ctx := context.Background()
	userID := int64(1)

	s.Run("should clear cart successfully", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("DELETE FROM cart_items").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repo.ClearCart(ctx, userID)
		s.Require().NoError(err)
	})

	s.Run("should rollback transaction on error", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("DELETE FROM cart_items").
			WithArgs(userID).
			WillReturnError(sql.ErrConnDone)
		s.mock.ExpectRollback()

		err := s.repo.ClearCart(ctx, userID)
		s.Require().Error(err)
		if e, ok := err.(*common.Error); ok {
			s.Equal(domain.FailedToClearCart, e.Code())
		} else {
			s.Fail("Expected *common.Error type")
		}
	})
}
