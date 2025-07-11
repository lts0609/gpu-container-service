package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ClusterName        string `env:"CLUSTER_NAME"`
	SchedulingPolicy   string `env:"SCHEDULING_POLICY"`
	GpuTypeLabel       string `env:"GPU_TYPE_LABEL"`
	AvaliableNodeLabel string `env:"AVALIABLE_NODE_LABEL"`
}
