#!/bin/bash

VERSION=`cat ./version.json| grep Version | sed 's/"//g' | sed 's/  Version: //g'`
GIT_COMMIT=`git rev-parse HEAD`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
REGISTRY="ghcr.io/antonio-alexander/go-bludgeon"
IMAGE_NAME="employees"
GO_ARCH="amd64"
PLATFORM_ARCH="linux/amd64"

if test -n "$GO_ARCH" 
then
    GO_ARCH=$GO_ARCH
fi

if test -n "$PLATFORM_ARCH" 
then
    PLATFORM_ARCH=$PLATFORM_ARCH
fi

case "$GOOS" in
    "windows") OUTPUT_FILE="employees.exe"
    ;;
    *) OUTPUT_FILE="employees"
    ;;
esac

case "$1" in
    "go") go build -ldflags \
        "-X github.com/antonio-alexander/go-bludgeon/employees/cmd/internal.Version=${VERSION} \
        -X github.com/antonio-alexander/go-bludgeon/employees/cmd/internal.GitCommit=${GIT_COMMIT} \
        -X github.com/antonio-alexander/go-bludgeon/employees/cmd/internal.GitBranch=${GIT_BRANCH}" \
        -o $OUTPUT_FILE
    ;;
    "latest") docker build -f ./cmd/service/Dockerfile . -t ${REGISTRY}/${IMAGE_NAME}:latest --build-arg GIT_COMMIT=${GIT_COMMIT} \
        --build-arg GIT_BRANCH=${GIT_BRANCH} --build-arg PLATFORM=${PLATFORM_ARCH} --build-arg GO_ARCH=${GO_ARCH}
    ;;
    *) docker build -f ./cmd/service/Dockerfile . -t ${REGISTRY}/${IMAGE_NAME}:${GO_ARCH}_${VERSION} --build-arg GIT_COMMIT=${GIT_COMMIT} \
        --build-arg GIT_BRANCH=${GIT_BRANCH} --build-arg PLATFORM=${PLATFORM_ARCH} --build-arg GO_ARCH=${GO_ARCH}
    ;;
esac

docker image prune -f
