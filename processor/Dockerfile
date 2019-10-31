FROM golang:1.13 as builder

WORKDIR /app

COPY internal ./internal
COPY Makefile .
COPY go.mod .

RUN CGO_ENABLED=0 make build_scan

###

FROM ubuntu:18.04 as runner

WORKDIR /app
ENV REST_PORT=8081

COPY --from=builder /app/service_scan .
COPY .deploy/wait-for-it.sh .

EXPOSE 8081
CMD ["./service_scan", "scan", "30", "1539992129"]