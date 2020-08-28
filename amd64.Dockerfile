FROM balenalib/intel-nuc-alpine-golang:1.12 as builder
WORKDIR /go/src/github.com/tekn0ir/camera
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -o camera -a -ldflags '-extldflags "-static"' .

FROM balenalib/intel-nuc-alpine:3.8
WORKDIR /
COPY --from=builder /go/src/github.com/tekn0ir/camera/camera .
CMD ["/camera"]
