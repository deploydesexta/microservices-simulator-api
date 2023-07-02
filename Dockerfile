FROM golang:1.20-alpine AS builder

WORKDIR /simulator

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -ldflags="-s -w" -v -o api ./cmd/api

FROM alpine:latest AS runtime

# Copy the binary to the production image from the builder stage.
COPY --from=builder /simulator/.env.local /simulator/.env
COPY --from=builder /simulator/api /simulator/api

CMD ["/simulator/api"]
