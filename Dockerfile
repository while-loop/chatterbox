FROM golang:alpine as builder

ENV PATH=${PATH}:${GOPATH}/bin

RUN apk update && apk add git
COPY . /go/src/chatterbox/
WORKDIR /go/src/chatterbox/
RUN go get ./... && go generate ./... && go build

FROM alpine:latest
COPY --from=builder /go/src/chatterbox/chatterbox  /usr/local/bin/chatterbox
ENTRYPOINT ["chatterbox"]