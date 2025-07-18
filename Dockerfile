FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/webserver/main.go

FROM alpine:latest AS security
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

FROM scratch AS runtime

COPY --from=builder /app/main /main
COPY --from=security /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENTRYPOINT ["/main", "start", "--registered-number=25438296000158", "--token=1d52a9b6b78cf07b08586152459a5c90", "--platform-code=5AKVkHqCn", "--dispatcher-zip-code=29161376", "--clickhouse-addr=clickhouse:9000"]
