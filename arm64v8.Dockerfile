FROM balenalib/generic-aarch64-alpine-golang:1.12 as builder
RUN [ "cross-build-start" ]
WORKDIR /go/src/github.com/tekn0ir/camera
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -o camera -a -ldflags '-extldflags "-static"' .
RUN [ "cross-build-end" ]

FROM balenalib/generic-aarch64-alpine:3.8
WORKDIR /
COPY --from=builder /go/src/github.com/tekn0ir/camera/camera .
CMD ["/camera"]
