package model

type ErrorCode string
type ErrorMessage string

const (
	// ----------------------------------------
	// Auth Errors
	// ----------------------------------------
	ErrUserNotVerified        ErrorCode = "ERROR_USER_NOT_VERIFIED"
	ErrInvalidRefreshToken    ErrorCode = "ERROR_INVALID_REFRESH_TOKEN"
	ErrInvalidAccessToken     ErrorCode = "ERROR_INVALID_ACCESS_TOKEN"
	ErrGeneratingAccessToken  ErrorCode = "ERROR_GENERATING_ACCESS_TOKEN"
	ErrGeneratingRefreshToken ErrorCode = "ERROR_GENERATING_REFRESH_TOKEN"
	ErrInvalidHashingPassword ErrorCode = "ERROR_INVALID_HASHING_PASSWORD"
	ErrUserNotFound           ErrorCode = "ERROR_USER_NOT_FOUND"
	ErrInvalidUserId          ErrorCode = "ERROR_INVALID_USER_ID"
)
