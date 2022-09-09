package config

type Crypto struct {
	Enable    bool   `json:"enable" yaml:"enable"`
	Algorithm string `json:"algorithm" yaml:"algorithm"`
}
