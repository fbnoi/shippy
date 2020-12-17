FROM alpine:latest

RUN apk add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY consignment-service.run .

CMD ["./consignment-service.run"]