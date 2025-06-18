package router

import (
	"github.com/gin-gonic/gin"
	"gpu-container-service/api/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// RegisterAllRoutes 注册所有路由
func RegisterAllRoutes(r *gin.Engine, client clientset.Interface) {
	v1.RegisterAllRoutes(r, client)
}
