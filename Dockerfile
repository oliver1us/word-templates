# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN go mod tidy
RUN go build -o word-templates-api .

# Stage 2: Runtime
FROM alpine:latest

# Install libreoffice-writer (slimmed down from full libreoffice) and fontconfig
RUN apk add --no-cache libreoffice-writer fontconfig ttf-freefont

# Create a non-root system group and user
RUN addgroup -S appgroup && adduser -S -D -h /home/appuser -G appgroup -s /bin/false appuser

WORKDIR /app

# Copy custom fonts if they exist in the build context
RUN --mount=type=bind,target=/src \
    mkdir -p /usr/share/fonts && \
    if [ -d /src/fonts ]; then cp -r /src/fonts/. /usr/share/fonts/ 2>/dev/null || true; fi && \
    fc-cache -f -v

# Copy binary from builder
COPY --from=builder /app/word-templates-api .

# Prepare directories, copy templates, and set correct permissions for non-root user
RUN --mount=type=bind,target=/src \
    mkdir -p docs tmp && \
    if [ -d /src/docs ]; then cp -r /src/docs/. docs/ 2>/dev/null || true; fi && \
    chown -R appuser:appgroup /app && \
    chmod 755 docs tmp

# Run as non-root user
USER appuser

EXPOSE 5000

CMD ["./word-templates-api"]