#!/bin/sh

if [ $# -ne 1 ]; then
    echo "USAGE: <output path>"
    echo "EXAMPLE: pkg/sdk/entity/api"

    exit 1
fi

path=$1

# download OpenAPI specs
curl -s -O https://api.corbado.com/docs/api/openapi/backend_api_public.yml

# generate Go entities and clients
oapi-codegen -package api -generate "types,client" backend_api_public.yml > ${path}/api.gen.go

# remove OpenAPI specs
rm backend_api_public.yml