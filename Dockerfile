# syntax=docker/dockerfile:1

FROM golang:1.20 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /gold-price-getter

FROM alpine:latest AS release-stage
COPY --from=build-stage /gold-price-getter /gold-price-getter
RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
CMD ["/gold-price-getter"]
