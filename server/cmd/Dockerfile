###########################################################
#  bludgeon-server [dockerfile]
# 
# This dockerfile can be used to create a docker image of
# the bludgeon-server hosted in a docker container
# 
##########################################################

# build from alpine container with golang
FROM golang:1.16.2-alpine3.13 AS builder

# set the working directory to cmd, so we can build easily
WORKDIR /build

# set environment
ENV CGO_ENABLED=0

# copy current directory into the gopath
COPY . .

# build the executable
RUN \
    mkdir -p /app && \
    cd ./server/cmd/server-rest && \
    go build -gcflags "all=-N -l" -o /app/bludgeon_server && \
    chmod +x /app/bludgeon_server

# build from a clean container (without golang)
FROM alpine:latest

# set working directory to app
WORKDIR /app

# copy delve and the app
COPY --from=builder /app/bludgeon_server /app/bludgeon_server

# use this to just run the server
ENTRYPOINT /app/bludgeon_server
