package error

// ErrCode エラーコードを表す
// クライアントがエラーハンドリングするために必要な情報を含む
type ErrCode string

const (
	// 400系
	ErrCodeValidation ErrCode = "VALIDATION_ERROR"
	ErrCodeBadRequest ErrCode = "BAD_REQUEST"

	// 404系
	ErrCodeNotFound ErrCode = "NOT_FOUND"

	// 500系
	ErrCodeInternal ErrCode = "INTERNAL_SERVER_ERROR"
)
