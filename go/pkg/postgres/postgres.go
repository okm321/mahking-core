package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB configuration
type DB struct {
	Host     string `env:"PG_HOST" envDefault:"127.0.0.1"`
	Port     string `env:"PG_PORT" envDefault:"5432"`
	User     string `env:"PG_USER" envDefault:"postgres"`
	Pass     string `env:"PG_PASS" envDefault:"password"`
	DbName   string `env:"PG_DBNAME" envDefault:"postgres"`
	Schema   string `env:"PG_SCHEMA" envDefault:"mahking_local"`
	Params   string `env:"PG_PARAMS" envDefault:"sslmode=disable timezone=Asia/Tokyo lock_timeout=50000"`
	MinConns int    `env:"MIN_CONNS" envDefault:"2"`
	MaxConns int    `env:"MAX_CONNS" envDefault:"10"`
}

// Connect 新しいプールを作成（シングルトン管理なし）
func Connect(cnf DB) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(DSN(cnf))
	if err != nil {
		return nil, err
	}

	if cnf.MinConns > 0 {
		cfg.MinConns = int32(cnf.MinConns)
	}
	if cnf.MaxConns > 0 {
		cfg.MaxConns = int32(cnf.MaxConns)
	}

	cfg.PrepareConn = prepareSession
	cfg.AfterRelease = resetGroupSetting

	return pgxpool.NewWithConfig(context.Background(), cfg)
}

func DSN(c DB) string {
	parts := []string{
		fmt.Sprintf("user=%s", c.User),
		fmt.Sprintf("password=%s", c.Pass),
		fmt.Sprintf("database=%s", c.DbName),
		fmt.Sprintf("host=%s", strings.TrimSpace(c.Host)),
		fmt.Sprintf("port=%s", strings.TrimSpace(c.Port)),
		fmt.Sprintf("search_path=%s", strings.TrimSpace(c.Schema)),
	}

	params := strings.TrimSpace(c.Params)
	if params != "" {
		parts = append(parts, params)
	}

	return strings.Join(parts, " ")
}

const groupSettingName = "app.group_id"

func prepareSession(ctx context.Context, conn *pgx.Conn) (bool, error) {
	session := sessionFromContext(ctx)
	switch session.kind {
	case sessionKindGroup:
		if err := setGroupConfig(ctx, conn, session.groupID); err != nil {
			return false, err
		}
	default:
		if err := clearGroupConfig(ctx, conn); err != nil {
			return false, err
		}
	}
	return true, nil
}

func resetGroupSetting(conn *pgx.Conn) bool {
	return clearGroupConfig(context.Background(), conn) == nil
}

func setGroupConfig(ctx context.Context, conn *pgx.Conn, groupID string) error {
	_, err := conn.Exec(ctx, "select set_config($1, $2, false)", groupSettingName, groupID)
	return err
}

func clearGroupConfig(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, "select set_config($1, '', false)", groupSettingName)
	return err
}

type sessionKind int

const (
	sessionKindDefault sessionKind = iota
	sessionKindGroup
)

type sessionContext struct {
	kind    sessionKind
	groupID string
}

type ctxKeySession struct{}

func sessionFromContext(ctx context.Context) sessionContext {
	if ctx == nil {
		return sessionContext{kind: sessionKindDefault}
	}
	val := ctx.Value(ctxKeySession{})
	if val == nil {
		return sessionContext{kind: sessionKindDefault}
	}
	if sess, ok := val.(sessionContext); ok {
		return sess
	}
	return sessionContext{kind: sessionKindDefault}
}

// WithGroupID marks the context to run in a group session (RLS enforced)
func WithGroupID(ctx context.Context, groupID string) context.Context {
	return context.WithValue(ctx, ctxKeySession{}, sessionContext{
		kind:    sessionKindGroup,
		groupID: groupID,
	})
}
