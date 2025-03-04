# Build the manager binary
FROM golang:1.22 AS builder

WORKDIR /workspace

# Install libzmq to generate certificate (build takes longer, slightly larger image)
RUN apt-get update && apt-get install -y libsodium-dev libzmq3-dev libczmq-dev

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
# This used to copy specific directories, now we
# want to reproduce the local build environment
# COPY . /workspace

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY cmd cmd/
COPY controllers/ controllers/
COPY pkg pkg/
COPY hack hack
COPY Makefile Makefile

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
# Build, ensuring we use the correct keygen
# Note that the original build command did not work here, so updated to mimic the Makefile logic
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CGO_CFLAGS="-I/usr/include" CGO_LDFLAGS="-L/usr/lib -lstdc++ -lczmq -lzmq" go build -a -o manager main.go
RUN make build-container && chmod +x ./manager

# We can't use distroless https://github.com/GoogleContainerTools/distroless
# now that we need the external libraries
FROM debian:stable-slim
WORKDIR /
COPY --from=builder /workspace/manager /manager
COPY --from=builder /workspace/bin/fluxoperator-gen /usr/bin/fluxoperator-gen
RUN apt-get update && apt-get install -y libsodium-dev libzmq3-dev libczmq-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
USER 65532:65532

ENTRYPOINT ["/manager"]
