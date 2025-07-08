package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"net/http"
	"strconv"
)

type GetClusterResourceLogic struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.GpuContainerServiceContext
}

func NewGetClusterResourceLogic(ctx context.Context, svcCtx *svc.GpuContainerServiceContext) *GetClusterResourceLogic {
	return &GetClusterResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClusterResourceLogic) Service(ctx context.Context, client clientset.Interface) (resp *types.Response, err error) {
	selector := fmt.Sprintf("%s=true", l.svcCtx.Config.AvaliableNodeLabel)
	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_RESOURCE_NOT_FOUND),
			Message: "failed to list nodes",
		}, err
	}

	var nodeInfos []types.NodeResourceInfo
	for _, node := range nodes.Items {
		capacity, allocated, allocatable := GetNodeGpuResource(node)
		i := types.NodeResourceInfo{
			Name:        node.Name,
			GPUType:     node.Labels[l.svcCtx.Config.GpuTypeLabel],
			Capacity:    capacity,
			Allocated:   allocated,
			Allocatable: allocatable,
		}
		nodeInfos = append(nodeInfos, i)
	}
	rspData := types.ClusterResourceResponseData{
		ClusterName: l.svcCtx.Config.ClusterName,
		Nodes:       nodeInfos,
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "get cluster resource success",
		Data:    rspData,
	}, nil
}

func GetNodeGpuResource(node v1.Node) (capacity int, allocated int, allocatable int) {
	capacityQuantity, ok := node.Status.Capacity["nvidia.com/gpu"]
	if !ok {
		logx.Errorf("Node %s has no nvidia.com/gpu capacity", node.Name)
		return 0, 0, 0
	}
	capacity, _ = strconv.Atoi(capacityQuantity.String())
	allocatableQuantity, ok := node.Status.Allocatable["nvidia.com/gpu"]
	if !ok {
		logx.Errorf("Node %s has no nvidia.com/gpu allocatable", node.Name)
		return 0, 0, 0
	}
	allocatable, _ = strconv.Atoi(allocatableQuantity.String())
	allocated = capacity - allocatable

	return capacity, allocatable, capacity
}
