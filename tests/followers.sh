#!/bin/bash

AUTH=https://auth.prod.monitalks.io
API=https://int-api.prod.monitalks.io/function
FAAS=https://openfaas.prod.monitalks.io/function

COGNITO_URL=${AUTH}/cognito

FOLLOWING_URL=${API}/user-rels/following

ACCESS=$(curl -X POST -H "Cache-Control: no-cache" -H "Content-Type: application/x-www-form-urlencoded" -d $CREDENTIAL_EMAIL -d $CREDENTIAL_PASS $COGNITO_URL)

echo $ACCESS

# Random User
USER=4084ece8-64dc-4ccb-917d-b9df47931aeb

# Sample URLs
COLL_URL=${FOLLOWING_URL}?userId=$USER

# Retrieve
curl -H "Authorization: $ACCESS" $COLL_URL | jq '.[]'
