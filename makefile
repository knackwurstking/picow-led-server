clean:
	@git clean -f -x -d

test:
	@go test -v --race ./...
	@which deadcode && deadcode ./...

build:
	@rm -rf dist && \
		go mod tidy -v && \
		cd ui && \
		npm install && \
		npm run build && \
		cd .. && \
		go build -v --race -o dist/picow-led-server ./cmd/picow-led-server

dev:
	#DEBUG=nodemon:*,nodemon nodemon -L --signal SIGTERM --exec 'go run ./cmd/picow-led-server -d -c .api.dev.json' --ext '' --delay 3
	nodemon -L --signal SIGTERM --exec 'go run -v --race ./cmd/picow-led-server -d -c .api.dev.json' --ext 'go,mod,sum' --delay 3 --ignore ./ui

