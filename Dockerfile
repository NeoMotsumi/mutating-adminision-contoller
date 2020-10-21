FROM golang:1.11 as builder
RUN mkdir /webhook-src
WORKDIR /webhook-src
COPY ./ .
RUN CGO_ENABLED=0 go build -o bin/webhook

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /webhook-src/bin/webhook ./

CMD ["./webhook"]