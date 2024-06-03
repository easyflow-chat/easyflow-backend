FROM node:20-alpine as base

RUN npm uninstall -g yarn
RUN npm uninstall -g npm

RUN addgroup -g 2000 -S appgroup
RUN adduser -DH -s /sbin/nologin -u 2000 -G appgroup -S appuser

RUN mkdir /app
RUN chown -R appuser:appgroup /app

WORKDIR /app
COPY --chown=appuser:appgroup ./dist /app/dist
COPY --chown=appuser:appgroup ./node_modules /app/node_modules
COPY --chown=appuse:appgroup ./prisma ./app/prisma
COPY --chown=appuser:appgroup ./entrypoint.sh /app/entrypoint.sh
COPY --chown=appuser:appgroup /.env /app/.env

USER appuser

LABEL org.opencontainers.image.authors="nico.benninger43@gmail.com"
LABEL org.opencontainers.image.source="https://github.com/Dragon437619/easyflow-backend"
LABEL org.opencontainers.image.title="Backend Frontend"
LABEL org.opencontainers.image.description="Backend for Easyflow chat application"

ENV APPLICATION_ROOT="/app"
ENV NODE_ENV="production"

RUN chmod +x ./entrypoint.sh
ENTRYPOINT ./entrypoint.sh
