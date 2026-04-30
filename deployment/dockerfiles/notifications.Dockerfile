FROM golang:1.26-alpine AS builder

LABEL org.opencontainers.image.authors="Alejandro G <contacto@zhacked.me>"
LABEL org.opencontainers.image.source="https://github.com/zhacked/oxyl"

RUN apk add --update --no-cache git mailcap openssl

WORKDIR /workspace

COPY .git ./

COPY service/notifications/go.* ./service/notifications/
COPY shared/go.* ./shared/
COPY go.* ./

RUN go work edit -dropuse=./agent -dropuse=./api -dropuse=./ingress -dropuse=./protocol -dropuse=./service/thresholds && go mod download

COPY service/notifications ./service/notifications
COPY shared ./shared

RUN go generate ./shared/pkg/version/version.go && go build -o /workspace/build/oxyl-notification-service ./service/notifications/main.go

FROM alpine:latest

COPY --from=builder /workspace/build/oxyl-notification-service /oxyl-notification-service

CMD ["./oxyl-notification-service"]