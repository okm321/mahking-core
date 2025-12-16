package application

import (
	"context"

	appin "github.com/okm321/mahking-go/internal/application/in"
	appout "github.com/okm321/mahking-go/internal/application/out"
	"github.com/okm321/mahking-go/internal/domain"
)

type GroupUsecase struct {
	groupRepo domain.GroupRepository
}

type NewGroupUsecaseArgs struct {
	GroupRepo domain.GroupRepository
}

func NewGroupUsecase(args *NewGroupUsecaseArgs) *GroupUsecase {
	return &GroupUsecase{
		groupRepo: args.GroupRepo,
	}
}

func (u *GroupUsecase) List(ctx context.Context) ([]appout.Group, error) {
	groups, err := u.groupRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]appout.Group, 0, len(groups))
	for _, g := range groups {
		res = append(res, appout.NewGroup(g))
	}
	return res, nil
}

func (u *GroupUsecase) Create(ctx context.Context, in appin.CreateGroupCommand) (appout.Group, error) {
	err := in.Validate()
	if err != nil {
		return appout.Group{}, err
	}

	g, err := u.groupRepo.Create(ctx, in.Name)
	if err != nil {
		return appout.Group{}, err
	}

	return appout.NewGroup(g), nil
}
