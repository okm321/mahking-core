package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okm321/mahking-go/internal/domain"
	"github.com/okm321/mahking-go/internal/infrastructure/postgres/sqlc"
)

type GroupRepository struct {
	q *sqlc.Queries
}

func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{
		q: sqlc.New(pool),
	}
}

func (r *GroupRepository) List(ctx context.Context) ([]domain.Group, error) {
	rows, err := r.q.ListGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}

	groups := make([]domain.Group, 0, len(rows))
	for _, row := range rows {
		groups = append(groups, domain.Group{
			ID:   row.ID,
			UID:  row.Uid.String(),
			Name: row.Name,
		})
	}
	return groups, nil
}

func (r *GroupRepository) Create(ctx context.Context, name string) (domain.Group, error) {
	row, err := r.q.CreateGroup(ctx, name)
	if err != nil {
		return domain.Group{}, fmt.Errorf("create group: %w", err)
	}

	return domain.Group{
		ID:   row.ID,
		UID:  row.Uid.String(),
		Name: row.Name,
	}, nil
}

var _ domain.GroupRepository = (*GroupRepository)(nil)
