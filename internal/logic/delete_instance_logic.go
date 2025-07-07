package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"net/http"
	"time"
)

type DeleteInstanceLogic struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.GpuContainerServiceContext
}

func NewDeleteInstanceLogic(ctx context.Context, svcCtx *svc.GpuContainerServiceContext) *DeleteInstanceLogic {
	return &DeleteInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteInstanceLogic) Service(ctx context.Context, req *types.DeleteInstanceRequest, client clientset.Interface) (resp *types.Response, err error) {
	uuid := req.Uuid
	instance, ok := l.svcCtx.Instances[uuid]
	if !ok {
		return &types.Response{
			Code:    int(types.ERR_RESOURCE_NOT_FOUND),
			Message: "instance not found",
		}, nil
	}

	namespace := "default"
	err = client.CoreV1().Pods(namespace).Delete(ctx, instance.Name, metav1.DeleteOptions{})
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_RESOURCE_DELETE_FAILED),
			Message: "delete instance failed",
		}, nil
	}

	instance.Status = types.StatusTerminating
	instance.DeleteTime = time.Now()
	l.svcCtx.Instances[uuid] = instance
	logx.Infof("delete instance success")

	return &types.Response{
		Code:    http.StatusOK,
		Message: "delete instance success",
	}, nil
}
