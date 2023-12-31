# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app

RUN CGO_ENABLED=0 GOOS=linux go build -o ./atabot ./app

CMD ["./atabot"]
