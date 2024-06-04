#!/bin/sh
node migrate.js

NODE_ENV="production" node dist/src/main.js