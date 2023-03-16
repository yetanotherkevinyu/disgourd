.PHONY: build clean format run test

DISCORD_APP=discord_app
DISCORD_APP_LOC=cmd/$(DISCORD_APP)/main.go

build:
	go build -o $(DISCORD_APP) $(DISCORD_APP_LOC)

clean:
	go mod tidy -v

env:
	docker-compose up -d

denv:
	docker-compose down

format:
	go vet ./...
	gofmt -l -w .

run:
	go run $(DISCORD_APP_LOC)

test:
	go test -v ./...