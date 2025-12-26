package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shortener/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSimplePinger(t *testing.T) {
	h := NewHandler(&mockShortService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.SimplePinger(c)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"message":"pong"}`, w.Body.String())
}

func TestCreateKey_OK(t *testing.T) {
	mockSvc := &mockShortService{
		createKeyFn: func(ctx context.Context, link *model.Link) (*model.Link, error) {
			link.ShortKey = "abc123"
			return link, nil
		},
	}

	h := NewHandler(mockSvc)

	body := `{"redirect":"https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateKey(c)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp model.Link
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "abc123", resp.ShortKey)
}

func TestCreateKey_BadJSON(t *testing.T) {
	h := NewHandler(&mockShortService{})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{bad json"))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateKey(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRedirect_OK(t *testing.T) {
	mockSvc := &mockShortService{
		getRedirFn: func(ctx context.Context, key, ua string) (string, error) {
			require.Equal(t, "abc", key)
			return "https://example.com", nil
		},
	}

	h := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/abc", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "short_url", Value: "abc"}}

	h.Redirect(c)

	require.Equal(t, http.StatusMovedPermanently, w.Code)
	require.Equal(t, "https://example.com", w.Header().Get("Location"))
}

func TestGetAll_OK(t *testing.T) {
	mockSvc := &mockShortService{
		getAllFn: func(ctx context.Context, limit, offset int) ([]model.Link, error) {
			require.Equal(t, 10, limit)
			require.Equal(t, 5, offset)
			return []model.Link{
				{ShortKey: "a"},
				{ShortKey: "b"},
			}, nil
		},
	}

	h := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/?limit=10&offset=5", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.GetAll(c)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetAnalytics_OK(t *testing.T) {
	mockSvc := &mockShortService{
		getAnalyticsFn: func(
			ctx context.Context,
			req *model.AnalyticsRequest,
			limit, offset int,
		) (*model.AnalyticsResponse, error) {
			require.Equal(t, "abc", req.Shortkey)
			require.Equal(t, "day", req.GroupBy)

			return &model.AnalyticsResponse{
				GroupBy: "day",
				Total:   1,
			}, nil
		},
	}

	h := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/?grouping=day", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "short_url", Value: "abc"}}

	h.GetAnalytics(c)

	require.Equal(t, http.StatusOK, w.Code)
}
