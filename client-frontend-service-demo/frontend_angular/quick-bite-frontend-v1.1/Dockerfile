FROM node:20-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
# Eliminăm declarațiile ENV de aici

COPY --from=build /app/dist/quick-bite-frnt_angular/browser /usr/share/nginx/html

RUN echo '#!/bin/sh' > /docker-entrypoint.sh && \
    echo 'set -e' >> /docker-entrypoint.sh && \
    echo 'envsubst < /usr/share/nginx/html/assets/env.js.template > /usr/share/nginx/html/assets/env.js' >> /docker-entrypoint.sh && \
    echo 'envsubst "\$NGINX_PORT" < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf' >> /docker-entrypoint.sh && \
    echo 'exec "$@"' >> /docker-entrypoint.sh && \
    chmod +x /docker-entrypoint.sh

COPY nginx.conf /etc/nginx/conf.d/default.conf.template
COPY --from=build /app/src/assets/env.js /usr/share/nginx/html/assets/env.js.template

# Expunem un port default, dar poate fi suprascris
EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]