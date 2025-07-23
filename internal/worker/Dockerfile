FROM golang:1.24.5-alpine3.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/worker ./cmd/worker
COPY internal/worker ./internal/worker
COPY pkg ./pkg

WORKDIR /app/cmd/worker

RUN go build -o /worker

CMD ["/worker"]
