FROM golang:1.20-alpine as builder

ARG TARGETOS TARGETARCH
RUN echo "Building for $TARGETOS/$TARGETARCH"

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev file build-base binutils
RUN apk add --no-cache $PACKAGES

ENV GOCACHE=/root/.cache/go-build

WORKDIR /go/src/github.com/tokenize-titan/titan

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 2687afbdae1bc6c7c8b05ae20dfb8ffc7ddc5b4e056697d0f37853dfe294e913
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 465e3a088e96fd009a11bfd234c69fb8a0556967677e54511c084f815cf9ce63
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvmstatic_darwin.a /lib/libwasmvmstatic_darwin.a
RUN sha256sum /lib/libwasmvmstatic_darwin.a | grep e45a274264963969305ab9b38a992dbc4401ae97252c7d59b217740a378cb5f2

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

COPY --from=builder /go/src/github.com/tokenize-titan/titan/build/titand /usr/bin/titand

EXPOSE 26656 26657 1317 9090
ENTRYPOINT [ "titand" ]
