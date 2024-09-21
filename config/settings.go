package config

type Settings struct {
	Port int `json:"port"`

	// JWT
	SecretKey      string `json:"jwt_secret"`
	AccessTokenTTL int64  `json:"access_token_ttl"`
}
