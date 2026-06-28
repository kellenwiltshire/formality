REGISTRY_HOST := kellenwiltshire
BUILD_IMAGE := $(REGISTRY_HOST)/formality
GIT_SHA :=$(shell git rev-parse HEAD)
BUILD_TAG ?= $(GIT_SHA)

.DEFAULT_GOAL := run

run:
	go run main.go

build:
	go build -o ./formality main.go

build-image:
	docker buildx build \
		--platform "linux/amd64" \
		--tag "$(BUILD_IMAGE):$(GIT_SHA)-build" \
		--target "build" \
		.
	docker buildx build \
		--cache-from "$(BUILD_IMAGE):$(GIT_SHA)-build" \
		--platform "linux/amd64" \
		--tag "$(BUILD_IMAGE):$(GIT_SHA)" \
		.

build-image-login:
	echo $(DOCKERHUB_TOKEN) | docker login -u $(DOCKER_USERNAME) --password-stdin

build-image-push: docker-image-login
	docker image tag $(BUILD_IMAGE):$(GIT_SHA) $(BUILD_IMAGE):$(BUILD_TAG)
	docker image push $(BUILD_IMAGE):$(BUILD_TAG)