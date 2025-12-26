package service

import (
	"context"
	"time"

	"shortener/internal/model"
)

// мок репозитоиря
type repoMock struct {
	createFn            func(ctx context.Context, link *model.Link) error
	checkFreeFn         func(ctx context.Context, key string) (bool, error)
	getAllFn            func(ctx context.Context, limit, offset int) ([]model.Link, error)
	getLIDFn            func(ctx context.Context, key string) (int, error)
	getRedirFn          func(ctx context.Context, key string) (string, error)
	addReferralFn       func(ctx context.Context, key, ua string) error
	getAnalyticsByDay   func(ctx context.Context, lid int, s, e *time.Time, l, o int) (*model.AnalyticsResponse, error)
	getAnalyticsByMonth func(ctx context.Context, lid int, s, e *time.Time, l, o int) (*model.AnalyticsResponse, error)
	getAnalyticsByUA    func(ctx context.Context, lid int, s, e *time.Time, l, o int) (*model.AnalyticsResponse, error)
}

func (m *repoMock) Create(ctx context.Context, link *model.Link) error {
	return m.createFn(ctx, link)
}

func (m *repoMock) CheckKeyIsFree(ctx context.Context, key string) (bool, error) {
	return m.checkFreeFn(ctx, key)
}

func (m *repoMock) GetAll(ctx context.Context, limit, offset int) ([]model.Link, error) {
	return m.getAllFn(ctx, limit, offset)
}

func (m *repoMock) GetLIDByKey(ctx context.Context, key string) (int, error) {
	return m.getLIDFn(ctx, key)
}

func (m *repoMock) GetRedirLinkByKey(ctx context.Context, key string) (string, error) {
	return m.getRedirFn(ctx, key)
}

func (m *repoMock) AddReferralByKey(ctx context.Context, key, ua string) error {
	return m.addReferralFn(ctx, key, ua)
}

func (m *repoMock) GetGroupByUserAgent(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error) {
	return m.getAnalyticsByUA(ctx, lid, start, end, limit, offset)
}

func (m *repoMock) GetGroupByDay(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error) {
	return m.getAnalyticsByDay(ctx, lid, start, end, limit, offset)
}

func (m *repoMock) GetGroupByMonth(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error) {
	return m.getAnalyticsByMonth(ctx, lid, start, end, limit, offset)
}

// мок кеша
type cacheMock struct {
	getFn    func(ctx context.Context, key string) (string, error)
	setFn    func(ctx context.Context, key, val string) error
	deleteFn func(ctx context.Context, key string) error
}

func (m *cacheMock) GetByShortkey(ctx context.Context, key string) (string, error) {
	return m.getFn(ctx, key)
}

func (m *cacheMock) SetByShortkey(ctx context.Context, key, val string) error {
	return m.setFn(ctx, key, val)
}

func (m *cacheMock) DeleteByShortkey(ctx context.Context, key string) error {
	return m.deleteFn(ctx, key)
}
