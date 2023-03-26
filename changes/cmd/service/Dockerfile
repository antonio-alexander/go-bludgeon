#---------------------------------------------------------------------------------------------------
# bludgeon-change [Dockerfile]
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
ARG HEALTHCHECK_VERSION="v1.0.3"

FROM golang:alpine AS build

ENV GOPROXY=https://proxy.golang.org,direct

ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GO_ARCH
ARG GO_ARM
ARG HEALTHCHECK_VERSION

RUN env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux CGO_ENABLED=0 \
    go install github.com/antonio-alexander/go-bludgeon/healthcheck/cmd/healthcheck-client@${HEALTHCHECK_VERSION} \
    && mv /go/bin/linux_arm/healthcheck-client /go/bin/healthcheck-client 2>/dev/null || :

WORKDIR /go/src/go-bludgeon/changes

COPY go.mod /go/src/go-bludgeon/changes

RUN go mod download

COPY . /go/src/go-bludgeon/changes

RUN VERSION=`cat ./version.json| grep Version | sed 's/"//g' | sed 's/  Version: //g'` \
    && cd cmd/service \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux go build -ldflags \
    "-X github.com/antonio-alexander/go-bludgeon/changes/cmd/internal.Version=$VERSION \
    -X github.com/antonio-alexander/go-bludgeon/changes/cmd/internal.GitCommit=$GIT_COMMIT \
    -X github.com/antonio-alexander/go-bludgeon/changes/cmd/internal.GitBranch=$GIT_BRANCH" \
    -o bludgeon-changes \
    && chmod +X bludgeon-changes \
    && tar -czvf bludgeon-changes.tar.gz ./bludgeon-changes

FROM --platform=${PLATFORM} alpine:latest

ARG GIT_BRANCH
ARG GIT_COMMIT

ENV BLUDGEON_LOG_PREFIX=bludgeon-changes-service

WORKDIR /app

COPY --from=build /go/src/go-bludgeon/changes/cmd/service/bludgeon-changes.tar.gz /app/bludgeon-changes.tar.gz
COPY --from=build /go/src/go-bludgeon/changes/cmd/service/config/default.json /app/config/default.json
COPY --from=build /go/bin/healthcheck-client /app/healthcheck-client

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-bludgeon

HEALTHCHECK --start-period=10s CMD /app/healthcheck-client --client-type=rest --client-address=localhost --client-port=8080 || exit 1

EXPOSE 8080 2345

CMD tar -xzf bludgeon-changes.tar.gz && ./bludgeon-changes

