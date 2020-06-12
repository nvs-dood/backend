# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.14.4 AS builder
RUN go version

COPY . /go/src/https://github.com/nvs-dood/backend
WORKDIR /go/src/https://github.com/nvs-dood/backend
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

# Stage 2 (to create a downsized "container executable", ~7MB)

# If you need SSL certificates for HTTPS, replace `FROM SCRATCH` with:
#
#   FROM alpine:3.7
#   RUN apk --no-cache add ca-certificates
#
FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/https://github.com/nvs-dood/backend/app .
EXPOSE 3000
ENTRYPOINT ["./app"]