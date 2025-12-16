package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// ctxKeyTx contextにTxを格納するためのキー
type ctxKeyTx struct{}

// TxToContext ContextにTxをセットする
func TxToContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

// TxFromContext ContextからTxを取得する
// トランザクション外の場合はnilを返す
func TxFromContext(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(ctxKeyTx{}).(pgx.Tx)
	return tx
}
