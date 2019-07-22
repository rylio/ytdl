FROM golang:alpine

#COPY . $GOPATH/src/github.com/rylio/ytdl/
RUN apk update && apk upgrade && \
    apk add --no-cache git
RUN go get github.com/brucewangno1/ytdl/cmd/ytdl/
RUN apk del git
WORKDIR /ytdl/

ENTRYPOINT ["/go/bin/ytdl"]
