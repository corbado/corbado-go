#!/bin/sh

if [ $# -ne 1 ]; then
    echo "USAGE: <output path>"
    echo "EXAMPLE: pkg/sdk/entity"

    exit 1
fi

path=$1

# download OpenAPI specs
curl -s -O https://api.corbado.com/docs/api/openapi/common.yml
curl -s -O https://api.corbado.com/docs/api/openapi/backend_api_public.yml

# generate Go entities
./scripts/generate-openapi-entities-go.sh ${path}

# remove OpenAPI specs
rm common.yml
rm backend_api_public.yml