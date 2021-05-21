package core

type (
	Config struct {
		URL   string          `json:"url"`
		Auths map[string]Auth `json:"auths"`
	}

	Auth struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
)
