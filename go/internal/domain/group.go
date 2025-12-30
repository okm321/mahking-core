package domain

import (
	"context"

	pkgerror "github.com/okm321/mahking-go/pkg/error"
)

type Group struct {
	ID      int64     // id
	UID     string    // uuid
	Name    string    // グループ名
	Members []*Member // グループメンバー
	Rule    *Rule     // グループに紐づくルール
}

const (
	MaxGroupNameLength   = 100
	MaxGroupMembersCount = 10
)

type NewGroupArgs struct {
	Name    string
	Members []*Member
	Rule    *Rule
}

func NewGroup(args NewGroupArgs) (_ *Group, err error) {
	grp := Group{
		Name:    args.Name,
		Members: args.Members,
		Rule:    args.Rule,
	}

	err = grp.Validate()
	if err != nil {
		return nil, err
	}

	return &Group{
		Name:    grp.Name,
		Members: grp.Members,
		Rule:    grp.Rule,
	}, nil
}

func (g *Group) Validate() error {
	if g.Name == "" {
		return pkgerror.NewError("グループ名は必須です")
	}

	if len(g.Name) > MaxGroupNameLength {
		return pkgerror.NewErrorf("グループ名は%d文字以内で入力してください", MaxGroupNameLength)
	}

	if g.Rule.MahjongType.RequiredMemberCount() > len(g.Members) {
		return pkgerror.NewErrorf(
			"%sは最低%d人のメンバーが必要です。 人数: %d人",
			g.Rule.MahjongType.String(),
			g.Rule.MahjongType.RequiredMemberCount(),
			len(g.Members),
		)
	}

	if len(g.Members) > MaxGroupMembersCount {
		return pkgerror.NewErrorf("グループメンバーは最大%d人までです", MaxGroupMembersCount)
	}

	return nil
}

// GroupRepository 永続化層のインタフェース
type GroupRepository interface {
	List(ctx context.Context) ([]Group, error)
	Create(ctx context.Context, group *Group) (*Group, error)
}
