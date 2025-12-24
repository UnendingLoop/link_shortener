// Package repository provides access to DB CRUD operations +analytics data
package repository

import (
	"context"
	"time"

	"shortener/internal/model"
)

type ShortRepository interface {
	Create(ctx context.Context, n *model.Link) error
	GetByKey(ctx context.Context, key string) (*model.Link, error)
	GetAll(ctx context.Context) ([]*model.Link, error)
	GetByUserAgentByPeriod(ctx context.Context, ua string, start, end time.Time) (int, error)
	ExistsByKey(ctx context.Context, key string) (bool, error)
	UpdateStatusByKey(ctx context.Context, key string) error
	AddReferralByKey(ctx context.Context, key string, userAgent string) error
}
