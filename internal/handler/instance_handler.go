package handler

import (
	"github.com/gin-gonic/gin"
	"gpu-container-service/internal/model"
	"gpu-container-service/internal/service"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"net/http"
)

// InstanceCreateHandler 处理创建Pod的请求
func InstanceCreateHandler(client clientset.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.DeployCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			HandleError(c, "Error binding request", err, http.StatusBadRequest)
			return
		}

		if err := req.Validate(); err != nil {
			HandleError(c, "Validate request params", err, http.StatusBadRequest)
			return
		}

		creator := service.NewCreator(client)
		deployment, secret, service, err := creator.CreateInstance(c.Request.Context(), req)
		if err != nil {
			HandleError(c, "Create Instance Error", err, http.StatusBadRequest)
			return
		}

		klog.Infof("Create Deployment: %s Secret %s Service %sin Namespace %s Successfully", deployment.Name, secret.Name, service.Name, req.Namespace)

		c.JSON(http.StatusCreated, gin.H{
			"Message":    "All Related Resource Created Successfully",
			"Namespace":  req.Namespace,
			"Deployment": deployment.Name,
			"Secret":     secret.Name,
			"Service":    service.Name,
		})
	}
}
