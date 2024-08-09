# syntax=docker/dockerfile:1
FROM golang:1.22-bookworm AS build
WORKDIR /app

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install --no-install-recommends --assume-yes \
    nodejs \
    npm \
    curl \
    gcc \
    build-essential
RUN npm install tailwindcss@latest

COPY go/go.mod go/go.sum go/
RUN cd go && go mod download

COPY Makefile ./

# Server executable
COPY go go
ENV CGO_ENABLED=1
ENV GOOS=linux
RUN make --always-make go

# CSS (need html files to cull unused styles)
COPY tailwind.config.js ./
COPY templates templates
COPY static static
RUN make css

FROM gcc:14.2-bookworm
WORKDIR /app

ENV GIN_MODE=release

COPY --from=build /app/build /app
COPY templates /app/templates
COPY static/js /app/static/js

EXPOSE 8080

ENTRYPOINT ["/app/server", "-database-path=/app/database.db"]
