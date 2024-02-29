# Copyright (c) 2024- CrowdStrike, Inc. All rights reserved.

PROJECT_NAME=cs-streaming-take-home-task
PROTO_IMAGE_NAME=$(PROJECT_NAME)-proto
EVENT_GENERATOR_IMAGE_NAME=$(PROJECT_NAME)-event-generator
API_IMAGE_NAME=$(PROJECT_NAME)-api
GOLANG_IMAGE_NAME=golang:1.21.7-alpine3.19
KAFKA_EVENT_PIPELINE_IMAGE_NAME=$(PROJECT_NAME)-kafka-event-pipeline

HOST_WORKDIR=$(PWD)
CONTAINER_WORKDIR=/home/dsci/$(PROJECT_NAME)
CONTAINER_PROTO_WORKDIR=/opt

export HOST_WORKDIR

###
# Building/updating the protobuf files
# You shouldn't need to run this for the interview take home task
###
docker-build-proto:
	@docker build \
		--build-arg BUILDER_IMAGE_NAME=$(GOLANG_IMAGE_NAME) \
		--tag $(PROTO_IMAGE_NAME) \
		--file Dockerfile.proto \
		.
	@touch $@

sensor_data.pb.go: docker-build-proto
	@docker run \
		--rm \
		--volume $(HOST_WORKDIR):$(CONTAINER_PROTO_WORKDIR) \
		$(PROTO_IMAGE_NAME) \
		protoc --go_out=$(CONTAINER_PROTO_WORKDIR)/protos/ --proto_path=$(CONTAINER_PROTO_WORKDIR)/protos/ sensor_data.proto


##
# Docker build / run commands
# used for building and running this application
##
docker-build-event-generator: Dockerfile.event_generator protos/sensor_data.pb.go
	@docker build \
		--build-arg BUILDER_IMAGE_NAME=$(GOLANG_IMAGE_NAME) \
		--tag $(EVENT_GENERATOR_IMAGE_NAME) \
		--file Dockerfile.event_generator \
		--no-cache \
		--target deploy \
		.
	@touch $@

docker-build-api: Dockerfile.api
	@docker build \
		--build-arg BUILDER_IMAGE_NAME=$(GOLANG_IMAGE_NAME) \
		--tag $(API_IMAGE_NAME) \
		--file Dockerfile.api \
		--no-cache \
		--target deploy \
		.
	@touch $@

docker-build-kafka-event-pipeline: Dockerfile.kafka_event_pipeline
	@docker build \
		--build-arg BUILDER_IMAGE_NAME=$(GOLANG_IMAGE_NAME) \
		--tag $(KAFKA_EVENT_PIPELINE_IMAGE_NAME) \
		--file Dockerfile.kafka_event_pipeline \
		--no-cache \
		--target deploy \
		.
	@touch $@

run: docker-build-kafka-event-pipeline docker-build-event-generator docker-build-api
	@docker-compose up
