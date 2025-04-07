package modles

import (
	c "yadwy-backend/internal/common"
)

const (
	UserNotFoundError           c.ErrorCode = "user_not_found"
	EmailAlreadyExistsError     c.ErrorCode = "email_already_exists"
	InvalidUserCredentialsError c.ErrorCode = "invalid_user_credentials"
	InvalidUserRoleError        c.ErrorCode = "invalid_user_role"
)
