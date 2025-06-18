package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gpu-container-service/internal/config"
	"gpu-container-service/internal/router"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	SERVICE_PORT    = "8080"
	KUBECONFIG_PATH = ""
)

func main() {
	klog.InitFlags(nil)
	klog.Errorf("gpu-container-service starting")

	cfg := config.NewConfig(SERVICE_PORT, KUBECONFIG_PATH)
	// 创建Kubernetes客户端
	client, err := config.NewKubernetesClient(cfg.KubeconfigPath)
	if err != nil {
		klog.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.RegisterAllRoutes(r, client)

	srv := &http.Server{
		Addr:    ":" + SERVICE_PORT,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			klog.Fatalf("Failed to start server: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		klog.Fatalf("Server Shutdown Failed: %v", err)
	}
	klog.Infof("Shutting down gracefully")
}
