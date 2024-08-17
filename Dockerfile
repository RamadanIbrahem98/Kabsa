FROM golang:1.22-alpine

RUN apk update && apk add --no-cache sqlite-dev gcc g++ musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/kabsa cmd/main.go

CMD ["./bin/kabsa"]
