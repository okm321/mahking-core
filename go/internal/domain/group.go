package domain

import (
	"context"

	pkgerror "github.com/okm321/mahking-go/pkg/error"
)

type Group struct {
	ID  int64  // id
	UID string // uuid
	//govalid:required
	//govalid:maxlength=100
	Name string // グループ名
	//govalid:required
	//govalid:maxitems=10
	Members []*Member // グループメンバー
	//govalid:required
	Rule *Rule // グループに紐づくルール
}

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

	if err = grp.validateRules(); err != nil {
		return nil, err
	}

	return &Group{
		Name:    grp.Name,
		Members: grp.Members,
		Rule:    grp.Rule,
	}, nil
}

func (g *Group) validateRules() error {
	if g.Rule.MahjongType.RequiredMemberCount() > len(g.Members) {
		return pkgerror.NewErrorf(
			"%sは最低%d人のメンバーが必要です。 人数: %d人",
			g.Rule.MahjongType.String(),
			g.Rule.MahjongType.RequiredMemberCount(),
			len(g.Members),
		)
	}

	return nil
}

// GroupRepository 永続化層のインタフェース
type GroupRepository interface {
	List(ctx context.Context) ([]Group, error)
	Create(ctx context.Context, group *Group) (*Group, error)
}
