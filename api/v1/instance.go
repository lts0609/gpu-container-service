package v1

import (
	"github.com/gin-gonic/gin"
	"gpu-container-service/internal/handler"
	clientset "k8s.io/client-go/kubernetes"
)

// RegisterInstanceCreateRoutes 注册实例创建相关路由
func RegisterInstanceCreateRoutes(r *gin.Engine, client clientset.Interface) {
	r.POST("/gpu-container/instance", handler.InstanceCreateHandler(client))
}
