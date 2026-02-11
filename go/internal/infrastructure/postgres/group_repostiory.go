package postgres

import (
	"context"
	"fmt"

	"github.com/guregu/null/v6"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okm321/mahking-go/internal/domain"
	"github.com/okm321/mahking-go/internal/infrastructure/postgres/sqlc"
	pkgerror "github.com/okm321/mahking-go/pkg/error"
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

func (r *GroupRepository) Create(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	row, err := r.q.CreateGroup(ctx, group.Name)
	if err != nil {
		return nil, err
	}

	group.ID = row.ID

	err = r.createRelatedInfo(ctx, group)
	if err != nil {
		return nil, err
	}

	return &domain.Group{
		ID:   row.ID,
		UID:  row.Uid.String(),
		Name: row.Name,
	}, nil
}

func (r *GroupRepository) createRelatedInfo(ctx context.Context, group *domain.Group) (err error) {
	memberParams := make([]sqlc.CreateMembersParams, 0, len(group.Members))
	for _, m := range group.Members {
		memberParams = append(memberParams, sqlc.CreateMembersParams{
			GroupID: group.ID,
			Name:    m.Name,
		})
	}
	_, err = r.q.CreateMembers(ctx, memberParams)
	if err != nil {
		return pkgerror.Wrap(err, "create related members")
	}

	ruleParam := sqlc.CreateRuleParams{
		GroupID:               group.ID,
		MahjongType:           int32(group.Rule.MahjongType),           //nolint:gosec // 麻雀タイプは1-2の範囲
		InitialPoints:         int32(group.Rule.InitialPoints),         //nolint:gosec // 点数はint32範囲内
		ReturnPoints:          int32(group.Rule.ReturnPoints),          //nolint:gosec // 点数はint32範囲内
		RankingPointsFirst:    int32(group.Rule.RankingPointsFirst),    //nolint:gosec // 点数はint32範囲内
		RankingPointsSecond:   int32(group.Rule.RankingPointsSecond),   //nolint:gosec // 点数はint32範囲内
		RankingPointsThird:    int32(group.Rule.RankingPointsThird),    //nolint:gosec // 点数はint32範囲内
		RankingPointsFourth:   null.IntFromPtr(group.Rule.RankingPointsFour.Ptr()),
		FractionalCalculation: int32(group.Rule.FractionalCalculation), //nolint:gosec // 計算方法は1-5の範囲
		UseBust:               group.Rule.UseBust,
		BustPoint:             null.IntFromPtr(group.Rule.BustPoint.Ptr()),
		UseChip:               group.Rule.UseChip,
		ChipPoint:             null.IntFromPtr(group.Rule.ChipPoint.Ptr()),
	}
	_, err = r.q.CreateRule(ctx, ruleParam)
	if err != nil {
		return pkgerror.Wrap(err, "create related rule")
	}

	return nil
}

var _ domain.GroupRepository = (*GroupRepository)(nil)
