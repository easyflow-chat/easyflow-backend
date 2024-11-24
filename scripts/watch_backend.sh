#!/bin/sh

reflex -r '^backend/(.*?).go|backend/.env|lib/(.*?).go$' -s -- sh -c 'cd backend && go run main.go'
