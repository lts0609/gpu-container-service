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
		total, used, remain := GetNodeGpuResource(ctx, node, client)
		i := types.NodeResourceInfo{
			Name:    node.Name,
			GPUType: node.Labels[l.svcCtx.Config.GpuTypeLabel],
			Total:   total,
			Used:    used,
			Remain:  remain,
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

func GetNodeGpuResource(ctx context.Context, node v1.Node, client clientset.Interface) (total int, used int, remain int) {
	allocatableQuantity, ok := node.Status.Allocatable["nvidia.com/gpu"]
	if !ok {
		logx.Errorf("Node %s has no nvidia.com/gpu allocatable", node.Name)
		return 0, 0, 0
	}
	total, _ = strconv.Atoi(allocatableQuantity.String())

	used = 0
	pods, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name),
	})
	if err != nil {
		logx.Errorf("Failed to list pods on node %s", node.Name)
		return 0, 0, 0
	}
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if gpuQuantity, ok := container.Resources.Requests[v1.ResourceName("nvidia.com/gpu")]; ok {
				request, _ := strconv.Atoi(gpuQuantity.String())
				used += request
			}
		}
	}
	remain = total - used

	return total, used, remain
}
