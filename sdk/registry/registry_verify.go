package registry

import (
	"github.com/jxo-me/plus-core/pkg/v2/security/auth"
)

type VerifyRegistry struct {
	registry[*auth.Verifier]
}
