FROM oven/bun:1-debian AS nuxt-builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /main

COPY ./package.json ./bun.lock ./.npmrc ./nuxt.config.ts ./tsconfig.json ./

RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install --ci

COPY ./app ./app
COPY ./public ./public

RUN bun run app:build

FROM golang:1.26.4-bookworm AS go-builder

WORKDIR /server

COPY ./api/go.mod ./api/go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./api ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build \
      -trimpath \
      -ldflags "-s -w" \
      -o /server/api .


FROM oven/bun:1-debian AS runner

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl tini bash && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=caddy:2-alpine /usr/bin/caddy /usr/bin/caddy

WORKDIR /main

COPY --from=nuxt-builder /main/.output ./.output
COPY --from=go-builder /server/api ./api

COPY <<'EOF' /etc/caddy/Caddyfile
{
	servers {
		trusted_proxies static 127.0.0.1/8 ::1/128 172.16.0.0/12 10.0.0.0/8 192.168.0.0/16
	}
}

:80

handle_path /go/* {
	reverse_proxy 127.0.0.1:8080
}

handle {
	reverse_proxy 127.0.0.1:3000
}
EOF

COPY <<'EOF' /entrypoint.sh
#!/usr/bin/env bash
set -euo pipefail

shutdown() {
	trap - TERM INT
	kill -TERM 0 2>/dev/null || true
	wait
}
trap shutdown TERM INT

/main/api &
bun /main/.output/server/index.mjs &
caddy run --config /etc/caddy/Caddyfile --adapter caddyfile &

wait -n
echo "[entrypoint] a process exited — shutting down container" >&2
kill -TERM 0 2>/dev/null || true
wait
exit 1
EOF

RUN chmod +x /entrypoint.sh

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
	CMD curl -fsS http://127.0.0.1:80/go/healthcheck || exit 1

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["/entrypoint.sh"]
