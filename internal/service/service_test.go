package service

import (
	"context"
	"testing"

	"shortener/internal/model"
	"shortener/internal/repository"

	"github.com/stretchr/testify/require"
)

func TestCreateKey_OK(t *testing.T) {
	repo := &repoMock{
		checkFreeFn: func(ctx context.Context, key string) (bool, error) {
			return true, nil
		},
		createFn: func(ctx context.Context, link *model.Link) error {
			return nil
		},
	}

	cache := &cacheMock{
		setFn: func(ctx context.Context, key, val string) error {
			return nil
		},
	}

	svc := NewKeyService(repo, cache)

	link := &model.Link{
		Redirect: "https://example.com",
	}

	res, err := svc.CreateKey(context.Background(), link)

	require.NoError(t, err)
	require.NotEmpty(t, res.ShortKey)
}

func TestCreateKey_BusyCustom(t *testing.T) {
	repo := &repoMock{
		checkFreeFn: func(ctx context.Context, key string) (bool, error) {
			return false, nil
		},
	}

	svc := NewKeyService(repo, nil)

	link := &model.Link{
		ShortKey: "abc",
		Redirect: "https://example.com",
	}

	_, err := svc.CreateKey(context.Background(), link)

	require.ErrorIs(t, err, ErrBusyKey)
}

func TestGetRedir_FromCache(t *testing.T) {
	repo := &repoMock{
		addReferralFn: func(ctx context.Context, key, ua string) error {
			return nil
		},
	}

	cache := &cacheMock{
		getFn: func(ctx context.Context, key string) (string, error) {
			return "https://cached.com", nil
		},
	}

	svc := NewKeyService(repo, cache)

	link, err := svc.GetRedirLinkByKey(context.Background(), "abc", "Chrome")

	require.NoError(t, err)
	require.Equal(t, "https://cached.com", link)
}

func TestGetRedir_FromDB(t *testing.T) {
	repo := &repoMock{
		getRedirFn: func(ctx context.Context, key string) (string, error) {
			return "https://db.com", nil
		},
		addReferralFn: func(ctx context.Context, key, ua string) error {
			return nil
		},
	}

	cache := &cacheMock{
		getFn: func(ctx context.Context, key string) (string, error) {
			return "", repository.ErrNotFound
		},
		setFn: func(ctx context.Context, key, val string) error {
			return nil
		},
	}

	svc := NewKeyService(repo, cache)

	link, err := svc.GetRedirLinkByKey(context.Background(), "abc", "Chrome")

	require.NoError(t, err)
	require.Equal(t, "https://db.com", link)
}
