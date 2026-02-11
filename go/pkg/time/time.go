package time

import (
	"context"
	"time"

	"github.com/newmo-oss/ctxtime"
)

func Now(ctx context.Context) time.Time {
	return ctxtime.Now(ctx).In(time.Local)
}
