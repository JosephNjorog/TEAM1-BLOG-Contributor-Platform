# This is a monorepo workspace (frontend depends on packages/shared), so
# the build context must be the repo root, not this directory:
#   docker build -f frontend/Dockerfile -t team1-frontend .
#
# Optionally pass --build-arg VITE_API_BASE_URL=https://api.example.com
# and --build-arg VITE_ADMIN_APP_URL=https://admin.example.com if this
# app and the backend/admin app aren't deployed on the same origin.

FROM node:20-alpine AS build
WORKDIR /repo

COPY package.json package-lock.json ./
COPY packages/shared/package.json packages/shared/package.json
COPY frontend/package.json frontend/package.json
RUN npm ci

COPY packages/shared packages/shared
COPY frontend frontend

ARG VITE_API_BASE_URL
ARG VITE_ADMIN_APP_URL
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_ADMIN_APP_URL=${VITE_ADMIN_APP_URL}
RUN npm run build -w frontend

FROM nginx:1.27-alpine
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /repo/frontend/dist /usr/share/nginx/html
EXPOSE 8080
