FROM golang:1.24-alpine
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@v1.49.0
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor -o ./tmp/main ./cmd/web/main.go
CMD ["air", "-c", ".air.toml"]