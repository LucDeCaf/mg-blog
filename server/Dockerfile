# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /app

# Install npm for tailwindcss CSS generation
RUN apk update
RUN apk add nodejs npm
RUN npm install --global tailwindcss@latest

# Ensure Go builds to the correct target
ENV CGO_ENABLED=0
ENV GOOS=linux

# Install deps
COPY go.mod go.sum ./
RUN go mod download

# Build server
COPY src src
RUN go build -o build/server src/main.go

# Generate CSS
COPY templates templates
COPY tailwind.config.js ./
RUN npx tailwindcss -i src/css/index.css -o build/static/css/index.css --minify

# Get JS
COPY ./src/js ./build/static/js

FROM scratch
WORKDIR /

EXPOSE 8080

ENV GIN_MODE=release

COPY --from=build /app/build .
COPY templates /templates

CMD ["/server"]
