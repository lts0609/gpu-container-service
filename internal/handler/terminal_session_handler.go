package handler

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gpu-container-service/internal/logic"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
)

type TerminalSessionHandler struct {
	svcCtx *svc.CreateInstanceContext
}

func NewTerminalSessionHandler(svcCtx *svc.CreateInstanceContext) *TerminalSessionHandler {
	return &TerminalSessionHandler{
		svcCtx: svcCtx,
	}
}

func (t *TerminalSessionHandler) TerminalSession(w http.ResponseWriter, r *http.Request) {
	var req types.CreateInstanceRequest
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	logx.Infof("Handle New Terminal Session Request")
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewTerminalSessionLogic(r.Context(), t.svcCtx)
	rsp, err := l.Service(ctx, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, rsp)
	}
}
