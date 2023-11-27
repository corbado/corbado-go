# Corbado Go SDK

Go SDK for Corbado Backend API

[![Go Reference](https://pkg.go.dev/badge/github.com/corbado/corbado-go.svg)](https://pkg.go.dev/github.com/corbado/corbado-go)
[![Test Status](https://github.com/corbado/corbado-go/workflows/tests/badge.svg)](https://github.com/corbado/corbado-go/actions?query=workflow%3Atests)
[![documentation](https://img.shields.io/badge/documentation-Corbado_Backend_API_Reference-blue.svg)](https://api.corbado.com/docs/api/)
[![Go Report Card](https://goreportcard.com/badge/github.com/corbado/corbado-go)](https://goreportcard.com/report/github.com/corbado/corbado-go)

## Requirements

The SDK supports Go version 1.18 and above.

## Usage

```
$ go get github.com/corbado/corbado-go@v0.4.2
```

Import SDK in your Go files:

```go
import "github.com/corbado/corbado-go"
```

Now create a new SDK client:

```go
config := corbado.MustNewConfig("pro-12345678", "yoursecret")
sdk, err := corbado.NewSDK(config)
if err != nil {
	// handle error
}

// list all users
users, err := sdk.Users().List(context.TODO(), nil)
if err != nil {
    if serverErr := corbado.AsServerError(err); serverErr != nil {
	    // handle server error	
    }
}
```

See [examples](https://github.com/corbado/corbado-go/tree/main/examples) for some real example code
