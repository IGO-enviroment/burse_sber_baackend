package config

type Settings struct {
	Port int  `json:"port"`
	Smtp Smtp `json:"smtp"`

	// JWT
	JwtSecret      string `json:"jwt_secret"`
	AccessTokenTTL int64  `json:"access_token_ttl"`

	// DB
	PgConnString string `json:"pg_conn_string"`
}

type Smtp struct {
	From     string `json:"from"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}
