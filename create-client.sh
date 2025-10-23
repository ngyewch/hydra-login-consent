#!/usr/bin/env bash

set -e

client=$(docker run -it --rm \
    --network host \
    docker.io/oryd/hydra:v2.3 \
    create client \
    --endpoint http://127.0.0.1:4445/ \
    --grant-type authorization_code,refresh_token \
    --response-type code,id_token \
    --scope openid \
    --scope email \
    --scope profile \
    --scope offline_access \
    --name gologin-test-app \
    --client-uri http://127.0.0.1:8080 \
    --redirect-uri http://127.0.0.1:8080/oauth2/oidc/callback \
    --format json \
    )
client_id=$(echo $client | jq -r '.client_id')
client_secret=$(echo $client | jq -r '.client_secret')
echo "export CLIENT_ID=$client_id"
echo "export CLIENT_SECRET=$client_secret"
