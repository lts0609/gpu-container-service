package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// HandleError 统一错误处理函数
func HandleError(c *gin.Context, reason string, err error, code int) {
	errMsg := fmt.Sprintf("Reason: %v, Error: %v", reason, err)
	klog.Errorf(errMsg)
	c.JSON(code, gin.H{"error": errMsg})
}
