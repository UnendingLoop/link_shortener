package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"shortener/internal/model"
)

func (p PostgresRepo) Create(ctx context.Context, n *model.Link) error {
	query := `INSERT INTO links (lid, shortkey, redirect, created_at)
	VALUES (DEFAULT, $1, $2, DEFAULT)`

	_, err := p.db.ExecContext(ctx, query,
		n.ShortKey, n.Redirect,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresRepo) GetLIDByKey(ctx context.Context, key string) (int, error) {
	query := `SELECT lid FROM links WHERE shortkey = $1`
	var lid int
	err := p.db.QueryRowContext(ctx, query,
		key,
	).Scan(&lid)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return -1, ErrNotFound // 404
		default:
			return -1, err

		}
	}
	return lid, nil
}

func (p PostgresRepo) GetAll(ctx context.Context, limit, offset int) ([]model.Link, error) {
	query := `SELECT lid, shortkey, redirect, created_at FROM links ORDER BY lid ASC LIMIT $1 OFFSET $2`
	rows, err := p.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	links := make([]model.Link, 0, limit)
	for rows.Next() {
		var l model.Link
		if err := rows.Scan(&l.LID, &l.ShortKey, &l.Redirect, &l.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return links, nil
}

func (p PostgresRepo) GetGroupByUserAgent(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error) {
	query := strings.Builder{}
	args := []interface{}{lid}
	i := 2

	query.WriteString(`SELECT useragent, COUNT(1) as "count", SUM(COUNT(1)) OVER() AS total
	FROM referrals
	WHERE lid = $1`)

	if start != nil {
		query.WriteString(fmt.Sprintf(` AND created_at >= $%d`, i))
		args = append(args, *start)
		i++
	}

	if end != nil {
		query.WriteString(fmt.Sprintf(` AND created_at <= $%d`, i))
		args = append(args, *end)
		i++
	}

	query.WriteString(fmt.Sprintf(` GROUP BY useragent 
	ORDER BY useragent ASC 
	LIMIT $%d OFFSET $%d`, i, i+1))

	args = append(args, limit, offset)

	rows, err := p.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Item, 0, limit)
	total := 0

	for rows.Next() {
		var v model.Item
		if err := rows.Scan(&v.Line, &v.Count, &total); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &model.AnalyticsResponse{Total: total, List: result, GroupBy: "by useragent"}, nil
}

func (p PostgresRepo) GetGroupByDay(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error) {
	query := strings.Builder{}
	args := []interface{}{lid}
	i := 2

	query.WriteString(`SELECT TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD') AS day, COUNT(1) as "count", SUM(COUNT(1)) OVER() AS total
	FROM referrals
	WHERE lid = $1`)

	if start != nil {
		query.WriteString(fmt.Sprintf(` AND created_at >= $%d`, i))
		args = append(args, *start)
		i++
	}

	if end != nil {
		query.WriteString(fmt.Sprintf(` AND created_at <= $%d`, i))
		args = append(args, *end)
		i++
	}

	query.WriteString(fmt.Sprintf(` GROUP BY day 
	ORDER BY day DESC 
	LIMIT $%d OFFSET $%d`, i, i+1))

	args = append(args, limit, offset)

	rows, err := p.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Item, 0, limit)
	total := 0

	for rows.Next() {
		var v model.Item
		if err := rows.Scan(&v.Line, &v.Count, &total); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &model.AnalyticsResponse{Total: total, List: result, GroupBy: "by day"}, nil
}

func (p PostgresRepo) GetGroupByMonth(ctx context.Context, lid int, start, end *time.Time, limit, offset int) (*model.AnalyticsResponse, error) {
	query := strings.Builder{}
	args := []interface{}{lid}
	i := 2

	query.WriteString(`SELECT TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM') AS month, COUNT(1) as "count", SUM(COUNT(1)) OVER() AS total
	FROM referrals
	WHERE lid = $1`)

	if start != nil {
		query.WriteString(fmt.Sprintf(` AND created_at >= $%d`, i))
		args = append(args, *start)
		i++
	}

	if end != nil {
		query.WriteString(fmt.Sprintf(` AND created_at <= $%d`, i))
		args = append(args, *end)
		i++
	}

	query.WriteString(fmt.Sprintf(` GROUP BY month 
	ORDER BY month DESC 
	LIMIT $%d OFFSET $%d`, i, i+1))

	args = append(args, limit, offset)

	rows, err := p.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Item, 0, limit)
	total := 0

	for rows.Next() {
		var v model.Item
		if err := rows.Scan(&v.Line, &v.Count, &total); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &model.AnalyticsResponse{Total: total, List: result, GroupBy: "by month"}, nil
}

func (p PostgresRepo) AddReferralByKey(ctx context.Context, key string, userAgent string) error {
	query := `INSERT INTO referrals (rid, lid, created_at, useragent) VALUES (DEFAULT, $1, DEFAULT, $2)`
	_, err := p.db.ExecContext(ctx, query, key, userAgent)
	if err != nil {
		return fmt.Errorf("failed to create new referral for shortkey %q: %w", key, err)
	}
	return nil
}

func (p PostgresRepo) GetRedirLinkByKey(ctx context.Context, key string) (string, error) {
	query := `SELECT redirect FROM links WHERE shortkey = $1`
	var redirLink string
	err := p.db.QueryRowContext(ctx, query,
		key,
	).Scan(&redirLink)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrNotFound // 404
		default:
			return "", err

		}
	}
	return redirLink, nil
}
