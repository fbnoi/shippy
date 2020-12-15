FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

COPY . .

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consignment-service

# Run container
FROM alpine:latest

RUN apk add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/consignment-service .

CMD ["./consignment-service"]