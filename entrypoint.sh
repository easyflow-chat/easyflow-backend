#!/bin/sh
npx prisma migrate deploy
npx prisma generate
NODE_ENV="production" node dist/src/main.js