# syntax=docker/dockerfile:1

FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /gold-price-getter
CMD ["/gold-price-getter"]
