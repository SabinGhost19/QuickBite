FROM nginx:alpine

WORKDIR /usr/share/nginx/html

# Copy static files
COPY index.html .
COPY styles.css .
COPY script.js .

# Expose port
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]