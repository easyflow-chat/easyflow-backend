FROM node:20-alpine as final

RUN addgroup -g 2000 -S appgroup
RUN adduser -DH -s /sbin/nologin -u 2000 -G appgroup -S appuser

RUN mkdir /app
RUN chown -R appuser:appgroup /app

WORKDIR /app
COPY --chown=appuser:appgroup /prisma /app/prisma
COPY --chown=appuser:appgroup /entrypoint.sh /app/entrypoint.sh

#Get deleted after build
COPY --chown=appuser:appgroup /src /app/src
COPY --chown=appuser:appgroup /package.json /app/package.json
COPY --chown=appuser:appgroup /package-lock.json /app/package-lock.json
COPY --chown=appuser:appgroup /.npmrc /app/.npmrc
COPY --chown=appuser:appgroup /enums /app/enums
COPY --chown=appuser:appgroup /tsconfig.json /app/tsconfig.json
COPY --chown=appuser:appgroup /tsconfig.build.json /app/tsconfig.build.json

#Build
RUN npm ci
RUN npm run build
RUN npm run prisma:generate
RUN cp /app/node_modules/prisma/*.node prisma
RUN rm -rf node_modules
RUN npm ci --omit=dev --omit=optional

#Migrating db
RUN npm run prisma:migrate

#Romve build dependencies
RUN rm -rf /app/src
RUN rm -rf /app/package.json
RUN rm -rf /app/package-lock.json
RUN rm -rf /app/.npmrc
RUN rm -rf /app/enums
RUN rm -rf /app/tsconfig.build.json
RUN rm -rf /app/tsconfig.json


#Uninstall yarn and npm not needed anymore
RUN npm uninstall -g yarn
RUN npm uninstall -g npm

USER appuser

LABEL org.opencontainers.image.authors="nico.benninger43@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/Dragon437619/easyflow-backend"
LABEL org.opencontainers.image.title="Backend Frontend"
LABEL org.opencontainers.image.description="Backend for Easyflow chat application"

ENV APPLICATION_ROOT="/app"
ENV NODE_ENV="production"

RUN chmod +x ./entrypoint.sh
ENTRYPOINT ./entrypoint.sh