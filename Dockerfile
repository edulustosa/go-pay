FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/go-pay ./cmd/go-pay

CMD ["./bin/go-pay"]
