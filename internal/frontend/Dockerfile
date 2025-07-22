FROM golang:1.24.5-alpine3.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/frontend ./cmd/frontend
COPY internal/frontend ./internal/frontend
COPY pkg ./pkg

COPY internal/frontend/.env /app/.env

WORKDIR /app/cmd/frontend

RUN go build -o /frontend

EXPOSE 8080

CMD ["/frontend"]
