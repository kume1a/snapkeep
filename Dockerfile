FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download

COPY *.go ./
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /snapkeepbin

RUN apk add --no-cache make

EXPOSE 8001

ENV ENVIRONMENT=production
ENV PATH="$PATH:/root/.local/bin"

# ENTRYPOINT ["tail", "-f", "/dev/null"]

CMD ["/snapkeepbin"]