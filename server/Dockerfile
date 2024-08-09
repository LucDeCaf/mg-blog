# syntax=docker/dockerfile:1
FROM golang:1.22-bookworm AS build
WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install --no-install-recommends --assume-yes \
    curl \
    gcc \
    build-essential 
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 tailwindcss

COPY go/go.mod go/go.sum go/
RUN cd go && go mod download

COPY go go
ENV CGO_ENABLED=1
ENV GOOS=linux
RUN cd go && go build -o /app/build/server main.go

COPY static static
COPY templates templates
COPY tailwind.config.js ./
RUN ./tailwindcss -i static/css/index.css -o build/static/css/index.css --minify

FROM gcc:14.2-bookworm
WORKDIR /app

ENV GIN_MODE=release

ARG BLOG_DB_PATH=/app/database.db
ENV BLOG_DB_PATH=${BLOG_DB_PATH}

COPY --from=build /app/build /app
COPY templates /app/templates
COPY static/js /app/static/js

EXPOSE 8080

ENTRYPOINT ["/app/server"]
