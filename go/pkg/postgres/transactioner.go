package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okm321/mahking-go/internal/domain"
)

type Transactioner struct {
	pool *pgxpool.Pool
}

func NewTransactioner(pool *pgxpool.Pool) *Transactioner {
	return &Transactioner{pool: pool}
}

// BeginTx トランザクションを開始する
func (t *Transactioner) BeginTx(ctx context.Context) (domain.Tx, context.Context, error) {
	pgxTx, err := t.pool.Begin(ctx)
	if err != nil {
		return nil, ctx, fmt.Errorf("failed to begin transaction: %w", err)
	}

	tx := NewTx(pgxTx)
	ctx = TxToContext(ctx, pgxTx)

	return tx, ctx, nil
}

// Commit トランザクションをコミットする
func (t *Transactioner) Commit(tx domain.Tx) error {
	rawTx, err := t.extractRawTx(tx)
	if err != nil {
		return err
	}

	if err := rawTx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// Rollback トランザクションをロールバックする
func (t *Transactioner) Rollback(tx domain.Tx) error {
	rawTx, err := t.extractRawTx(tx)
	if err != nil {
		return err
	}

	if err := rawTx.Rollback(context.Background()); err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	return nil
}

// extractRawTx domain.Txからpgx.Txを取り出す
func (t *Transactioner) extractRawTx(tx domain.Tx) (pgx.Tx, error) {
	pgTx, ok := tx.(*Tx)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type: expected *postgres.Tx, got %T", tx)
	}
	return pgTx.Raw(), nil
}

// インターフェースを満たすことをコンパイル時にチェック
var _ domain.Transactioner = (*Transactioner)(nil)
