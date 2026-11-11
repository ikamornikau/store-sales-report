FROM node:22-bookworm-slim AS build

WORKDIR /src

ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0

RUN corepack enable && corepack prepare pnpm@12.16.1 --activate

COPY package.json pnpm-lock.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm build

FROM nginxinc/nginx-unprivileged:alpine

COPY --from=build --chown=nginx:nginx /src/dist /usr/share/nginx/html
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

USER nginx

EXPOSE 8080
