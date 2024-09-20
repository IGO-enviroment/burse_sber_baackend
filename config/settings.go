package config

type Settings struct {
	Port      int    `json:"port"`
	SecretKey string `json:"jwt_secret"`
}
