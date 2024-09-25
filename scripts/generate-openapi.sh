#!/bin/sh

if [ $# -ne 1 ]; then
    echo "USAGE: <output path>"
    echo "EXAMPLE: pkg/generated"

    exit 1
fi

path=$1

# To generate openapi client, you need to copy common.yml and backend_api_public.yml to the root of the project

# generate Go entities and clients
oapi-codegen -package common -generate "types" common.yml > ${path}/common/common.gen.go
oapi-codegen -package api -import-mapping common.yml:github.com/corbado/corbado-go/v2/pkg/generated/common -generate "types,client" backend_api_public_v2.yml > ${path}/api/api.gen.go
