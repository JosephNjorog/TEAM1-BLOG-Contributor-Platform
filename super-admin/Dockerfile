# This is a monorepo workspace (super-admin depends on packages/shared),
# so the build context must be the repo root, not this directory:
#   docker build -f super-admin/Dockerfile -t team1-super-admin .
#
# Optionally pass --build-arg VITE_API_BASE_URL=https://api.example.com
# and --build-arg VITE_FRONTEND_URL=https://app.example.com if this app
# and the backend/main app aren't deployed on the same origin.

FROM node:20-alpine AS build
WORKDIR /repo

COPY package.json package-lock.json ./
COPY packages/shared/package.json packages/shared/package.json
COPY super-admin/package.json super-admin/package.json
RUN npm ci

COPY packages/shared packages/shared
COPY super-admin super-admin

ARG VITE_API_BASE_URL
ARG VITE_FRONTEND_URL
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_FRONTEND_URL=${VITE_FRONTEND_URL}
RUN npm run build -w super-admin

FROM nginx:1.27-alpine
COPY super-admin/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /repo/super-admin/dist /usr/share/nginx/html
EXPOSE 8080
