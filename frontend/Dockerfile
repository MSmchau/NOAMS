ARG NPM_REGISTRY=https://registry.npmmirror.com

FROM node:20-alpine AS builder

ARG NPM_REGISTRY
RUN npm config set registry ${NPM_REGISTRY}

WORKDIR /build

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:1.25-alpine

# Replace default nginx config
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf

# Copy built frontend
COPY --from=builder /build/dist /usr/share/nginx/html

# Create placeholder favicon if not exists (avoid 404/500)
RUN touch /usr/share/nginx/html/favicon.ico

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
