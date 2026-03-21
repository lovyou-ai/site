.PHONY: generate build run dev deploy

generate:
	templ generate

build: generate
	go build -o site ./cmd/site/

run: build
	./site

dev:
	templ generate --watch &
	go run ./cmd/site/

deploy:
	fly deploy
