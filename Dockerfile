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
COPY ./shared ./shared
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
    apt-get install -y --no-install-recommends ca-certificates curl && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /main
 
COPY --from=nuxt-builder /main/.output ./.output
COPY --from=go-builder /server/api ./api

EXPOSE 3000

CMD ["sh", "-c", "/main/api & exec bun /main/.output/server/index.mjs"]
