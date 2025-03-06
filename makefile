CONFIG_LOCATION=.api.dev.json

clean:
	git clean -f -x -d

init:
	cd ui && npm install && cd .. && go mod tidy -v

build:
	rm -rf dist && \
		cd ui && \
		npm run build && \
		cd .. && \
		go test -v --race ./... && \
		go build -v -o dist/picow-led-server ./cmd/picow-led-server

test:
	go test -v --race ./...
	which deadcode && deadcode ./... || exit 0

dev:
	#DEBUG=nodemon:*,nodemon nodemon -L --signal SIGTERM --exec 'go run ./cmd/picow-led-server -d -c .api.dev.json' --ext '' --delay 3
nodemon -L --signal SIGTERM --exec 'go run --race -v ./cmd/picow-led-server -d -c ${CONFIG_LOCATION}' --ext 'go,mod,sum' --delay 3 --ignore ./ui
