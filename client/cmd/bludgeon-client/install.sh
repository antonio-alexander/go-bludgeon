#!/bin/bash

VERSION=`cat ../../version.json| grep Version | sed 's/"//g' | sed 's/  Version: //g'`
GIT_COMMIT=`git rev-parse HEAD`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`

case "$GOOS" in
    "windows") OUTPUT_FILE="bludgeon-client.exe"
    ;;
    *) OUTPUT_FILE="bludgeon-client"
    ;;
esac

case "$1" in
    "build") go build -ldflags \
        "-X github.com/antonio-alexander/go-bludgeon/client/cmd/internal.Version=${VERSION} \
        -X github.com/antonio-alexander/go-bludgeon/client/cmd/internal.GitCommit=${GIT_COMMIT} \
        -X github.com/antonio-alexander/go-bludgeon/client/cmd/internal.GitBranch=${GIT_BRANCH}" \
        -o $OUTPUT_FILE
    ;;
    *) go install -ldflags \
        "-X github.com/antonio-alexander/go-bludgeon/client/cmd/internal.Version=${VERSION} \
        -X github.com/antonio-alexander/go-bludgeon/client/cmd/internal.GitCommit=${GIT_COMMIT} \
        -X github.com/antonio-alexander/go-bludgeon/client/cmd/internal.GitBranch=${GIT_BRANCH}"
    ;;
esac

