###########################################################
#  bludgeon-server [dockerfile]
# 
# This dockerfile can be used to create a docker image of
# the bludgeon-server hosted in a docker container
# 
##########################################################
# build from alpine container with golang
FROM golang:1.14.1-alpine3.11 AS builder

# set the working directory to cmd, so we can build easily
WORKDIR /go/src/github.com/antonio-alexander/go-bludgeon/cmd/server-rest

# set environment
ENV GOPATH=/go PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# copy current directory into the gopath
COPY . "$GOPATH"/src/github.com/antonio-alexander/go-bludgeon/

# build/install delve
RUN \
    apk add --no-cache ca-certificates \
    dpkg \
    gcc \
    git \
    musl-dev \
    && mkdir -p "$GOPATH/src" "$GOPATH/bin" \
    && chmod -R 777 "$GOPATH" \
    && go get github.com/go-delve/delve/cmd/dlv 

# build the executable
RUN \
    go get -d ./... && \
    go build -gcflags "all=-N -l"  -o /app/bludgeon_server \
    && chmod +x /app/bludgeon_server

# build from a clean container (without golang)
FROM alpine:latest

# set working directory to app
WORKDIR /app

# copy delve and the app
COPY --from=builder /app/bludgeon_server /app/bludgeon_server
COPY --from=builder /go/bin/dlv /bin/

# use this to just run the server
ENTRYPOINT /app/bludgeon_server

# use this to run delve with the executable
# CMD dlv exec /app/bludgeon_server --headless --listen=":2345" --log --api-version=2 --check-go-version=false
