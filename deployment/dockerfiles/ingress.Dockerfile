FROM golang:1.26-alpine AS builder

LABEL org.opencontainers.image.authors="Alejandro G <contacto@zhacked.me>"
LABEL org.opencontainers.image.source="https://github.com/zhacked/oxyl"

RUN apk add --update --no-cache git mailcap openssl

WORKDIR /workspace

COPY .git ./

COPY ingress/go.* ./ingress/
COPY shared/go.* ./shared/
COPY protocol/go.* ./protocol/
COPY go.* ./

RUN go work edit -dropuse=./agent -dropuse=./api -dropuse=./service/thresholds -dropuse=./service/notifications && go mod download

COPY ingress ./ingress
COPY shared ./shared
COPY protocol ./protocol

RUN go generate ./shared/pkg/version/version.go && go build -o /workspace/build/oxyl-ingress-server ./ingress/main.go

FROM alpine:latest

COPY --from=builder /workspace/build/oxyl-ingress-server /oxyl-ingress-server

EXPOSE 19988

CMD ["./oxyl-ingress-server"]
