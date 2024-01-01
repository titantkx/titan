FROM golang:1.20 as builder

WORKDIR /go/src/github.com/tokenize-titan/titan
COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build

#############################################

FROM debian:12

ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm.x86_64.so /usr/lib/libwasmvm.x86_64.so

COPY --from=builder /go/src/github.com/tokenize-titan/titan/build/titand /usr/bin/titand

EXPOSE 26656 26657 1317 9090
ENTRYPOINT [ "titand" ]
