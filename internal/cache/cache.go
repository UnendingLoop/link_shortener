package cache

import (
	"context"
)

type ShortCache interface {
	SetByShortkey(ctx context.Context, key string, redirect string) error
	GetByShortkey(ctx context.Context, key string) (string, error)
	DeleteByShortkey(ctx context.Context, key string) error
}
