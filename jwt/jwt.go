package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func GetToken(claims interface{}, secret string) string {
	header := getHeader()
	payload := getPayload(claims)
	signature := GetSignature(header, payload, secret)
	return fmt.Sprintf("%s.%s.%s", header, payload, signature)
}

func getHeader() string {
	header := Header{
		Algorithm: HS256,
		Type:      TokenTypeJWT,
	}
	headerJson, err := json.Marshal(header)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(headerJson)
}

func getPayload(claims interface{}) string {
	payloadJson, err := json.Marshal(claims)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(payloadJson)
}

func GetSignature(header string, payload string, secret string) string {
	secretBase64 := base64.StdEncoding.EncodeToString([]byte(secret))
	signatureContent := fmt.Sprintf("%s.%s", header, payload)
	return computeHmacSha256(signatureContent, secretBase64)
}

func computeHmacSha256(content string, secret string) string {
	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
