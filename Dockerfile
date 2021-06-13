FROM golang:alpine as builder

ARG API_KEY
ARG API_SECRET

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/build
ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o binance-bot ./cmd

FROM scratch

WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/build/binance-bot /app/

ENV API_KEY=${API_KEY}
ENV API_SECRET=${API_SECRET}

CMD ["./binance-bot"]