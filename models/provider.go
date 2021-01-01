package models

type OAuth2Provider struct {
	AuthorizationEndpoint string `mapstructure:"authorization-endpoint"`
	TokenEndpoint         string `mapstructure:"token-endpoint"`
	ClientId              string `mapstructure:"client-id"`
	ClientSecret          string `mapstructure:"client-secret"`
}
