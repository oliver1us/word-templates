# Stage 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN go mod tidy
RUN go build -o word-templates-api .

# Stage 2: Runtime
FROM alpine:latest

# Install libreoffice and fontconfig
RUN apk update && \
    apk add --no-cache libreoffice fontconfig ttf-freefont

WORKDIR /app

# Copy custom fonts
COPY fonts /usr/share/fonts/
RUN fc-cache -f -v

# Copy binary from builder
COPY --from=builder /app/word-templates-api .

# Package templates into the image
COPY docs docs/

# Create tmp directory
RUN mkdir -p tmp && chmod 777 tmp

EXPOSE 5000

CMD ["./word-templates-api"]