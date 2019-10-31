FROM golang:1.13 as builder

WORKDIR /app

COPY internal ./internal
COPY main.go .
COPY go.mod .

RUN CGO_ENABLED=0 && go build -o app

###

FROM ubuntu:18.04 as runner

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]