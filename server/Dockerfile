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

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN rm -rf build

RUN ./tailwindcss -i src/css/index.css -o build/static/css/index.css --minify

RUN CGO_ENABLED=1 GOOS=linux go build -o build/server src/main.go

FROM gcc:14.2-bookworm
WORKDIR /app

ENV GIN_MODE=release
ARG BLOG_DB_PATH=/app/database.db
ENV BLOG_DB_PATH=${BLOG_DB_PATH}

COPY --from=build /app/build /app
COPY templates /app/templates
COPY src/js /app/static/js

EXPOSE 8080

ENTRYPOINT ["/app/server"]
