FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/
COPY main.go ./

COPY . .

# set environment variables
ENV BROKER_URL=localhost:9092 \
    BROKER=KAFKA

RUN go build -o app

CMD ["./app"]