all: build

.PHONY: build
build: css go
	cp -r static/js build/static
	cp -r templates build

.PHONY: run
run: build
# Explicitly set database path to avoid issues with future changes to default path
	cd build && ./server -database-path='../database.db'

.PHONY: go
go:
	CGO_ENABLED=1 cd go && go build -o ../build/server main.go

.PHONY: css
css:
	npx tailwindcss -o build/static/css/index.css -i static/css/index.css --minify

.PHONY: clean
clean:
	rm -rf build
