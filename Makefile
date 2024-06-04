migrate:
	go run cmd/migration/main.go

run:
	go run cmd/dating/main.go

test:
	go test -v ./...

build-image:
	docker build -t datingapp -f scripts/build/Dockerfile .

build-binary:
	go build -o ./dating ./cmd/dating/main.go
	go build -o ./migration ./cmd/dating/main.go
