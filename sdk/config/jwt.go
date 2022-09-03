package config

type Jwt struct {
	Secret      string `yaml:"secret" json:"secret"`
	SigningKey  string `yaml:"signingKey" json:"signing_key"`
	Timeout     int64  `yaml:"timeout" json:"timeout"`
	MaxRefresh  int64  `yaml:"maxRefresh" json:"max_refresh"`
	IdentityKey string `yaml:"identityKey" json:"identity_key"`
}

var JwtConfig = new(Jwt)
