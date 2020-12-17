# OAuth2 PKCE Proxy

![Build Status](https://github.com/bwilczynski/oauth2-pkce-proxy/workflows/verify/badge.svg) ![Go Report](https://goreportcard.com/badge/github.com/bwilczynski/oauth2-pkce-proxy)

Written in Go this small app will proxy [PKCE](https://tools.ietf.org/html/rfc7636#section-4.1) requests to OAuth2 providers that currently do not support it. It was written mostly to simplify the usage of [Strava CLI](https://github.com/bwilczynski/strava-cli) and [Sonos CLI](https://github.com/bwilczynski/sonos-cli) that currently require users to create their own client applications in order to use it.
