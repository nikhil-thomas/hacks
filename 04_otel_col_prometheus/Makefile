# Determine OS and architecture
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
    ARCH := amd64
endif

# Define variables
KIND_VERSION := v0.20.0
KIND_BIN := .bin/kind
CLUSTER_NAME := otel-to-prometheus

.PHONY: install-kind create-cluster delete-cluster clean deploy-otel deploy-prometheus deploy-all

install-kind:
	@mkdir -p .bin
	@if [ ! -f $(KIND_BIN) ]; then \
		echo "Downloading Kind $(KIND_VERSION) for $(OS)/$(ARCH)..."; \
		curl -Lo $(KIND_BIN) https://kind.sigs.k8s.io/dl/$(KIND_VERSION)/kind-$(OS)-$(ARCH); \
		chmod +x $(KIND_BIN); \
		echo "Kind installed successfully in .bin directory"; \
	else \
		echo "Kind $(KIND_VERSION) already exists in .bin directory"; \
	fi

create-cluster: install-kind
	@if ! $(KIND_BIN) get clusters | grep -q $(CLUSTER_NAME); then \
		$(KIND_BIN) create cluster --name $(CLUSTER_NAME); \
		echo "Kind cluster '$(CLUSTER_NAME)' created successfully"; \
	else \
		echo "Kind cluster '$(CLUSTER_NAME)' already exists"; \
	fi

delete-cluster: install-kind
	@$(KIND_BIN) delete cluster --name $(CLUSTER_NAME)
	@echo "Kind cluster '$(CLUSTER_NAME)' deleted successfully"

clean:
	@rm -rf .bin
	@echo ".bin directory deleted successfully"

deploy-otel: create-cluster
	@kubectl apply -f k8s/otel-collector.yaml
	@echo "Splunk OpenTelemetry Collector deployed successfully"

deploy-prometheus: create-cluster
	@kubectl apply -f k8s/prometheus.yaml
	@echo "Prometheus deployed successfully"

deploy-all: deploy-otel deploy-prometheus
	@echo "All components deployed successfully"

remove-all: create-cluster
	@kubectl delete -f k8s/otel-collector.yaml || true
	@kubectl delete -f k8s/prometheus.yaml || true
	@echo "All components removed from the cluster"