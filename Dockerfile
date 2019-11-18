FROM golang:alpine

COPY . $GOPATH/src/github.com/corny/ytdl/
RUN apk update && apk upgrade && \
    apk add --no-cache git
RUN go get -d github.com/corny/ytdl/cmd/ytdl/
RUN apk del git
WORKDIR $GOPATH/src/github.com/corny/ytdl/cmd/ytdl/
RUN go build -o /go/bin/ytdl
WORKDIR /ytdl/

ENTRYPOINT ["/go/bin/ytdl"]