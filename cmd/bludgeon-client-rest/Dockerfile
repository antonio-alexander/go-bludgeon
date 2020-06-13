#build from alpine container with golang
FROM golang:alpine AS build
#set the working directory to cmd, so we can build easily
WORKDIR /go/src/github.com/antonio-alexander/go-bludgeon/cmd/bludgeon-client-rest
#set environment
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
#copy current directory into the gopath
COPY . "$GOPATH"/src/github.com/antonio-alexander/go-bludgeon/
#install delve and build the hell-http app
RUN \
    apk add --no-cache ca-certificates \
    dpkg \
    gcc \
    git \
    musl-dev \
    && mkdir -p "$GOPATH/src" "$GOPATH/bin" \
    && chmod -R 777 "$GOPATH" \
    && go get github.com/go-delve/delve/cmd/dlv \
    && go build -gcflags "all=-N -l"  -o /app/bludgeon-client-rest \
    && chmod +x /app/bludgeon_server 
#build from a clean container (without golang)
FROM alpine:latest
#set working directory to app
WORKDIR /app
#copy delve and the app
COPY --from=build /app/bludgeon-client-rest /app/bludgeon-client-rest
COPY --from=build /go/bin/dlv /bin/
#use this to just run the server
CMD bludgeon-client-rest
#use this to run delve with the executable
#CMD dlv exec /app/bludgeon_client --headless --listen=":2345" --log --api-version=2 --check-go-version=false
