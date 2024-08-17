build:
	go build -o bin/kabsa cmd/main.go
run: build
	sudo bin/kabsa
test:
	go test -v ./...
