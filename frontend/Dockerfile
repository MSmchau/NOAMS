ARG NPM_REGISTRY=https://registry.npmmirror.com

FROM node:20-alpine AS builder

ARG NPM_REGISTRY
RUN npm config set registry ${NPM_REGISTRY}

WORKDIR /build

# 上下文为项目根目录，路径相对于 .
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

FROM nginx:1.25-alpine

# Nginx 配置（相对于项目根目录）
COPY nginx/nginx.conf /etc/nginx/nginx.conf

# 构建产物
COPY --from=builder /build/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
