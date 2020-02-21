FROM node:12-alpine as BUILD
WORKDIR /app

# install deps and the copy the src to speed up docker build
COPY package-lock.json package.json /app/
RUN npm install


COPY tsconfig.json /app/
COPY public /app/public/
COPY src /app/src/
RUN npm run build

# Stage 2 - the production environment
FROM nginxinc/nginx-unprivileged
COPY --from=BUILD /app/build /usr/share/nginx/html
COPY image/start.sh /usr/bin/

USER root
RUN chown nginx:root /usr/share/nginx/html/config.js
USER nginx

EXPOSE 8080

CMD /usr/bin/start.sh
