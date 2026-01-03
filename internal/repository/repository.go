// Package repository provides access to DB CRUD operations +analytics data
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/UnendingLoop/link_shortener/internal/model"
)

type ShortRepository interface {
	Create(ctx context.Context, n *model.Link) error
	GetLIDByKey(ctx context.Context, key string) (int, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Link, error)
	AddReferralByKey(ctx context.Context, key string, userAgent string) error
	GetGroupByUserAgent(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error)
	GetGroupByDay(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error)
	GetGroupByMonth(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error)
	CheckKeyIsFree(ctx context.Context, key string) (bool, error)
	GetRedirLinkByKey(ctx context.Context, key string) (string, error)
}

var ErrNotFound error = errors.New("specified shorturl doesn't exist")
