#---------------------------------------------------------------------------------------------------
# bludgeon-swagger [Dockerfile]
# 
# Reference: https://stackoverflow.com/questions/63178036/how-to-find-commit-hash-from-within-a-running-docker-image
# commit: git rev-parse HEAD
# 
# https://stackoverflow.com/questions/6245570/how-to-get-the-current-branch-name-in-git
# branch: git rev-parse --abbrev-ref HEAD
# 
#---------------------------------------------------------------------------------------------------

ARG GIT_BRANCH
ARG GIT_COMMIT=no_git_commit_provided
ARG GIT_REPO=https://github.com/antonio-alexander/go-bludgeon.git
ARG PLATFORM=linux/amd64
ARG GO_ARCH=amd64
ARG GO_ARM=7

FROM --platform=${PLATFORM} golang:alpine AS build

ARG GO_ARCH
ARG GO_ARM
ARG GIT_BRANCH
ARG GIT_REPO

RUN \
    apk add git openssh-client gcc g++ wget \
    && env GOARCH=${GO_ARCH} GOARM=${GO_ARM} GOOS=linux CGO_ENABLED=1 \
    go install github.com/go-swagger/go-swagger/cmd/swagger@v0.29.0

WORKDIR /go/src/github.com/antonio-alexander

RUN git clone --depth=1 --branch ${GIT_BRANCH} ${GIT_REPO}

WORKDIR /swagger

RUN \
    cd /go/src/github.com/antonio-alexander/go-bludgeon/employees/internal/swagger \
    && swagger generate spec -o /swagger/swagger_employees.json --scan-models \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/timers/internal/swagger\
    && swagger generate spec -o /swagger/swagger_timers.json --scan-models \
    && cd /go/src/github.com/antonio-alexander/go-bludgeon/changes/internal/swagger\
    && swagger generate spec -o /swagger/swagger_changes.json --scan-models 

FROM --platform=${PLATFORM} swaggerapi/swagger-ui:v4.12.0

ARG GIT_BRANCH
ARG GIT_COMMIT

COPY --from=build /swagger/swagger_timers.json /usr/share/nginx/html/api/swagger_timers.json
COPY --from=build /swagger/swagger_employees.json /usr/share/nginx/html/api/swagger_employees.json
COPY --from=build /swagger/swagger_changes.json /usr/share/nginx/html/api/swagger_changes.json

ENV URLS="[{ url: \"./api/swagger_employees.json\", name: \"Employees\" }, { url: \"./api/swagger_timers.json\", name: \"Timers\" }, { url: \"./api/swagger_changes.json\", name: \"Changes\" } ]"

RUN apk add wget

LABEL antonio-alexander.git.branch=${GIT_BRANCH}
LABEL antonio-alexander.git.commit=${GIT_COMMIT}
LABEL org.opencontainers.image.source=https://github.com/antonio-alexander/go-bludgeon

