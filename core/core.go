package core

type (
	Config struct {
		URL       string          `json:"url"`
		Auths     map[string]Auth `json:"auths"`
		AllowList []string        `json:"allowlist"`
	}

	Auth struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
)
