FROM golang:alpine as builder

COPY . $GOPATH/src/github.com/rylio/ytdl/
RUN apk update && apk upgrade && \
    apk add --no-cache git
RUN go get -d github.com/rylio/ytdl/cmd/ytdl/
RUN apk del git
WORKDIR $GOPATH/src/github.com/rylio/ytdl/cmd/ytdl/
RUN go build -o /go/bin/ytdl

FROM alpine
COPY --from=builder /go/bin/ytdl /go/bin/ytdl
ENTRYPOINT ["/go/bin/ytdl"]
