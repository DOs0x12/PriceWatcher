# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /price-watcher

FROM ubuntu:latest AS release-stage
COPY --from=build-stage /price-watcher /price-watcher
ENV TZ=Europe/Moscow \
    DEBIAN_FRONTEND=noninteractive
RUN apt-get -y update
#RUN apk add --no-cache tzdata
RUN apt-get install -y tzdata
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get clean
CMD ["/price-watcher"]
