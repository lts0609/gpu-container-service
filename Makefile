# 项目名称
PROJECT_NAME := gpu-container-service
# 镜像仓库
REGISTRY ?= containercloud-mirror.xaidc.com/containercloud/
# 镜像名称
IMAGE_NAME := $(PROJECT_NAME)
# 镜像标签
IMAGE_TAG ?= 1.0.0
# Dockerfile
DOCKERFILE := Dockerfile
# Helm Chart 目录
HELM_CHART_DIR := charts/$(PROJECT_NAME)
# Helm Release 名称
HELM_RELEASE_NAME := $(PROJECT_NAME)
# 命名空间
NAMESPACE ?= default


# 同步 go 依赖
tidy:
	go mod tidy

# 构建 Docker 镜像
build: tidy
	@echo "build image $(REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG)"
	docker build -t $(REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG) -f $(DOCKERFILE) .

# 推送 Docker 镜像
push: build
	@echo "push image $(REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG)"
	docker push $(REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG)

# 删除 Docker 镜像
clean:
	@echo "clean image $(REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG)"
	docker rmi $(REGISTRY)$(IMAGE_NAME):$(IMAGE_TAG) || true

# 打包 Helm Chart
helm-package:
	@echo "package helm chart to charts/output"
	helm package $(HELM_CHART_DIR) -d charts/output

# Helm部署
helm-deploy:
	@echo "deploy helm chart to namespace (NAMESPACE)"
	helm upgrade --install $(HELM_RELEASE_NAME) $(HELM_CHART_DIR) \
		--namespace $(NAMESPACE) \
		--set image.repository=$(REGISTRY)$(IMAGE_NAME) \
		--set image.tag=$(IMAGE_TAG)

# Helm删除
helm-delete:
	@echo "delete helm chart"
	helm delete $(HELM_RELEASE_NAME) --namespace $(NAMESPACE)

# 显示帮助信息
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build           Build the Docker image"
	@echo "  push            Push the Docker image to the registry"
	@echo "  helm-package    Package the Helm Chart"
	@echo "  helm-deploy     Deploy the Helm Chart to Kubernetes"
	@echo "  helm-delete     Delete the Helm release from Kubernetes"
	@echo "  clean           Remove the Docker image"
	@echo "  help            Show this help message"

.PHONY: build push helm-package helm-deploy helm-delete clean help