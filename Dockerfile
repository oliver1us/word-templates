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

# Copy custom fonts if they exist in the build context
RUN --mount=type=bind,target=/src \
    mkdir -p /usr/share/fonts && \
    if [ -d /src/fonts ] && [ "$(ls -A /src/fonts)" ]; then cp -r /src/fonts/* /usr/share/fonts/; fi && \
    fc-cache -f -v

# Copy binary from builder
COPY --from=builder /app/word-templates-api .

# Package templates into the image if they exist in the build context
RUN --mount=type=bind,target=/src \
    mkdir -p docs tmp && \
    chmod 777 docs tmp && \
    if [ -d /src/docs ] && [ "$(ls -A /src/docs)" ]; then cp -r /src/docs/* docs/; fi

EXPOSE 5000

CMD ["./word-templates-api"]