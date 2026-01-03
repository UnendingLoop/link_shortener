package api

import (
	"context"

	"github.com/UnendingLoop/link_shortener/internal/model"
)

type mockShortService struct {
	createKeyFn    func(ctx context.Context, link *model.Link) (*model.Link, error)
	getRedirFn     func(ctx context.Context, key, ua string) (string, error)
	getAnalyticsFn func(ctx context.Context, req *model.AnalyticsRequest, limit, offset int) (*model.AnalyticsResponse, error)
	getAllFn       func(ctx context.Context, limit, offset int) ([]model.Link, error)
}

func (m *mockShortService) CreateKey(ctx context.Context, link *model.Link) (*model.Link, error) {
	return m.createKeyFn(ctx, link)
}

func (m *mockShortService) GetRedirLinkByKey(ctx context.Context, key, ua string) (string, error) {
	return m.getRedirFn(ctx, key, ua)
}

func (m *mockShortService) GetAnalytics(ctx context.Context, req *model.AnalyticsRequest, limit, offset int) (*model.AnalyticsResponse, error) {
	return m.getAnalyticsFn(ctx, req, limit, offset)
}

func (m *mockShortService) GetAll(ctx context.Context, limit, offset int) ([]model.Link, error) {
	return m.getAllFn(ctx, limit, offset)
}
