#  build golang code
FROM golang:1.20.6-alpine3.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN mkdir vibe && \
    mkdir numexa-common && \
    mkdir auth && \
    mkdir monger

# Copy go mod and sum files
COPY vibe/go.mod ./vibe
COPY vibe/go.sum ./vibe
COPY auth/go.mod ./auth
COPY auth/go.sum ./auth
COPY monger/go.mod ./monger
COPY numexa-common/go.mod ./numexa-common
# numexa-common doesnot have go.sum file
# we need auth here too, because auth's pkg/db is used by vibe

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN cd vibe && go mod download
RUN cd numexa-common && go mod download 
RUN cd auth && go mod download
RUN cd monger && go mod download

COPY ./vibe ./vibe
COPY ./auth ./auth
COPY ./numexa-common ./numexa-common
COPY ./monger ./monger

# Build the Go app
RUN cd vibe && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nxa-vibe .

# Path: Dockerfile
#  build final image
FROM alpine:3.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/vibe/nxa-vibe /usr/local/bin/nxa-vibe
COPY ./vibe/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose port 8081 to the outside world
EXPOSE 8082

ENTRYPOINT [ "/entrypoint.sh" ]
CMD [ "nxa-vibe" ]

