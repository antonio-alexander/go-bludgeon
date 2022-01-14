#!/bin/bash

VERSION=`cat ../../version.json| grep Version | sed 's/"//g' | sed 's/  Version: //g'`
GIT_COMMIT=`git rev-parse HEAD`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
REGISTRY="github.com/antonio-alexander/go-bludgeon"
IMAGE_NAME="bludgeon-server"
GO_ARCH="amd64"
PLATFORM_ARCH="linux/amd64"

case "$GOOS" in
    "windows") OUTPUT_FILE="bludgeon-server.exe"
    ;;
    *) OUTPUT_FILE="bludgeon-server"
    ;;
esac

case "$1" in
    "go") go build -ldflags \
        "-X github.com/antonio-alexander/go-bludgeon/server/cmd/internal.Version=${VERSION} \
        -X github.com/antonio-alexander/go-bludgeon/server/cmd/internal.GitCommit=${GIT_COMMIT} \
        -X github.com/antonio-alexander/go-bludgeon/server/cmd/internal.GitBranch=${GIT_BRANCH}" \
        -o $OUTPUT_FILE
    ;;
    *) docker build -f ./Dockerfile . -t ${REGISTRY}/${IMAGE_NAME}:${GO_ARCH}_${VERSION} --build-arg GIT_COMMIT=${GIT_COMMIT} \
        --build-arg GIT_BRANCH=${GIT_BRANCH} --build-arg PLATFORM=${PLATFORM_ARCH} --build-arg GO_ARCH=${GO_ARCH}
    ;;
esac

