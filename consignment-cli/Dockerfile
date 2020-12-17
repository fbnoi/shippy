FROM alpine:latest

RUN mkdir -p /app
WORKDIR /app

ADD consignment.json /app/consignment.json
ADD cli.run /app/cli.run

CMD ["./cli.run"]