# syntax=docker/dockerfile:1
FROM golang:1.22-bookworm AS build
WORKDIR /app

# Install build deps
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install --no-install-recommends --assume-yes \
    curl \
    gcc \
    build-essential 
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 tailwindcss

# Copy project files
COPY go.mod go.sum ./
COPY src src
COPY templates templates
COPY tailwind.config.js ./

RUN mkdir build && mkdir build/static

# Install Go deps
RUN go mod download

# Generate CSS
RUN ./tailwindcss -i src/css/index.css -o build/static/css/index.css --minify

# Build backend
ENV CGO_ENABLED=1
ENV GOOS=linux
RUN go build -o build/server src/main.go

FROM gcc:14.2-bookworm
WORKDIR /app

ENV GIN_MODE=release
ENV BLOG_DB_PATH=/app/database.db

COPY --from=build /app/build /app
COPY templates /app/templates
COPY src/js /app/static/js

EXPOSE 8080

ENTRYPOINT ["/app/server"]
