# syntax=docker/dockerfile:1

FROM golang:1.22.3 AS build-stage
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN env GOOS=linux CGO_ENABLED=1 CC=gcc go build -o ./campus -ldflags '-s -w' .

FROM ubuntu:latest AS release-stage
RUN apt update && apt -y install bash openssh-client git
WORKDIR /app
RUN mkdir cache
COPY --from=build-stage /app/campus ./campus
COPY *.sh roots.pem ./
RUN chmod +x *.sh

ENTRYPOINT ["./campus"]
CMD ["--db", "--ab"]
