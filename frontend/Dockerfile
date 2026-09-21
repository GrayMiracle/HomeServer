FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
RUN apk add --no-cache certbot certbot-nginx
WORKDIR /app
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
RUN echo "0 0 * * * certbot renew --quiet && nginx -s reload" > /etc/crontabs/root
EXPOSE 80
EXPOSE 443
CMD sh -c "crond && nginx -g 'daemon off;'"