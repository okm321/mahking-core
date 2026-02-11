package application

import (
	"context"

	appin "github.com/okm321/mahking-go/internal/application/in"
	appout "github.com/okm321/mahking-go/internal/application/out"
	"github.com/okm321/mahking-go/internal/domain"
	pkgerror "github.com/okm321/mahking-go/pkg/error"
	pkgtrace "github.com/okm321/mahking-go/pkg/trace"
)

type GroupUsecase struct {
	groupRepo  domain.GroupRepository
	memberRepo domain.MemberRepository
}

type NewGroupUsecaseArgs struct {
	GroupRepo  domain.GroupRepository
	MemberRepo domain.MemberRepository
}

func NewGroupUsecase(args *NewGroupUsecaseArgs) *GroupUsecase {
	return &GroupUsecase{
		groupRepo:  args.GroupRepo,
		memberRepo: args.MemberRepo,
	}
}

func (u *GroupUsecase) List(ctx context.Context) (res []appout.Group, err error) {
	ctx = pkgtrace.StartSpan(ctx, "GroupUsecase.List")
	defer func() { pkgtrace.EndSpan(ctx, err) }()

	groups, err := u.groupRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	res = make([]appout.Group, 0, len(groups))
	for _, g := range groups {
		res = append(res, appout.NewGroup(g))
	}
	return res, nil
}

func (u *GroupUsecase) Create(ctx context.Context, in appin.CreateGroupWithRule) (res *appout.Group, err error) {
	ctx = pkgtrace.StartSpan(ctx, "GroupUsecase.Create")
	defer func() { pkgtrace.EndSpan(ctx, err) }()

	err = in.Validate()
	if err != nil {
		return nil, pkgerror.Errorf("invalid input: %w", err)
	}

	dms := make([]*domain.Member, 0, len(in.MemberNames))
	for _, mn := range in.MemberNames {
		dm, err := domain.NewMember(0, domain.NewMemberArgs{
			Name: mn,
		})
		if err != nil {
			return nil, pkgerror.Errorf("err creating member %s: %w", mn, err)
		}
		dms = append(dms, dm)
	}

	dr, err := domain.NewRule(0, domain.NewRuleArgs{
		MahjongType:           in.Rules.MahjongType,
		InitialPoints:         in.Rules.InitialPoints,
		ReturnPoints:          in.Rules.ReturnPoints,
		RankingPointsFirst:    in.Rules.RankingPointsFirst,
		RankingPointsSecond:   in.Rules.RankingPointsSecond,
		RankingPointsThird:    in.Rules.RankingPointsThird,
		RankingPointsFour:     in.Rules.RankingPointsFour,
		FractionalCalculation: in.Rules.FractionalCalculation,
		UseBust:               in.Rules.UseBust,
		BustPoint:             in.Rules.BustPoint,
		UseChip:               in.Rules.UseChip,
		ChipPoint:             in.Rules.ChipPoint,
	})
	if err != nil {
		return nil, err
	}

	dg, err := domain.NewGroup(domain.NewGroupArgs{
		Name:    in.Name,
		Members: dms,
		Rule:    dr,
	})
	if err != nil {
		return nil, err
	}

	_, err = u.groupRepo.Create(ctx, dg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
