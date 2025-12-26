package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"shortener/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/dbpg"
)

func setupRepo(t *testing.T) (*PostgresRepo, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := &PostgresRepo{
		db: &dbpg.DB{Master: db}, // см. ниже
	}

	return repo, mock
}

func TestPostgresRepo_Create(t *testing.T) {
	repo, mock := setupRepo(t)

	link := &model.Link{
		ShortKey: "abc123",
		Redirect: "https://example.com",
	}

	mock.ExpectQuery(`INSERT INTO links`).
		WithArgs(link.ShortKey, link.Redirect).
		WillReturnRows(
			sqlmock.NewRows([]string{"created_at"}).
				AddRow(time.Now()),
		)

	err := repo.Create(context.Background(), link)
	require.NoError(t, err)
	require.False(t, link.CreatedAt.IsZero())

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPostgresRepo_GetLIDByKey_NotFound(t *testing.T) {
	repo, mock := setupRepo(t)

	mock.ExpectQuery(`SELECT lid FROM links`).
		WithArgs("nope").
		WillReturnError(sql.ErrNoRows)

	lid, err := repo.GetLIDByKey(context.Background(), "nope")

	require.ErrorIs(t, err, ErrNotFound)
	require.Equal(t, -1, lid)
}

func TestPostgresRepo_GetRedirLinkByKey(t *testing.T) {
	repo, mock := setupRepo(t)

	mock.ExpectQuery(`SELECT redirect FROM links`).
		WithArgs("abc").
		WillReturnRows(
			sqlmock.NewRows([]string{"redirect"}).
				AddRow("https://example.com"),
		)

	redir, err := repo.GetRedirLinkByKey(context.Background(), "abc")

	require.NoError(t, err)
	require.Equal(t, "https://example.com", redir)
}

func TestPostgresRepo_GetAll(t *testing.T) {
	repo, mock := setupRepo(t)

	rows := sqlmock.NewRows([]string{
		"lid", "shortkey", "redirect", "created_at",
	}).
		AddRow(1, "a", "https://a.com", time.Now()).
		AddRow(2, "b", "https://b.com", time.Now())

	mock.ExpectQuery(`SELECT lid, shortkey, redirect, created_at FROM links`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	result, err := repo.GetAll(context.Background(), 10, 0)

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, "a", result[0].ShortKey)
}

func TestPostgresRepo_GetGroupByUserAgent(t *testing.T) {
	repo, mock := setupRepo(t)

	rows := sqlmock.NewRows([]string{
		"useragent", "count", "total",
	}).
		AddRow("Chrome/Windows", 3, 4).
		AddRow("Edge/Windows", 1, 4)

	mock.ExpectQuery(`FROM referrals`).
		WithArgs(1, 10, 0).
		WillReturnRows(rows)

	resp, err := repo.GetGroupByUserAgent(
		context.Background(),
		1,
		nil,
		nil,
		10,
		0,
	)

	require.NoError(t, err)
	require.Equal(t, 4, resp.Total)
	require.Len(t, resp.List, 2)
	require.Equal(t, "Chrome/Windows", resp.List[0].Line)
}
