package config

type Settings struct {
	Port int `json:"port"`

	// JWT
	JwtSecret      string `json:"jwt_secret"`
	AccessTokenTTL int64  `json:"access_token_ttl"`

	// DB
	PgConnString string `json:"pg_conn_string"`
}
