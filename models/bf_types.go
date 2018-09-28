package models

const (
	openIdRequestPath                   string = "https://login.botframework.com/v1/.well-known/openidconfiguration"
	authorizationHeaderValuePrefix      string = "Bearer "
	wrongAuthorizationHeaderFormatError string = "The provided authorization header is in the wrong format: %v"
	wrongSplitLengthError               string = "The authorize value split length with character \"%v\" is not valid: %v (%v)"
	splitCharacter                      string = "."
	issuerUrl                           string = "https://api.botframework.com"
)

type OpenIdDocument struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	JwksURI                           string   `json:"jwks_uri"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
}

type SigningKeys struct {
	Keys []struct {
		Kty          string   `json:"kty"`
		Use          string   `json:"use"`
		KeyId        string   `json:"kid"`
		X5T          string   `json:"x5t"`
		N            string   `json:"n"`
		E            string   `json:"e"`
		X5C          []string `json:"x5c"`
		Endorsements []string `json:"endorsements,omitempty"`
	} `json:"keys"`
}

type JwtHeader struct {
	Type            string `json:"typ"`
	Algorithm       string `json:"alg"`
	SigningKeyId    string `json:"kid"`
	SigningKeyIdX5T string `json:"x5t"`
}

type JwtPayload struct {
	ServiceUrl   string `json:"serviceurl"`
	Issuer       string `json:"iss"`
	Audience     string `json:"aud"`
	Expires      int    `json:"exp"`
	CreatedOnNbf int    `json:"nbf"`
}

type MicrosoftJsonWebToken struct {
	HeaderBase64, PayloadBase64 string
	Header                      JwtHeader
	Payload                     JwtPayload
	VerifySignature             []byte
}
