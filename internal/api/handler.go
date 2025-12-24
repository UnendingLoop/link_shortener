package api

import (
	"shortener/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	svc *service.ShortService
}

func NewHandler(srv *service.ShortService) *Handler {
	return &Handler{svc: srv}
}

func (H *Handler) SimplePinger(ctx *ginext.Context) {
	ctx.JSON(200, map[string]string{"message": "pong"})
}

func (H *Handler) CreateKey(ctx *ginext.Context) {
	// var req struct {
	// 	Text   string `json:"text"`
	// 	SendAt string `json:"send_at"` // формат "2006-01-02 15:04:05"
	// }

	// if err := ctx.BindJSON(&req); err != nil {
	// 	ctx.JSON(400, map[string]string{"error": "invalid payload"})
	// 	return
	// }

	// sendTime, err := time.Parse("2006-01-02 15:04:05", req.SendAt)
	// if err != nil {
	// 	ctx.JSON(400, map[string]string{"error": "invalid send_at format"})
	// 	return
	// }

	// task, err := H.svc.Create(ctx.Request.Context(), req.Text, sendTime)
	// if err != nil {
	// 	ctx.JSON(errorChecker(err), map[string]string{"error": err.Error()})
	// 	return
	// }

	// ctx.JSON(200, map[string]string{"id": task.ID, "status": task.Status})
}

func (H *Handler) Redirect(ctx *ginext.Context) {
	// reqID, ok := ctx.Params.Get("uid")
	// if !ok {
	// 	ctx.JSON(400, map[string]string{"error": "invalid task uid"})
	// 	return
	// }

	// task, err := H.svc.GetByID(ctx.Request.Context(), reqID)
	// if err != nil {
	// 	ctx.JSON(errorChecker(err), map[string]string{"error": err.Error()})
	// 	return
	// }

	// ctx.JSON(200, task)
}

func (H *Handler) GetAnalytics(ctx *ginext.Context) {
	// tasks, err := H.svc.GetPending(ctx.Request.Context())
	// if err != nil {
	// 	ctx.JSON(errorChecker(err), map[string]string{"error": err.Error()})
	// 	return
	// }

	// ctx.JSON(200, tasks)
}
