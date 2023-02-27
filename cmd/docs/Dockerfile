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

ARG GIT_REPO=https://github.com/antonio-alexander/go-bludgeon.git
ARG GIT_BRANCH=main
ARG GIT_COMMIT=no_git_commit_provided
ARG PLATFORM=linux/amd64
ARG GO_ARCH=amd64
ARG GO_ARM=7

FROM --platform=${PLATFORM} golang:alpine AS staging

ARG GO_ARCH
ARG GO_ARM
ARG GIT_BRANCH
ARG GIT_COMMIT
ARG GIT_REPO

RUN \
    apk add git openssh-client gcc g++ wget curl \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux CGO_ENABLED=1 \
    go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0 \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux \
    go install golang.org/x/tools/cmd/godoc@v0.1.10

RUN \
    git clone --depth=1 --branch ${GIT_BRANCH} ${GIT_REPO} /go/src/github.com/antonio-alexander/go-bludgeon \
    && git -c advice.detachedHead=false clone --depth 1 --branch v4.15.5 https://github.com/swagger-api/swagger-ui /tmp/swagger-ui

RUN \
    cd /go/src/github.com/antonio-alexander/go-bludgeon/employees \
    && go mod download \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/timers \
    && go mod download \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/changes \
    && go mod download \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/healthcheck \
    && go mod download \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/internal \
    && go mod download

WORKDIR /godocs

RUN \
    cd /go/src/github.com/antonio-alexander/go-bludgeon \
    && godoc -http=:8080 -v &> /dev/null & \
    pid=$! \
    && while ! curl http://localhost:8080 &> /dev/null; do sleep 10s; done \
    && sleep 10s \
    && wget --recursive --no-verbose \
    --convert-links --page-requisites \
    --adjust-extension --execute=robots=off \
    --exclude-directories="*" \
    --no-parent --no-clobber\
    --directory-prefix=/godocs \
    --no-host-directories \
    http://localhost:8080/pkg &> /dev/null | exit 0 \
    && cp /go/src/github.com/antonio-alexander/go-bludgeon/cmd/docs/godocs/index.html /godocs/index.html

WORKDIR /swagger

RUN \
    cd /go/src/github.com/antonio-alexander/go-bludgeon/employees/internal/swagger \
    && swagger generate spec -o /swagger/swagger-employees.json --scan-models \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/timers/internal/swagger \
    && swagger generate spec -o /swagger/swagger-timers.json --scan-models \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/changes/internal/swagger \
    && swagger generate spec -o /swagger/swagger-changes.json --scan-models \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/healthcheck/internal/swagger \
    && swagger generate spec -o /swagger/swagger-healthcheck.json --scan-models \
    && cp -r /tmp/swagger-ui/dist /swagger/dist \
    && cp /go/src/github.com/antonio-alexander/go-bludgeon/cmd/docs/swagger/index.html /swagger/index.html \
    && cp /go/src/github.com/antonio-alexander/go-bludgeon/cmd/docs/swagger/swagger-initializer.js /swagger/dist/swagger-initializer.js

FROM --platform=${PLATFORM} busybox:1.35

WORKDIR /http

COPY --from=staging /godocs /http/godocs
COPY --from=staging /swagger /http/swagger

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-bludgeon

EXPOSE 8080

CMD ["busybox", "httpd", "-f", "-v", "-p", "8080"]
