package svc

import (
	"gpu-container-service/internal/config"
)

type GpuContainerServiceContext struct {
	Config config.Config
}

func NewGpuContainerServiceContext(c config.Config) *GpuContainerServiceContext {
	return &GpuContainerServiceContext{
		Config: c,
	}
}
