FROM golang:1.24.3-alpine@sha256:ef18ee7117463ac1055f5a370ed18b8750f01589f13ea0b48642f5792b234044 AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download

COPY *.go ./
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /snapkeepbin

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

WORKDIR /app

COPY --from=builder /snapkeepbin /snapkeepbin

EXPOSE 8001

ENV ENVIRONMENT=production
ENV PATH="/root/.local/bin:$PATH"

CMD ["/snapkeepbin"]