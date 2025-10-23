FROM golang:1.25.3 AS go-builder

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /server

COPY ./api/go.mod ./api/go.sum ./
RUN go mod download

COPY ./api ./

RUN go build -o api .

FROM oven/bun:1-debian

ENV PRODUCTION="true"

RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates && update-ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /main

COPY ./package.json ./bun.lock ./.npmrc ./nuxt.config.ts ./tsconfig.json ./
RUN --mount=type=cache,target=/root/.bun/install/cache bun install
RUN bun add ofetch

COPY ./app ./app
COPY ./public ./public
RUN bun run app:build
RUN apt-get update && apt-get install -y rsync
RUN mkdir -p .output/server/node_modules && rsync -a --ignore-existing node_modules/ .output/server/node_modules/

COPY --from=go-builder /server/api /main/api

EXPOSE 3000

CMD ["sh", "-c", "/main/api & bun /main/.output/server/index.mjs; wait"]
