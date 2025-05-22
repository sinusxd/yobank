# Stage 1: Build the React app with Vite
FROM node:18 AS build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# Сборка Vite-приложения (будет папка dist)
RUN NODE_ENV=production npm run build

# Stage 2: Serve the build using Nginx
FROM nginx:alpine

# Копируем собранное приложение из dist
COPY --from=build /app/dist /usr/share/nginx/html

COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
