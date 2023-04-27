# Build the manager binary
FROM golang:1.18 as builder

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
COPY controllers/ controllers/
COPY pkg pkg/
COPY hack hack
COPY Makefile Makefile

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
# Build, ensuring we use the correct keygen
# Note that the original build command did not work here, so updated to mimic the Makefile logic
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CGO_CFLAGS="-I/usr/include" CGO_LDFLAGS="-L/usr/lib -lstdc++ -lczmq -lzmq" go build -a -o manager main.go
RUN make build-container && chmod +x ./manager
RUN apt-get clean && rm -rf -rf /var/lib/apt/lists/*

# We can't use distroless https://github.com/GoogleContainerTools/distroless
# now that we need the external libraries
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager /manager

# Carefully copy libraries over
COPY --from=builder /usr/lib/x86_64-linux-gnu/libsodium.so /usr/lib/x86_64-linux-gnu/libsodium.so
COPY --from=builder /usr/include/sodium /usr/include/sodium
COPY --from=builder /usr/lib/x86_64-linux-gnu/libsodium.a /usr/lib/x86_64-linux-gnu/libsodium.a
COPY --from=builder /usr/lib/x86_64-linux-gnu/pkgconfig/libsodium.pc /usr/lib/x86_64-linux-gnu/pkgconfig/libsodium.pc

COPY --from=builder /usr/lib/x86_64-linux-gnu/libzmq.so /usr/lib/x86_64-linux-gnu/libzmq.so
COPY --from=builder /usr/include/zmq.h /usr/include/zmq.h
COPY --from=builder /usr/include/zmq.hpp /usr/include/zmq.hpp
COPY --from=builder /usr/include/zmq_addon.hpp /usr/include/zmq_addon.hpp
COPY --from=builder /usr/include/zmq_utils.h /usr/include/zmq_utils.h
COPY --from=builder /usr/lib/x86_64-linux-gnu/libzmq.a /usr/lib/x86_64-linux-gnu/libzmq.a
COPY --from=builder /usr/lib/x86_64-linux-gnu/pkgconfig/libzmq.pc /usr/lib/x86_64-linux-gnu/pkgconfig/libzmq.pc

COPY --from=builder /usr/include/czmq.h /usr/include/czmq.h
COPY --from=builder /usr/include/czmq_library.h /usr/include/czmq_library.h
COPY --from=builder /usr/include/czmq_prelude.h /usr/include/czmq_prelude.h
COPY --from=builder /usr/include/zactor.h /usr/include/zactor.h
COPY --from=builder /usr/include/zarmour.h /usr/include/zarmour.h
COPY --from=builder /usr/include/zauth.h /usr/include/zauth.h
COPY --from=builder /usr/include/zbeacon.h /usr/include/zbeacon.h
COPY --from=builder /usr/include/zcert.h /usr/include/zcert.h
COPY --from=builder /usr/include/zcertstore.h /usr/include/zcertstore.h
COPY --from=builder /usr/include/zchunk.h /usr/include/zchunk.h
COPY --from=builder /usr/include/zclock.h /usr/include/zclock.h
COPY --from=builder /usr/include/zconfig.h /usr/include/zconfig.h
COPY --from=builder /usr/include/zdigest.h /usr/include/zdigest.h
COPY --from=builder /usr/include/zdir.h /usr/include/zdir.h
COPY --from=builder /usr/include/zdir_patch.h /usr/include/zdir_patch.h
COPY --from=builder /usr/include/zfile.h /usr/include/zfile.h
COPY --from=builder /usr/include/zframe.h /usr/include/zframe.h
COPY --from=builder /usr/include/zgossip.h /usr/include/zgossip.h
COPY --from=builder /usr/include/zhash.h /usr/include/zhash.h
COPY --from=builder /usr/include/zhashx.h /usr/include/zhashx.h
COPY --from=builder /usr/include/ziflist.h /usr/include/ziflist.h
COPY --from=builder /usr/include/zlist.h /usr/include/zlist.h
COPY --from=builder /usr/include/zlistx.h /usr/include/zlistx.h
COPY --from=builder /usr/include/zloop.h /usr/include/zloop.h
COPY --from=builder /usr/include/zmonitor.h /usr/include/zmonitor.h
COPY --from=builder /usr/include/zmsg.h /usr/include/zmsg.h
COPY --from=builder /usr/include/zpoller.h /usr/include/zpoller.h
COPY --from=builder /usr/include/zproxy.h /usr/include/zproxy.h
COPY --from=builder /usr/include/zrex.h /usr/include/zrex.h
COPY --from=builder /usr/include/zsock.h /usr/include/zsock.h
COPY --from=builder /usr/include/zstr.h /usr/include/zstr.h
COPY --from=builder /usr/include/zsys.h /usr/include/zsys.h
COPY --from=builder /usr/include/zuuid.h /usr/include/zuuid.h
COPY --from=builder /usr/lib/x86_64-linux-gnu/libczmq.a /usr/lib/x86_64-linux-gnu/libczmq.a
COPY --from=builder /usr/lib/x86_64-linux-gnu/pkgconfig/libczmq.pc /usr/lib/x86_64-linux-gnu/pkgconfig/libczmq.pc
COPY --from=builder /usr/lib/x86_64-linux-gnu/libczmq.so  /usr/lib/x86_64-linux-gnu/libczmq.so

USER 65532:65532

ENTRYPOINT ["/manager"]