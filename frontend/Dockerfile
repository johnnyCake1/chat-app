FROM node:latest as build

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

# Copy the rest of the application's source code
COPY . .

RUN npm run build

# Serve the build using NGINX
FROM nginx:alpine

# Copy the build output to replace the default nginx contents.
COPY --from=build /app/build /usr/share/nginx/html

RUN chmod -R a+r /usr/share/nginx/html && \
    chown -R nginx:nginx /usr/share/nginx/html

EXPOSE 80

# Run nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
