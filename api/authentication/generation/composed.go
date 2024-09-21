package generation

import (
	"boilerplate/config"
	"boilerplate/jwt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func NewJWTToken(token string, settings *config.Settings) (*JWT, error) {
	if token == "" {
		return nil, fmt.Errorf("Empty token")
	}
	jwtParts := strings.Split(token, ".")
	if len(jwtParts) != 3 {
		return nil, fmt.Errorf("Invalid JWT")
	}
	encodedPayloadBytes, err := base64.StdEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid JWT payload encoding")
	}
	var claims AccessTokenClaims
	if err := json.Unmarshal(encodedPayloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("Invalid JWT payload JSON")
	}
	return &JWT{
		settings:       settings,
		header:         jwtParts[0],
		payload:        jwtParts[1],
		encodedPayload: string(encodedPayloadBytes),
		claims:         claims,
		signature:      jwtParts[2],
	}, nil
}

type JWT struct {
	settings       *config.Settings
	header         string
	payload        string
	claims         AccessTokenClaims
	signature      string
	encodedPayload string
}

func (t *JWT) IsValid() bool {
	return t.isHeaderValid() && t.isSignatureAuthentic() && !t.isExpired()
}

func (t *JWT) Claims() AccessTokenClaims {
	return t.claims
}

func (t *JWT) isHeaderValid() bool {
	if t.header == "" || t.payload == "" || t.signature == "" {
		return false
	}

	headerStr, err := base64.StdEncoding.DecodeString(t.header)
	if err != nil {
		return false
	}

	var header jwt.Header
	if err = json.Unmarshal(headerStr, &header); err != nil {
		return false
	}

	return header.Type == jwt.TokenTypeJWT && header.Algorithm == jwt.HS256
}

func (t *JWT) isSignatureAuthentic() bool {
	computedSignature := jwt.GetSignature(t.header, t.payload, t.settings.JwtSecret)
	return computedSignature == t.signature
}

func (t *JWT) isExpired() bool {
	return t.claims.CreationTimestamp+t.claims.TTL < time.Now().UTC().Unix()
}
