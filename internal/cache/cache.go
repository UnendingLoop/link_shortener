package cache

import (
	"context"
	"shortener/internal/model"
)

type KeyCache interface {
	SetByUID(ctx context.Context, key string, data *model.Link) error
	GetByUID(ctx context.Context, key string) (*model.Link, error)
	DeleteByUID(ctx context.Context, key string) error
}
