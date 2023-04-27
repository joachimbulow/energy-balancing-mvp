FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/
COPY main.go ./

COPY . .

RUN go build -o app

CMD ["./app"]