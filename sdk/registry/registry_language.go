package registry

import (
	"github.com/gogf/gf/v2/i18n/gi18n"
)

type LanguageRegistry struct {
	registry[*gi18n.Manager]
}
