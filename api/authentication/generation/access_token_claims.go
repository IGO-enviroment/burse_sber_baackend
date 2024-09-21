package generation

type AccessTokenClaims struct {
	UserId              string `json:"user_id"`
	UserName            string `json:"user_name"`
	Email               string `json:"email"`
	IsStudent           bool   `json:"is_smartagent"`
	AccountId           int    `json:"account_id"`
	IsBackofficeManager bool   `json:"is_backoffice_manager"`
	IsEmployer          bool   `json:"is_employer"`
	CreationTimestamp   int64  `json:"iat"`
	TTL                 int64  `json:"exp"`
}
