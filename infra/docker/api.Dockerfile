FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY apps ./apps

RUN go build -o /out/api ./apps/api

FROM alpine:3.20

WORKDIR /app
ENV PORT=8080
COPY --from=build /out/api /app/api

EXPOSE 8080

CMD ["/app/api"]


