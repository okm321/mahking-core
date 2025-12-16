package domain

import "context"

type Group struct {
	ID   int64  // id
	UID  string // uuid
	Name string // グループ名
}

// GroupRepository 永続化層のインタフェース
type GroupRepository interface {
	List(ctx context.Context) ([]Group, error)
	Create(ctx context.Context, name string) (Group, error)
}
