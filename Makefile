build:
	go build -o bin/kabsa main.go
run: build
	sudo ./bin/kabsa
test:
	go test -v ./...
