package handler

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	clientset "k8s.io/client-go/kubernetes"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gpu-container-service/internal/logic"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
)

type CreateInstanceHandler struct {
	svcCtx *svc.CreateInstanceContext
	client clientset.Interface
}

func NewCreateInstanceHandler(svcCtx *svc.CreateInstanceContext, client clientset.Interface) *CreateInstanceHandler {
	return &CreateInstanceHandler{
		svcCtx: svcCtx,
		client: client,
	}
}

func (c *CreateInstanceHandler) CreateInstance(w http.ResponseWriter, r *http.Request) {
	var req types.CreateInstanceRequest
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	logx.Infof("Handle New Create Instance Request")
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewCreateInstanceLogic(r.Context(), c.svcCtx)
	rsp, err := l.Service(ctx, &req, c.client)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, rsp)
	}
}
