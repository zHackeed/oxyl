FROM golang:1.26-alpine AS builder

LABEL org.opencontainers.image.authors="Alejandro G <contacto@zhacked.me>"
LABEL org.opencontainers.image.source="https://github.com/zhacked/oxyl"

RUN apk add --update --no-cache git mailcap openssl

WORKDIR /workspace

COPY .git ./

COPY service/thresholds/go.* ./service/thresholds/
COPY shared/go.* ./shared/
COPY go.* ./

RUN go work edit -dropuse=./agent -dropuse=./api -dropuse=./ingress -dropuse=./protocol -dropuse=./service/notifications && go mod download

COPY service/thresholds ./service/thresholds
COPY shared ./shared

RUN go generate ./shared/pkg/version/version.go && go build -o /workspace/build/oxyl-threshold-manager ./service/thresholds/main.go

FROM alpine:latest

COPY --from=builder /workspace/build/oxyl-threshold-manager /oxyl-threshold-manager

CMD ["./oxyl-threshold-manager"]
