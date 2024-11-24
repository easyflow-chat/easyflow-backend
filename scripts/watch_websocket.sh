#!/bin/sh

reflex -r '^websocket/(.*?).go|websocket/.env|lib/(.*?).go$' -s -- sh -c 'cd websocket && go run main.go'
