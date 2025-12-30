package postgres

import (
	"context"

	pkgerror "github.com/okm321/mahking-go/pkg/error"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okm321/mahking-go/internal/domain"
	"github.com/okm321/mahking-go/internal/infrastructure/postgres/sqlc"
)

type MemberRepository struct {
	q *sqlc.Queries
}

func NewMemberRepository(pool *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{
		q: sqlc.New(pool),
	}
}

func (r *MemberRepository) BatchCreateMembers(ctx context.Context, members []*domain.Member) error {
	params := make([]sqlc.CreateMembersParams, 0, len(members))
	for _, m := range members {
		params = append(params, sqlc.CreateMembersParams{
			GroupID: m.GroupID,
			Name:    m.Name,
		})
	}

	_, err := r.q.CreateMembers(ctx, params)
	if err != nil {
		return pkgerror.Wrapf(err, "batch create members")
	}

	return nil
}
