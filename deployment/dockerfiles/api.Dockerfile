FROM golang:1.26-alpine AS builder

LABEL org.opencontainers.image.authors="Alejandro G <contacto@zhacked.me>"
LABEL org.opencontainers.image.source="https://github.com/zhacked/oxyl"

RUN apk add --update --no-cache git mailcap openssl

WORKDIR /workspace

COPY .git ./

COPY api/go.* ./api/
COPY go.* ./
COPY shared ./shared/

RUN go mod download

COPY api ./api

RUN go generate ./shared/pkg/version/version.go && go build -o /workspace/build/oxyl-api-server ./api/main.go

FROM alpine:latest

# copy curl for healthcheck

COPY --from=builder /workspace/build/oxyl-api-server /oxyl-api-server

EXPOSE 19999

CMD ["./oxyl-api-server"]
