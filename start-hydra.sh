#!/usr/bin/env bash

set -e

docker run -it --rm \
    --name ory_hydra_test \
    -p 4444:4444 \
    -p 4445:4445 \
    -e DSN=memory \
    -e SECRETS_COOKIE=some-cookie-secret \
    -e SECRETS_SYSTEM=some-system-secret \
    -e URLS_SELF_PUBLIC=http://127.0.0.1:4444/ \
    -e URLS_SELF_ADMIN=http://127.0.0.1:4445/ \
    -e URLS_SELF_ISSUER=http://127.0.0.1:4444/ \
    -e URLS_LOGIN=http://127.0.0.1:3001/login \
    -e URLS_CONSENT=http://127.0.0.1:3001/consent \
    -e URLS_LOGOUT=http://127.0.0.1:3001/logout \
    -e URLS_ERROR=http://127.0.0.1:3001/error \
    -e URLS_POST_LOGOUT_REDIRECT=http://127.0.0.1:3001/logout-successful \
    docker.io/oryd/hydra:v2.3 \
    serve all --dev
