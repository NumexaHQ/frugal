#  build golang code
FROM golang:1.20.6-alpine3.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN mkdir auth && \
    mkdir numexa-common

# Copy go mod and sum files
COPY auth/go.mod ./auth
COPY auth/go.sum ./auth
COPY numexa-common/go.mod ./numexa-common
# numexa-common doesnot have go.sum file

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN cd auth && go mod download
RUN cd numexa-common && go mod download 

COPY ./auth ./auth
COPY ./numexa-common ./numexa-common

# Build the Go app
RUN cd auth && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nxa-auth .

# Path: Dockerfile
#  build final image
FROM alpine:3.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy postgresql migration files
COPY ./numexa-common/postgresql  /usr/local/postgresql
COPY ./numexa-common/clickhouse /usr/local/clickhouse
# Install curl
RUN apk add --no-cache curl postgresql15-client openssl

# Install postgresql migration tool
RUN curl -L "https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz" | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

# Copy binary from builder
COPY --from=builder /app/auth/nxa-auth /usr/local/bin/nxa-auth
COPY ./auth/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port 8080 to the outside world
EXPOSE 8080

ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "nxa-auth" ]

