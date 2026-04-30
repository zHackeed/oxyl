FROM oven/bun:1 AS base
# https://bun.com/docs/guides/ecosystem/docker -- conventional

LABEL org.opencontainers.image.authors="Alejandro G <contacto@zhacked.me>"
LABEL org.opencontainers.image.source="https://github.com/zhacked/oxyl"

WORKDIR /app

FROM base AS deps
COPY websocket/package*.json ./
RUN bun install

FROM base AS builder
COPY --from=deps /app/node_modules node_modules
COPY websocket/ .
RUN bun run build

FROM base AS prod-deps
COPY websocket/package*.json ./
RUN bun install --production

FROM base AS release
COPY --from=prod-deps /app/node_modules node_modules
COPY --from=builder /app/dist dist
COPY websocket/package*.json ./

EXPOSE 19977/tcp

CMD ["node", "dist/index.js"]