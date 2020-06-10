# Multi-stage build setup (https://docs.docker.com/develop/develop-images/multistage-build/)

FROM golang:latest AS builder
RUN go version

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Export necessary port
EXPOSE 3000

# Command to run when starting the container
CMD ["/dist/main"]

EXPOSE 3000
ENTRYPOINT ["./app"]