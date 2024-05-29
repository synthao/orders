# Image and tag can be overridden via environment variables.
DOCKER_USERNAME ?= synthao
IMAGE_NAME ?= orders
TAG ?= latest

# Name of the Docker image.
IMAGE := ${DOCKER_USERNAME}/${IMAGE_NAME}:${TAG}

.PHONY: all gen

all: build push

gen:
	@protoc -I proto proto/sso/sso.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

build:
	@echo "Building Docker image ${IMAGE}"
	@docker build -t ${IMAGE} .

push:
	@echo "Pushing Docker image ${IMAGE}"
	@docker push ${IMAGE}