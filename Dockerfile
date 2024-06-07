FROM node:20-alpine as base

#Uninstall yarn
RUN npm uninstall -g yarn

RUN addgroup -g 2000 -S appgroup
RUN adduser -DH -s /sbin/nologin -u 2000 -G appgroup -S appuser

RUN mkdir /app
RUN chown -R appuser:appgroup /app

WORKDIR /app
COPY --chown=appuser:appgroup /prisma /app/prisma
COPY --chown=appuser:appgroup /entrypoint.sh /app/entrypoint.sh

#Get deleted after build
COPY --chown=appuser:appgroup /enums /app/enums
COPY --chown=appuser:appgroup /src /app/src
COPY --chown=appuser:appgroup /.eslintrc.json /app/.eslintrc.json
COPY --chown=appuser:appgroup /.prettierrc /app/.prettierrc
COPY --chown=appuser:appgroup /package-lock.json /app/package-lock.json
COPY --chown=appuser:appgroup /package.json /app/package.json
COPY --chown=appuser:appgroup /tsconfig.build.json /app/tsconfig.build.json
COPY --chown=appuser:appgroup /tsconfig.json /app/tsconfig.json

#Install packages
RUN npm ci

#Lint
RUN npm run lint

#Build
RUN npm run build

#Reinstall production packages
RUN rm -rf node_modules
RUN npm ci --omit=dev

#Generate @prisma/client
RUN npm run prisma:generate

#For debuging purposes
RUN ls -la


#Romve build dependencies
RUN rm -rf /enums
RUN rm -rf /src
RUN rm -rf /.eslintrc.json
RUN rm -rf /.prettierrc
RUN rm -rf /package-lock.json
RUN rm -rf /package.json
RUN rm -rf /tsconfig.build.json
RUN rm -rf /tsconfig.json


#Uninstall npm
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