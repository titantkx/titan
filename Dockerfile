FROM golang:1.23-alpine AS builder

ARG TARGETOS TARGETARCH
RUN echo "Building for $TARGETOS/$TARGETARCH"

ENV PACKAGES="curl make git libc-dev bash gcc linux-headers eudev-dev file build-base binutils"
RUN apk add --no-cache $PACKAGES

ENV GOCACHE=/root/.cache/go-build

WORKDIR /go/src/github.com/titantkx/titan

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.9/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep a43cb22bf85e89bea45c9af04a229bc92f70ddd8216fee7db21d349d7579cff6
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.9/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 797a235aefb5f8b2d60ecd0f3a430bab28d913b46b7661c0280dc7195a7e1144
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.9/libwasmvmstatic_darwin.a /lib/libwasmvmstatic_darwin.a
RUN sha256sum /lib/libwasmvmstatic_darwin.a | grep dfd377c760742fe6771345c3a5bf2c888d45caa1826f20bb67db8c7b889ce2ae

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN if [ "$TARGETARCH" = "amd64" ]; then \
    ARCH="x86_64"; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
    ARCH="aarch64"; \
    else \
    echo "Unsupported architecture: $TARGETARCH"  ; exit 1; \
    fi && \
    cp "/lib/libwasmvm_muslc.$ARCH.a" "/lib/libwasmvm.$ARCH.a"

COPY go.mod go.sum ./
RUN go mod download

COPY ./ .

RUN --mount=type=cache,target="/root/.cache/go-build" GOOS=$TARGETOS GOARCH=$TARGETARCH COSMOS_BUILD_OPTIONS="nostrip static" make build

#############################################

FROM alpine:3

# install netcat
RUN apk add --no-cache netcat-openbsd binutils

COPY --from=builder /go/src/github.com/titantkx/titan/build/titand /usr/bin/titand

EXPOSE 26656 26657 1317 9090
ENTRYPOINT [ "titand" ]
