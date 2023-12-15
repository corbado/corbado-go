#!/bin/sh

if [ $# -ne 1 ]; then
    echo "USAGE: <output path>"
    echo "EXAMPLE: pkg/sdk/entity"

    exit 1
fi

path=$1

oapi-codegen -package common -generate "types" common.yml > ${path}/common/common.gen.go
oapi-codegen -package api -generate "types,client" -import-mapping common.yml:github.com/corbado/corbado-go/pkg/sdk/entity/common backend_api_public.yml > ${path}/api/api.gen.go