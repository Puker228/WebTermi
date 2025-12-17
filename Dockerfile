FROM golang:1.25.5-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/app
RUN go build -v -o /usr/local/bin/app ./cmd/worker

CMD ["app"]
