FROM alpine:latest

WORKDIR /app

COPY bin/flydb-server /app/flydb-server

RUN chmod +x /app/flydb-server

ENTRYPOINT ["/app/flydb-server"]
