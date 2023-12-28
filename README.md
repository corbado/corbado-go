<img width="1070" alt="GitHub Repo Cover" src="https://github.com/corbado/corbado-php/assets/18458907/aa4f9df6-980b-4b24-bb2f-d71c0f480971">

# Corbado Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/corbado/corbado-go.svg)](https://pkg.go.dev/github.com/corbado/corbado-go)
[![License](https://poser.pugx.org/corbado/php-sdk/license.svg)](https://packagist.org/packages/corbado/php-sdk)
[![Test Status](https://github.com/corbado/corbado-go/workflows/tests/badge.svg)](https://github.com/corbado/corbado-go/actions?query=workflow%3Atests)
[![Go Report Card](https://goreportcard.com/badge/github.com/corbado/corbado-go)](https://goreportcard.com/report/github.com/corbado/corbado-go)
[![documentation](https://img.shields.io/badge/documentation-Corbado_Backend_API_Reference-blue.svg)](https://api.corbado.com/docs/api/)
[![Slack](https://img.shields.io/badge/slack-join%20chat-brightgreen.svg)](https://join.slack.com/t/corbado/shared_invite/zt-1b7867yz8-V~Xr~ngmSGbt7IA~g16ZsQ)

The [Corbado](https://www.corbado.com) Go SDK provides convenient access to the [Corbado Backend API](https://api.corbado.com/docs/api/) from applications written in the Go language.

:rocket: [Getting started](#rocket-getting-started) | :hammer_and_wrench: [Services](#hammer_and_wrench-services) | :books: [Advanced](#books-advanced) | :speech_balloon: [Support & Feedback](#speech_balloon-support--feedback)

## :rocket: Getting started

### Requirements

- Go 1.18 or later

### Installation

Use the following command to install the Corbado Go SDK:

```bash
go get github.com/corbado/corbado-go@v0.6.0
```

### Usage

To create a Corbado Go SDK instance you need to provide your `Project ID` and `API secret` which can be found at the [Developer Panel](https://app.corbado.com).

```Go
configuration, err := corbado.NewConfiguration("<Project ID>", "<API secret>")
if err != nil {
    // handle error
}

sdk, err := corbado.NewSDK(configuration)
if err != nil {
    // handle error
}
```

### Examples

A list of examples can be found in the integration tests [here](tests/integration).

## :hammer_and_wrench: Services

The Corbado Go SDK provides the following services:

- `emailMagicLinks` for managing email magic links ([examples](tests/integration/emailmagiclink))
- `emailOTPs` for managing email OTPs ([examples](tests/integration/emailotp))
- `smsOTPs` for managing SMS OTPs ([examples](tests/integration/smsotp))
- `sessions` for managing sessions
- `users` for managing users ([examples](tests/integration/user))
- `authTokens` for managing authentication tokens needed for own session management ([examples](tests/integration/authtoken))
- `validations` for validating email addresses and phone numbers ([examples](tests/integration/validation))

To use a specific service, such as `sessions`, invoke it as shown below:

```Go
user, err := sdk->sessions()->GetCurrentUser()
if err != nil {
    // handle error
}
``` 

## :books: Advanced

### Error handling

The Corbado PHP SDK throws exceptions for all errors. The following exceptions are thrown:

- `AssertException` for failed assertions (client side)
- `ConfigurationException` for configuration errors (client side)
- `ServerException` for server errors (server side)
- `StandardException` for everything else (client side)

If the Backend API returns a HTTP status code other than 200, the Corbado PHP SDK throws a `ServerException`. The `ServerException`class provides convenient methods to access all important data:

```Go
TODO
```

## :speech_balloon: Support & Feedback

### Report an issue

If you encounter any bugs or have suggestions, please [open an issue](https://github.com/corbado/corbado-go/issues/new).

### Slack channel

Join our Slack channel to discuss questions or ideas with the Corbado team and other developers.

[![Slack](https://img.shields.io/badge/slack-join%20chat-brightgreen.svg)](https://join.slack.com/t/corbado/shared_invite/zt-1b7867yz8-V~Xr~ngmSGbt7IA~g16ZsQ)

### Email

You can also reach out to us via email at vincent.delitz@corbado.com.

### Vulnerability reporting

Please report suspected security vulnerabilities in private to security@corbado.com. Please do NOT create publicly viewable issues for suspected security vulnerabilities.