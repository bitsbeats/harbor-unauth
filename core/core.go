package core

type (
	Config struct {
		URL        string          `json:"url"`
		Auths      map[string]Auth `json:"auths"`
		AllowList  []string        `json:"allowlist"`
		ProxyCount int             `json:"proxy_count"`
	}

	Auth struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
)
