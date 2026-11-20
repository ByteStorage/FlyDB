FROM golang:alpine as go-builder

RUN go env -w CGO_ENABLED=0

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -trimpath -ldflags '-extldflags "-static -fpic" -s -w' -o ./flydb-server cmd/server/cli/flydb-server.go

FROM alpine:latest

WORKDIR /app

COPY --from=go-builder /build/flydb-server /app/flydb-server

RUN chmod +x /app/flydb-server

EXPOSE 8999

ENTRYPOINT ["/app/flydb-server"]
