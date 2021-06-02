package core

type (
	Config struct {
		URL        string   `json:"url"`
		Auth       Auth     `json:"auth"`
		Projects   []string `json:"projects"`
		AllowList  []string `json:"allowlist"`
		ProxyCount int      `json:"proxy_count"`
	}

	Auth struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
)
