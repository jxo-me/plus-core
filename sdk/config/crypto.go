package config

import "github.com/jxo-me/plus-core/sdk/pkg/security"

type Crypto struct {
	Enable    bool                      `json:"enable" yaml:"enable"`
	Algorithm string                    `json:"algorithm" yaml:"algorithm"`
	Rc4       security.Rc4CipherConfig  `json:"rc4" yaml:"rc4"`
	Rsa       security.RsaCiphersConfig `json:"rsa" yaml:"rsa"`
	Aes       security.AesCipherConfig  `json:"aes" yaml:"aes"`
}
