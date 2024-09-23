BIN_DIR := .bin
KIND := $(BIN_DIR)/kind
KIND_CLUSTER_NAME := otel-test

.PHONY: all ensure-bin-dir install-kind create-cluster deploy-prometheus deploy-otel-collector clean


all: ensure-bin-dir install-kind create-cluster deploy-prometheus deploy-otel-collector

ensure-bin-dir:
	mkdir -p $(BIN_DIR)

install-kind: ensure-bin-dir
	@if [ ! -f $(KIND) ]; then \
		echo "Installing kind..."; \
		OS=$(shell uname | tr '[:upper:]' '[:lower:]'); \
		ARCH=$(shell uname -m); \
		if [ "$$ARCH" = "x86_64" ]; then ARCH="amd64"; fi; \
		curl -Lo $(KIND) https://kind.sigs.k8s.io/dl/latest/kind-$$OS-$$ARCH; \
		chmod +x $(KIND); \
	else \
		echo "kind is already installed."; \
	fi

create-cluster: install-kind
	$(KIND) create cluster --name $(KIND_CLUSTER_NAME)

clean:
	$(KIND) delete cluster --name $(KIND_CLUSTER_NAME)
	rm -rf $(BIN_DIR)

deploy-prometheus:
	kubectl apply -f ./k8s/prometheus.yaml

deploy-otel-collector:
	kubectl apply -f ./k8s/otel-collector.yaml
