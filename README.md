# OAuth2 PKCE Proxy

![Build Status](https://github.com/bwilczynski/oauth2-pkce-proxy/workflows/verify/badge.svg) ![Go Report](https://goreportcard.com/badge/github.com/bwilczynski/oauth2-pkce-proxy)

Written in Go this small app will proxy [PKCE](https://tools.ietf.org/html/rfc7636#section-4.1) requests to OAuth2 providers that currently do not support it. It was written mostly to simplify the usage of [Strava CLI](https://github.com/bwilczynski/strava-cli) and [Sonos CLI](https://github.com/bwilczynski/sonos-cli) that currently require users to create their own client applications in order to use it.

# Installation

Currently there is no installer provided. You need to build it from source:

```
go build .
```

Docker and Homebrew are coming soon!

# Usage

Use `authorization-endpoint` and `token-endpoint` flags passed to oauth2-pkce-proxy to configure OAuth provider (by default it points to [Strava OAuth2](https://developers.strava.com/docs/authentication/)). Pass your application's `client-id` and `client-secret` that will be used for further authentication calls.

You can also use environment variables, all CLI flags have equivalents starting with PKCE_PROXY and using snake case, for example:

```
export PKCE_PROXY_AUTHORIZATION_ENDPOINT=https://www.strava.com/oauth/authorize
```

Run the server (by default it will listen on 8080):

```
oauth2-pkce-proxy \
      --client-id "OAuth Client ID" \
      --client-secret "OAuth Client Secret" \
      --console
```

Call it using excellent [step](https://github.com/smallstep/cli) utility:

```console
❯ step oauth \
    --authorization-endpoint http://localhost:8080/authorize \
    --token-endpoint http://localhost:8080/token \
    --client-id "OAuth Client ID" \
    --scope activity:read_all

Your default web browser has been opened to visit:

http://localhost:8080/authorize?client_id=31056&code_challenge=q8NDohZJemwxjQNVd1zQySckYjLQQm8lADn4obl_K-s&code_challenge_method=S256&nonce=208e331e3fbaa84d1579d4e16c7e0de17ce9b900e84d0218b6b2b27b18bb3706&redirect_uri=http%3A%2F%2F127.0.0.1%3A51195&response_type=code&scope=activity%3Aread_all&state=oAVV5t1wMmLhHJ1sUFszRAZDNK3AYtxJ
{
  "access_token": "5181012ebe596a4b6ff9cc7f5e82fee4a8e685df",
  "id_token": "",
  "refresh_token": "707e29ebe8ee2f2115c3d24c3f5de00c7a35bd47",
  "expires_in": 21600,
  "token_type": "Bearer"
}
```

For usage and help content, pass in the --help parameter, for example:

```sh
./oauth2-pkce-proxy --help
```

# Features

- Graceful shutdown
- Prometheus metrics
- Structured logging with zerolog
- 12-factor app with viper

# Endpoints

| Endpoint         | Description                                                                                                                                                                   |
| ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `GET /authorize` | Used by the client to [send the code challenge](https://tools.ietf.org/html/rfc7636#section-4.3) as part of the OAuth 2.0 Authorization Request                               |
| `GET /code`      | Callback endpoint used by the server to [return the code](https://tools.ietf.org/html/rfc7636#section-4.4)                                                                    |
| `POST /token`    | Used by the client to [send Authorization Code and the Code Verifier to the Token Endpoint](https://tools.ietf.org/html/rfc7636#section-4.5) and exchange it for access token |
| `GET /metrics`   | Get Prometheus metrics                                                                                                                                                        |
