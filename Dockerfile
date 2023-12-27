FROM golang:1.18

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0

WORKDIR /app