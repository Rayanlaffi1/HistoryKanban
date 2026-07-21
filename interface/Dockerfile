FROM node:22-alpine AS construction
WORKDIR /application
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
ARG VITE_URLAPI
ARG VITE_URLWS
ARG VITE_KEYCLOAKURL
ARG VITE_KEYCLOAKROYAUME
ARG VITE_KEYCLOAKCLIENT
ENV VITE_URLAPI=$VITE_URLAPI \
    VITE_URLWS=$VITE_URLWS \
    VITE_KEYCLOAKURL=$VITE_KEYCLOAKURL \
    VITE_KEYCLOAKROYAUME=$VITE_KEYCLOAKROYAUME \
    VITE_KEYCLOAKCLIENT=$VITE_KEYCLOAKCLIENT
RUN npm run build

FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=construction /application/dist /usr/share/nginx/html
EXPOSE 80
