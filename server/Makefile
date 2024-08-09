all: build

build:
# Server executable
	CGO_ENABLED=1 cd go && go build -o ../build/server main.go
# Static files
	npx tailwindcss -o build/static/css/index.css -i static/css/index.css
	cp -r static/js build/static
	cp -r templates build

run: build
# Explicitly set database path to avoid issues with future changes to default path
	./build/server -database-path='database.db'

clean:
	rm -rf build
