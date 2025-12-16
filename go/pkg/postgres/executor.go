package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Executor SQLを実行するインターフェース
// pgxpool.Poolえおpgx.Txの共通操作を抽象化
type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// GetExecutor Contextから適切なExecutorを取得する
// トランザクション中ならTxを、そうでなければPoolを返す
func GetExecutor(ctx context.Context, pool *pgxpool.Pool) Executor {
	if tx := TxFromContext(ctx); tx != nil {
		return tx
	}

	return pool
}
