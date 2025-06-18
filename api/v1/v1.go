package v1

import (
	"github.com/gin-gonic/gin"
	clientset "k8s.io/client-go/kubernetes"
)

// RegisterAllRoutes 注册所有v1路由
func RegisterAllRoutes(r *gin.Engine, client clientset.Interface) {
	RegisterInstanceCreateRoutes(r, client)
}
