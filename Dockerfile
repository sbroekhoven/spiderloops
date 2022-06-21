# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o /app/spiderloops

EXPOSE 8080

CMD [ "/app/spiderloops" ]