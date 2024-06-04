#!/bin/sh
node dist/src/migrate.js

NODE_ENV="production" node dist/src/main.js