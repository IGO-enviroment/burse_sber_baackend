package jwt

type Header struct {
	Algorithm TokenAlgorithm `json:"alg"`
	Type      TokenType      `json:"typ"`
}

const (
	HS256        = TokenAlgorithm("HS256") // HMAC SHA256
	TokenTypeJWT = TokenType("JWT")
)

type (
	TokenAlgorithm string
	TokenType      string
)
