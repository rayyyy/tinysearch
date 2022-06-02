FROM golang:1.18-buster

RUN mkdir -p /app
COPY . /app
WORKDIR /app