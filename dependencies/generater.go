package dependencies

import "context"

type IGenerater interface {
	Take(ctx context.Context, key string, step int64) (int64, error)
	Cursor(ctx context.Context, key string) (int64, error)
	Cursors(ctx context.Context) (map[string]int64, error)
}
