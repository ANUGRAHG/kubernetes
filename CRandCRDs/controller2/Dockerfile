# Stage 1: build the Go binary
FROM golang:1.21 AS builder

WORKDIR /app

# Cache dependencies before copying source
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a statically linked binary
COPY custom_controller.go .
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o controller .

# Stage 2: minimal runtime image
FROM scratch

COPY --from=builder /app/controller /controller

ENTRYPOINT ["/controller"]