FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /go/src/auth-service

COPY . .

RUN go mod download

ARG CGO_ENABLED=0 GOOS=linux

RUN go build -o /out/auth-service cmd/auth/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=build /out/auth-service ./

EXPOSE 50016

ENTRYPOINT [ "/app/auth-service" ]