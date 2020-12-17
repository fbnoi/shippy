FROM alpine:latest

RUN apk add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY  vessel-service.run .

CMD ["./vessel-service.run"]