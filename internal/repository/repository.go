package repository

import (
	"context"
	"shortener/internal/model"
	"time"
)

type ShortRepository interface {
	Create(ctx context.Context, n *model.Link) error
	GetByKey(ctx context.Context, id string) (*model.Link, error)
	GetAll(ctx context.Context) ([]*model.Link, error)
	GetByUserAgentByPeriod(ctx context.Context, ua string, start, end time.Time) (int, error)
}
