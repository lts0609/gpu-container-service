package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ClusterName        string `json:",env=CLUSTER_NAME"`
	SchedulingPolicy   string `json:",env=SCHEDULING_POLICY"`
	GpuTypeLabel       string `json:",env=GPU_TYPE_LABEL"`
	AvaliableNodeLabel string `json:",env=AVALIABLE_NODE_LABEL"`
}
