#!/bin/bash
export GO111MODULE=on
dir=$( pwd )

#latest golang
#export PATH=$PATH:/usr/local/go/bin && \
# go mod init gpu-metric-collector
# go mod vendor
# go mod tidy

image_name="hpc-metric-collector"
registry="ketidevit2"
version="v1.1"

#binary file
go build -o $dir/../build/_output/bin/$image_name -mod=vendor $dir/../cmd/main.go

# make image
docker build -t $image_name:$version $dir/../build 

# add tag
docker tag $image_name:$version $registry/$image_name:$version

# login
docker login && \

# push image
docker push $registry/$image_name:$version
