package domain

import (
	"context"
)

type Tx any

type Transactioner interface {
	// BeginTx トランザクションを開始し、TxをContextにセットして返す
	BeginTx(ctx context.Context) (Tx, context.Context, error)
	// Commit トランザクションをコミット
	Commit(tx Tx) error
	// Rollback トランザクションをロールバック
	Rollback(tx Tx) error
}

// WithTransaction トランザクション内で処理を実行するヘルパー関数
// fnがerrorを返した場合は自動でRokkback, nilの場合は自動でCommit
func WithTransaction(ctx context.Context, t Transactioner, fn func(ctx context.Context) error) error {
	tx, ctx, err := t.BeginTx(ctx)
	if err != nil {
		return err
	}

	err = fn(ctx)
	if err != nil {
		// fnのエラーを優先するため、Rollbackのエラーは無視
		_ = t.Rollback(tx)
		return err
	}

	return t.Commit(tx)
}
