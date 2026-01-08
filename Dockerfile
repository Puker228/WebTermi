FROM golang:1.25.5-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o app ./cmd/app

FROM alpine:3.22
WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/app .

CMD ["./app"]
