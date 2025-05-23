FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download

COPY *.go ./
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /sonifybin

# # Install goose
# RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN apk add --no-cache make

EXPOSE 8000

ENV ENVIRONMENT=production
ENV PATH="$PATH:/root/.local/bin"

# ENTRYPOINT ["tail", "-f", "/dev/null"]

CMD ["/bin/bash", "-c", "make migrate-prod;/sonifybin"]