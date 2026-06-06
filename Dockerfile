FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod ./
COPY *.go ./

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o word-templates-api .

FROM alpine:latest

RUN apk add --no-cache libreoffice-writer fontconfig ttf-freefont

RUN addgroup -S appgroup && adduser -S -D -h /home/appuser -G appgroup -s /bin/false appuser

WORKDIR /app

RUN --mount=type=bind,target=/src \
    mkdir -p /usr/share/fonts && \
    if [ -d /src/fonts ]; then cp -r /src/fonts/. /usr/share/fonts/ 2>/dev/null || true; fi && \
    fc-cache -f -v

COPY --from=builder /app/word-templates-api .

RUN --mount=type=bind,target=/src \
    mkdir -p docs tmp && \
    if [ -d /src/docs ]; then cp -r /src/docs/. docs/ 2>/dev/null || true; fi && \
    chown -R appuser:appgroup /app && \
    chmod 755 docs tmp

USER appuser

EXPOSE 5000

CMD ["./word-templates-api"]
