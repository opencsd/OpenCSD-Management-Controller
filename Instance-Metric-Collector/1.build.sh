#!/bin/bash
docker_id="ketidevit2"
controller_name="instance-metric-collector"
tag="v1.0"

export GO111MODULE=on
go mod vendor

go build -o build/bin/$controller_name -gcflags all=-trimpath=`pwd` -asmflags all=-trimpath=`pwd` -mod=vendor $controller_name/src/main && \

docker build -t $docker_id/$controller_name:$tag build && \
docker push $docker_id/$controller_name:$tag
