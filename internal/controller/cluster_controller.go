package controller

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gpu-container-service/internal/logic"
	"gpu-container-service/internal/svc"
	clientset "k8s.io/client-go/kubernetes"
	"net/http"
	"time"
)

type ClusterController struct {
	svcCtx *svc.GpuContainerServiceContext
	client clientset.Interface
}

func NewClusterController(svcCtx *svc.GpuContainerServiceContext, client clientset.Interface) *ClusterController {
	return &ClusterController{
		svcCtx: svcCtx,
		client: client,
	}
}

func (c *ClusterController) GetClusterResource(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	logx.Infof("Handle Get Cluster Resources Request")
	l := logic.NewGetClusterResourceLogic(r.Context(), c.svcCtx)
	rsp, err := l.Service(ctx, c.client)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, rsp)
	}
}
