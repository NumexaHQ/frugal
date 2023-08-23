#  build golang code
FROM golang:1.20.6-alpine3.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN cd monger && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nxa-monger .

# Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nxa-auth .

# Path: Dockerfile
#  build final image
FROM alpine:3.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/monger/nxa-monger /usr/local/bin/nxa-monger
COPY ./monger/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port 8081 to the outside world
EXPOSE 8081

ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "nxa-monger" ]

