SERVER_IN=src/main.go
SERVER_OUT=server
CSS_IN=src/css/index.css
CSS_OUT=index.css

all: build run

build:
	# Build server
	CGO_ENABLED=0 GOOS=linux go build -o build/${SERVER_OUT} ${SERVER_IN}

	# Bundle HTML/CSS/JS
	cp -r templates build/templates
	npx tailwindcss -i ${CSS_IN} -o build/static/${CSS_OUT} --minify

run: build
	./$SERVER_OUT

docker:
	docker build -t server:latest .

dockerrun:
	docker run server:latest

clean:
	rm -rf build
