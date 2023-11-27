# Use the official Golang image as the base image
FROM --platform=$BUILDPLATFORM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Set the build arguments
ARG TARGETOS TARGETARCH

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o server ./cmd/.

# Use a lightweight base image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

# Command to run the application
CMD [ "./server" ]
