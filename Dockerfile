FROM golang as builder
RUN mkdir /webhook-src
WORKDIR /webhook-src
COPY ./ .
RUN CGO_ENABLED=0 go build -o bin/webhook

FROM alpine:3.8
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /app
WORKDIR /app
COPY --from=builder /webhook-src/bin/webhook ./

CMD ["./webhook"]