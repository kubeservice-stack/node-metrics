FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache gcc musl-dev libc6-compat build-base libc-dev

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum

#RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY vendor/ vendor/

# Build
RUN go build -o node-metrics ./cmd/node-metrics.go


FROM alpine:latest

WORKDIR /
COPY --from=builder /workspace/node-metrics .

ENTRYPOINT ["/node-metrics"]
