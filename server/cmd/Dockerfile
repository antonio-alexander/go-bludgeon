#---------------------------------------------------------------------------------------------------
# bludgeon-server [Dockerfile]
# 
# Reference: https://stackoverflow.com/questions/63178036/how-to-find-commit-hash-from-within-a-running-docker-image
# commit: git rev-parse HEAD
# 
# https://stackoverflow.com/questions/6245570/how-to-get-the-current-branch-name-in-git
# branch: git rev-parse --abbrev-ref HEAD
# 
# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
# The GOOS will always be linux, but in case you were curious
# 
# Sample docker build commands:
#  docker build -f ./server/cmd/Dockerfile . -t ghcr.io/antonio-alexander/bludgeon-server:amd64_latest --build-arg GIT_COMMIT=$GITHUB_SHA --build-arg GIT_BRANCH=$GITHUB_REF --build-arg PLATFORM=linux/amd64 --build-arg GO_ARCH=amd64
#  docker build -f ./server/cmd/Dockerfile . -t ghcr.io/antonio-alexander/bludgeon-server:armv7_latest --build-arg GIT_COMMIT=$GITHUB_SHA --build-arg GIT_BRANCH=$GITHUB_REF --build-arg PLATFORM=linux/arm/v7 --build-arg GO_ARCH=arm --build-arg GO_ARM=7
#
#---------------------------------------------------------------------------------------------------

ARG GIT_BRANCH=no_git_branch_provided
ARG GIT_COMMIT=no_git_commit_provided
ARG PLATFORM=linux/amd64
ARG GO_ARCH=amd64
ARG GO_ARM=7

FROM golang:alpine AS build

ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GO_ARCH
ARG GO_ARM

WORKDIR /go/src/go-bludgeon/server

COPY . /go/src/go-bludgeon/server

RUN apk add jq

RUN \
    VERSION=`cat /go/src/go-bludgeon/server/version.json | jq '.Version' | sed 's/"//g'` \
    && cd cmd \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux go build -ldflags "-X github.com/antonio-alexander/go-bludgeon/server/cmd/internal.Version=${VERSION} -X github.com/antonio-alexander/go-bludgeon/server/cmd/internal.GitCommit=${GIT_COMMIT} -X github.com/antonio-alexander/go-bludgeon/server/cmd/internal.GitBranch=${GIT_BRANCH}" -o bludgeon-server \
    && chmod +X bludgeon-server \
    && tar -czvf bludgeon-server.tar.gz ./bludgeon-server

FROM --platform=${PLATFORM} alpine:latest

ARG GIT_BRANCH
ARG GIT_COMMIT

WORKDIR /app

COPY --from=build /go/src/go-bludgeon/server/cmd/bludgeon-server.tar.gz /app/bludgeon-server.tar.gz
COPY --from=build /go/src/go-bludgeon/server/cmd/config/bludgeon_server_config-default.json /app/config/bludgeon_server_config-default.json

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-bludgeon

CMD tar -xzf bludgeon-server.tar.gz && ./bludgeon-server

