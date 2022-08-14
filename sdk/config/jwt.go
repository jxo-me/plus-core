package config

type Jwt struct {
	Secret      string
	SigningKey  string
	Timeout     int64
	MaxRefresh  int64
	IdentityKey string
}

var JwtConfig = new(Jwt)
