#  build golang code
FROM golang:1.20.6-alpine3.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN mkdir monger && \
    mkdir numexa-common && \
    mkdir auth

# Copy go mod and sum files
COPY monger/go.mod ./monger
COPY monger/go.sum ./monger
COPY numexa-common/go.mod ./numexa-common
# numexa-common doesnot have go.sum file
# we need auth here too, because auth's pkg/db is used by monger
COPY auth/go.mod ./auth
COPY auth/go.sum ./auth

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN cd monger && go mod download
RUN cd numexa-common && go mod download 
RUN cd auth && go mod download

COPY ./monger ./monger
COPY ./auth ./auth
COPY ./numexa-common ./numexa-common

# Build the Go app
RUN cd monger && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nxa-monger .

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

