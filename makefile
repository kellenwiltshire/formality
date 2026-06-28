REGISTRY_HOST := docker.io
USERNAME := $(DOCKER_USERNAME)
BUILD_IMAGE := ${REGISTRY_HOST}/formality
GIT_SHA :=$(shell git rev-parse HEAD)
BUILD_TAG := $(if $(BUILD_TAG), $(BUILD_TAG), latest)
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
	"$$DOCKERHUB_TOKEN" | docker login -u $(USERNAME) --password-stdin \

build-image-push:
	docker push $(BUILD_IMAGE):$(GIT_SHA)