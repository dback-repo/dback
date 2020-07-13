FROM golang:1.14.1-alpine3.11 as builder
RUN apk update && apk add build-base=0.5-r1 && apk add --no-cache git ca-certificates && update-ca-certificates
COPY node_modules/restic-linux/restic /bin/restic
COPY src /go/src
ENV CGO_ENABLED=0
RUN cd /go/src/dback && go build -a -installsuffix cgo -ldflags="-s -w" && go build -a -installsuffix cgo -o dback-dev && chmod -R 777 /bin


# dev version is compiled with debug info, and not compressed with UPX
FROM scratch as dev
COPY --from=builder /go/src/dback/dback-dev /bin/dback
COPY --from=builder /bin/restic /bin/restic
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
VOLUME /tmp
ENV DOCKER_API_VERSION 1.37
ENTRYPOINT ["/bin/dback"]

FROM builder as compressor
RUN apk add upx=3.95-r2 && cd /go/src/dback && upx --brute dback

# dev version is compiled without debug info, and also compressed with UPX
FROM scratch as prod
COPY --from=compressor /go/src/dback/dback /bin/dback
COPY --from=builder /bin/restic /bin/restic
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
VOLUME /tmp
ENV DOCKER_API_VERSION 1.37
ENTRYPOINT ["/bin/dback"]