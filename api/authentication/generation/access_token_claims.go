package generation

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AccessTokenClaims struct {
	UserId            int    `json:"user_id"`
	Email             string `json:"email"`
	IsStudent         bool   `json:"is_student"`
	IsAdmin           bool   `json:"is_admin"`
	IsOrganization    bool   `json:"is_organization"`
	IsUniversity      bool   `json:"is_university"`
	CreationTimestamp int64  `json:"iat"`
	TTL               int64  `json:"exp"`
}

func (c *AccessTokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.CreationTimestamp+c.TTL, 0)), nil
}
func (c *AccessTokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}
func (c *AccessTokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}
func (c *AccessTokenClaims) GetIssuer() (string, error) {
	return "", nil
}
func (c *AccessTokenClaims) GetSubject() (string, error) {
	return "", nil
}
func (c *AccessTokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}
