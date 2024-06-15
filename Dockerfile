# syntax=docker/dockerfile:1

FROM golang:1.22 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /gold-price-watcher

FROM ubuntu:latest AS release-stage
COPY --from=build-stage /gold-price-watcher /gold-price-watcher

RUN apt -y update
RUN apt install -y software-properties-common
RUN add-apt-repository -y ppa:xtradeb/apps
RUN apt -y update
RUN apt install -y chromium

RUN apt install -y ca-certificates
COPY ${data_path}/*.crt /usr/local/share/ca-certificates
RUN update-ca-certificates

#RUN apk add --no-cache tzdata
RUN apt install -y tzdata
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt clean
RUN mkdir /gold-price-watcher-data
CMD ["/gold-price-watcher"]
