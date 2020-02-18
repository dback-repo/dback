FROM golang:1.13.7-alpine3.11 as build
RUN apk update && apk add build-base upx 
COPY go-app/src /go/src
RUN cd /go/src/dback && go build -a -installsuffix cgo -ldflags="-s -w" && upx --brute dback

FROM scratch as unpacked
COPY --from=build /go/src/dback/dback /bin/dback

ENV DOCKER_API_VERSION 1.37
ENTRYPOINT ["/bin/dback"]