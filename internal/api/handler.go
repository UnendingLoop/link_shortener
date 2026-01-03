package api

import (
	"errors"
	"strconv"
	"time"

	"github.com/UnendingLoop/link_shortener/internal/model"
	"github.com/UnendingLoop/link_shortener/internal/repository"
	"github.com/UnendingLoop/link_shortener/internal/service"

	"github.com/wb-go/wbf/ginext"
)

const layoutDate = "2006-01-02"

type Handler struct {
	svc service.ShortService
}

func NewHandler(srv service.ShortService) *Handler {
	return &Handler{svc: srv}
}

func (h *Handler) SimplePinger(ctx *ginext.Context) {
	ctx.JSON(200, map[string]string{"message": "pong"})
}

func (h *Handler) CreateKey(ctx *ginext.Context) {
	newLink := &model.Link{}

	if err := ctx.BindJSON(newLink); err != nil {
		ctx.JSON(400, map[string]string{"error": "invalid request"})
		return
	}

	newLink, err := h.svc.CreateKey(ctx.Request.Context(), newLink)
	if err != nil {
		ctx.JSON(defineCode(err), map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(201, newLink)
}

func (h *Handler) Redirect(ctx *ginext.Context) {
	shortKey, ok := ctx.Params.Get("short_url")
	if !ok {
		ctx.JSON(400, map[string]string{"error": "invalid shortkey"})
		return
	}
	redir, err := h.svc.GetRedirLinkByKey(ctx.Request.Context(), shortKey, ctx.Request.UserAgent())
	if err != nil {
		ctx.JSON(defineCode(err), map[string]string{"error": err.Error()})
		return
	}

	ctx.Redirect(301, redir)
}

func (h *Handler) GetAnalytics(ctx *ginext.Context) {
	req := &model.AnalyticsRequest{}
	// достаем шортки
	shortKey, ok := ctx.Params.Get("short_url")
	if !ok {
		ctx.JSON(400, map[string]string{"error": "invalid shortkey"})
		return
	}
	req.Shortkey = shortKey

	// достаем тип групинга для аналитики
	req.GroupBy, ok = ctx.GetQuery("grouping")
	if !ok {
		ctx.JSON(400, map[string]string{"error": "invalid grouping"})
		return
	}

	// достаем начальную дату для фильтрации
	startRaw, exists := ctx.GetQuery("start")
	switch exists {
	case true:
		start, err := time.Parse(layoutDate, startRaw)
		if err != nil {
			ctx.JSON(400, map[string]string{"error": "invalid start-date"})
			return
		}
		req.Start = &start
	case false:
		req.Start = nil
	}

	// достаем конечную дату для фильтрации
	endRaw, exists := ctx.GetQuery("end")
	switch exists {
	case true:
		end, err := time.Parse(layoutDate, endRaw)
		if err != nil {
			ctx.JSON(400, map[string]string{"error": "invalid end-date"})
			return
		}
		req.End = &end
	case false:
		req.End = nil
	}

	// читаем офсет и лимит
	limit := 20
	offset := 0
	if limitRaw, ok := ctx.GetQuery("limit"); ok {
		if n, err := strconv.Atoi(limitRaw); err == nil && n > 0 {
			limit = n
		}
	}
	if offsetRaw, ok := ctx.GetQuery("offset"); ok {
		if n, err := strconv.Atoi(offsetRaw); err == nil && n >= 0 {
			offset = n
		}
	}

	// идем в сервис
	data, err := h.svc.GetAnalytics(ctx.Request.Context(), req, limit, offset)
	if err != nil {
		ctx.JSON(defineCode(err), map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, data)
}

func (h *Handler) GetAll(ctx *ginext.Context) {
	limit := 20
	offset := 0

	if limitRaw, ok := ctx.GetQuery("limit"); ok {
		if n, err := strconv.Atoi(limitRaw); err == nil && n > 0 {
			limit = n
		}
	}

	if offsetRaw, ok := ctx.GetQuery("offset"); ok {
		if n, err := strconv.Atoi(offsetRaw); err == nil && n >= 0 {
			offset = n
		}
	}

	data, err := h.svc.GetAll(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(defineCode(err), map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, data)
}

func defineCode(err error) int {
	switch {
	case errors.Is(err, repository.ErrNotFound) || errors.Is(err, service.ErrKeyNotFound):
		return 404
	case errors.Is(err, service.ErrAnalytics) || errors.Is(err, service.ErrNoRedirLink):
		return 422
	case errors.Is(err, service.ErrBusyKey):
		return 409
	default:
		return 500
	}
}
