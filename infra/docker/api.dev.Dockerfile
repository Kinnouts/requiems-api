FROM golang:1.26-alpine

RUN apk add --no-cache git && \
    go install github.com/cosmtrek/air@v1.52.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"

CMD ["air"]
