package domain

import (
	"context"

	pkgerror "github.com/okm321/mahking-go/pkg/error"
)

type Member struct {
	ID      int64
	GroupID int64
	Name    string // 名前
}

type NewMemberArgs struct {
	Name string
}

const MaxMemberNameLength = 10

func (a NewMemberArgs) validate() error {
	if a.Name == "" {
		return pkgerror.NewError("名前は必須です")
	}
	if len(a.Name) > MaxMemberNameLength {
		return pkgerror.NewErrorf("名前は%d文字以内で入力してください: %s", MaxMemberNameLength, a.Name)
	}
	return nil
}

func NewMember(groupID int64, args NewMemberArgs) (_ *Member, err error) {
	if err = args.validate(); err != nil {
		return nil, err
	}

	return &Member{
		GroupID: groupID,
		Name:    args.Name,
	}, nil
}

type MemberRepository interface {
	BatchCreateMembers(ctx context.Context, members []*Member) error
}
