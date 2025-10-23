#!/usr/bin/env bash

set -e

eval $(./create-client.sh)

docker build --tag gologin-test-app:0.5.0 https://github.com/ngyewch/gologin-test-app.git#v0.5.0
docker run --rm -it \
    --network host \
    -p 8080:8080 \
    -e GOLOGIN_OIDC_ISSUERURL=http://127.0.0.1:4444/ \
    -e GOLOGIN_OIDC_CLIENTID=${CLIENT_ID} \
    -e GOLOGIN_OIDC_CLIENTSECRET=${CLIENT_SECRET} \
    -e GOLOGIN_OIDC_SCOPES=openid,profile,email \
    gologin-test-app:0.5.0 \
    gologin-test-app serve
