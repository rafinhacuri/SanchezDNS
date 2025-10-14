FROM golang:1.25.0 AS go-builder

WORKDIR /server

COPY ./api/go.mod ./api/go.sum ./
RUN go mod download

COPY ./api ./

RUN go build -o api .

FROM oven/bun:1-debian

ENV PRODUCTION="true"

RUN apt-get update && apt-get install -y --no-install-recommends curl && rm -rf /var/lib/apt/lists/*

WORKDIR /main

COPY ./package.json ./bun.lock ./.npmrc ./nuxt.config.ts ./tsconfig.json ./
RUN --mount=type=cache,target=/root/.bun/install/cache bun install --ci

COPY ./app ./app
COPY ./public ./public
RUN bun run app:build

COPY --from=go-builder /server/api /main/api

EXPOSE 3000

CMD ["sh", "-c", "/main/api & bun /main/.output/server/index.mjs; wait"]
