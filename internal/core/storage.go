package core

import "context"

type AccountStorage interface {
	Transfer(ctx context.Context, from, to uint64, amount uint64) error
}
