FROM balenalib/raspberrypi3-alpine-golang:1.13 as builder
RUN [ "cross-build-start" ]
WORKDIR /go/src/github.com/tekn0ir/camera
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -o camera -a -ldflags '-extldflags "-static"' .
RUN [ "cross-build-end" ]

FROM balenalib/raspberrypi3-alpine:3.8
WORKDIR /
COPY --from=builder /go/src/github.com/tekn0ir/camera/camera .
CMD ["/camera"]