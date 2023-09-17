# syntax=docker/dockerfile:1

FROM golang:1.20 AS build-stage
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
RUN apt install wget -y
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN apt install -y ./google-chrome-stable_current_amd64.deb
RUN rm google-chrome-stable_current_amd64.deb
#RUN apk add --no-cache tzdata
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get clean
CMD ["/price-watcher"]
