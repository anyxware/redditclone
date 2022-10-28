# syntax=docker/dockerfile:1

FROM golang

WORKDIR /usr/src/app

COPY go.mod /usr/src/app
COPY go.sum /usr/src/app
RUN go mod download

COPY . /usr/src/app

CMD make run

