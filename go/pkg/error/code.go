package error

// ErrCode エラーコードを表す
// クライアントがエラーハンドリングするために必要な情報を含む
type ErrCode string

const (
	// 404系
	ErrCodeNotFound ErrCode = "NOT_FOUND"
)
