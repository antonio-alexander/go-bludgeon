#---------------------------------------------------------------------------------------------------
# bludgeon-healthcheck [Dockerfile]
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
#---------------------------------------------------------------------------------------------------

ARG GIT_BRANCH=no_git_branch_provided
ARG GIT_COMMIT=no_git_commit_provided
ARG PLATFORM=linux/amd64
ARG GO_ARCH=amd64
ARG GO_ARM=7

FROM golang:alpine AS build

ENV GOPROXY=https://proxy.golang.org,direct

ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GO_ARCH
ARG GO_ARM

WORKDIR /go/src/go-bludgeon/healthcheck

COPY . /go/src/go-bludgeon/healthcheck

RUN go mod download
RUN VERSION=`cat ./version.json| grep Version | sed 's/"//g' | sed 's/  Version: //g'` \
    && cd cmd/healthcheck-service \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux go build -ldflags \
    "-X github.com/antonio-alexander/go-bludgeon/healthcheck/cmd/internal.Version=$VERSION \
    -X github.com/antonio-alexander/go-bludgeon/healthcheck/cmd/internal.GitCommit=$GIT_COMMIT \
    -X github.com/antonio-alexander/go-bludgeon/healthcheck/cmd/internal.GitBranch=$GIT_BRANCH" \
    -o bludgeon-healthcheck \
    && chmod +X bludgeon-healthcheck \
    && tar -czvf bludgeon-healthcheck.tar.gz ./bludgeon-healthcheck

FROM --platform=${PLATFORM} alpine:latest

ARG GIT_BRANCH
ARG GIT_COMMIT

WORKDIR /app

COPY --from=build /go/src/go-bludgeon/healthcheck/cmd/healthcheck-service/bludgeon-healthcheck.tar.gz /app/bludgeon-healthcheck.tar.gz

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-bludgeon

EXPOSE 8080 2345

CMD tar -xzf bludgeon-healthcheck.tar.gz && ./bludgeon-healthcheck

