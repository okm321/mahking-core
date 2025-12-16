package postgres

import "github.com/jackc/pgx/v5"

// Tx domain.Txを満たすpgxトランザクションのラッパー
type Tx struct {
	raw pgx.Tx
}

// NewTx Txを作成する
func NewTx(raw pgx.Tx) *Tx {
	return &Tx{raw: raw}
}

// Raw内部のpgx.Txを取得する（pkg内での使用を想定）
func (t *Tx) Raw() pgx.Tx {
	return t.raw
}
