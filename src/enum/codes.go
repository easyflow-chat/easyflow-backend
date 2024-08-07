package enum

// ErrorCode represents an error code as a string.
type ErrorCode string

// Error codes constants, unexported.
const (
	Unauthorized        ErrorCode = "UNAUTHORIZED"
	ApiError            ErrorCode = "API_ERROR"
	NotAllowed          ErrorCode = "NOT_ALLOWED"
	NotFound            ErrorCode = "NOT_FOUND"
	AlreadyExists       ErrorCode = "ALREADY_EXISTS"
	WrongCredentials    ErrorCode = "WRONG_CREDENTIALS"
	MalformedRequest    ErrorCode = "MALFORMED_REQUEST"
	InvalidCookie       ErrorCode = "INVALID_COOKIE"
	InvalidRefresh      ErrorCode = "INVALID_REFRESH"
	ExpiredToken        ErrorCode = "EXPIRED_TOKEN"
	ExpiredRefreshToken ErrorCode = "EXPIRED_REFRESH"
	UserNotFound        ErrorCode = "USER_NOT_FOUND"
)
