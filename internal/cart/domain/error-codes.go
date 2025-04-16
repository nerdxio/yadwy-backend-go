package domain

import "yadwy-backend/internal/common"

const (
	CartNotFoundError      common.ErrorCode = "cart-not-found"
	CartItemNotFoundError  common.ErrorCode = "cart-item-not-found"
	InvalidQuantityError   common.ErrorCode = "invalid-quantity"
	ProductNotFoundError   common.ErrorCode = "product-not-found"
	FailedToCreateCart     common.ErrorCode = "failed-to-create-cart"
	FailedToAddItem        common.ErrorCode = "failed-to-add-item"
	FailedToUpdateItem     common.ErrorCode = "failed-to-update-item"
	FailedToRemoveItem     common.ErrorCode = "failed-to-remove-item"
	FailedToClearCart      common.ErrorCode = "failed-to-clear-cart"
	FailedToGetCart        common.ErrorCode = "failed-to-get-cart"
	InsufficientStockError common.ErrorCode = "insufficient-stock"
)
