FROM golang:1.20-alpine as builder
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev
RUN apk add --no-cache $PACKAGES
WORKDIR /go/src/github.com/titanlab/titan
COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

FROM --platform=linux alpine:3
COPY --from=builder /go/src/github.com/titanlab/titan/build/titand /usr/bin/titand

EXPOSE 26656 26657 1317 9090
ENTRYPOINT [ "titand" ]
