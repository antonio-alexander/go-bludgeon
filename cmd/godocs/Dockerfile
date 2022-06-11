#---------------------------------------------------------------------------------------------------
# bludgeon-godoc [Dockerfile]
# 
# Reference: https://stackoverflow.com/questions/63178036/how-to-find-commit-hash-from-within-a-running-docker-image
# commit: git rev-parse HEAD
# 
# https://stackoverflow.com/questions/6245570/how-to-get-the-current-branch-name-in-git
# branch: git rev-parse --abbrev-ref HEAD
# 
#---------------------------------------------------------------------------------------------------

ARG GIT_BRANCH=no_git_branch_provided
ARG GIT_COMMIT=no_git_commit_provided
ARG PLATFORM=linux/amd64
ARG GO_ARCH=amd64
ARG GO_ARM=7
ARG REPO=https://github.com/antonio-alexander/go-bludgeon.git
ARG BRANCH=develop

FROM --platform=${PLATFORM} golang:alpine

ARG GO_ARCH
ARG GO_ARM
ARG GIT_BRANCH
ARG GIT_COMMIT
ARG BRANCH
ARG REPO

RUN \
    apk add git openssh-client wget \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux \
    go install golang.org/x/tools/cmd/godoc@v0.1.10

WORKDIR /go/src/github.com/antonio-alexander

RUN \
    git clone --depth=1 --branch ${BRANCH} ${REPO} \
    && rm -r /go/src/github.com/antonio-alexander/go-bludgeon/.git

WORKDIR /go/src/github.com/antonio-alexander/go-bludgeon

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-bludgeon

EXPOSE 8080

CMD [ "godoc", "-http", ":8080" ]
