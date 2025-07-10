package controller

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gpu-container-service/internal/logic"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
	clientset "k8s.io/client-go/kubernetes"
	"net/http"
	"time"
)

type Controller interface {
	GetAllInstance(http.ResponseWriter, *http.Request)
	GetInstance(http.ResponseWriter, *http.Request)
	CreateInstance(http.ResponseWriter, *http.Request)
	DeleteInstance(http.ResponseWriter, *http.Request)
}

type InstanceController struct {
	svcCtx *svc.GpuContainerServiceContext
	client clientset.Interface
}

func NewInstanceController(svcCtx *svc.GpuContainerServiceContext, client clientset.Interface) *InstanceController {
	return &InstanceController{
		svcCtx: svcCtx,
		client: client,
	}
}

func (c *InstanceController) GetAllInstance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	logx.Infof("Handle Get All Instance Request")
	l := logic.NewGetAllInstanceLogic(r.Context(), c.svcCtx)
	rsp, err := l.Service(ctx)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, rsp)
	}
}

func (c *InstanceController) GetInstance(w http.ResponseWriter, r *http.Request) {
	var req types.GetInstanceRequest
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	logx.Infof("Handle Get All Instance Request")
	if err := httpx.ParsePath(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	logx.Errorf("Handle Get All Instance Request: %v", req)
	l := logic.NewGetInstanceLogic(r.Context(), c.svcCtx)
	rsp, err := l.Service(ctx, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, rsp)
	}
}

func (c *InstanceController) CreateInstance(w http.ResponseWriter, r *http.Request) {
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

func (c *InstanceController) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	var req types.DeleteInstanceRequest
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	logx.Infof("Handle New Delete Instance Request")
	if err := httpx.ParsePath(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	l := logic.NewDeleteInstanceLogic(r.Context(), c.svcCtx)
	rsp, err := l.Service(ctx, &req, c.client)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else {
		httpx.OkJsonCtx(r.Context(), w, rsp)
	}
}
