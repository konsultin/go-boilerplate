FROM golang:1.22.5-alpine AS build

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1 GOOS=linux
RUN go build -o /app/main ./app

FROM alpine:latest

ARG PORT=8080

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata libgcc

COPY --from=build /app/main .
COPY config ./config

RUN mkdir -p /app/files

ENV PORT=${PORT}
EXPOSE ${PORT}

ENTRYPOINT ["./main"]
