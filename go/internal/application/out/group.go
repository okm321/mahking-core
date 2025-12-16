package out

import "github.com/okm321/mahking-go/internal/domain"

// Group is a view model returned by usecases.
type Group struct {
	ID   int64
	UID  string
	Name string
}

func NewGroup(g domain.Group) Group {
	return Group{
		ID:   g.ID,
		UID:  g.UID,
		Name: g.Name,
	}
}
