FROM golang:1.14.1-alpine3.11 as builder
RUN apk update && apk add build-base=0.5-r1
COPY go-app/src /go/src
ENV CGO_ENABLED=0
RUN (cd /go/src/dback && go build -a -installsuffix cgo -ldflags="-s -w") & (cd /go/src/dback && go build -a -installsuffix cgo -o dback-dev) ; wait

# dev version is compiled with debug info, and not compressed with UPX
FROM scratch as dev
COPY --from=builder /go/src/dback/dback-dev /bin/dback
COPY node_modules/restic-linux/restic /bin/restic
ENV DOCKER_API_VERSION 1.37
ENTRYPOINT ["/bin/dback"]

FROM builder as compressor
RUN apk add upx=3.95-r2 && cd /go/src/dback && upx --brute dback

# dev version is compiled without debug info, and also compressed with UPX
FROM scratch as prod
COPY --from=compressor /go/src/dback/dback /bin/dback
ENV DOCKER_API_VERSION 1.37
ENTRYPOINT ["/bin/dback"]