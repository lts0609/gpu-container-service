package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ClusterName        string `yaml:"cluster_name"`
	SchedulingPolicy   string `yaml:"scheduling_policy"`
	GpuTypeLabel       string `yaml:"gpu_type_label"`
	AvaliableNodeLabel string `yaml:"avaliable_node_label"`
}
