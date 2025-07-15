package auth

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	ErrMissingHeader     = gerror.NewCode(gcode.New(210, "missing required headers", nil))
	ErrInvalidTimestamp  = gerror.NewCode(gcode.New(211, "invalid or expired timestamp", nil))
	ErrInvalidAPIKey     = gerror.NewCode(gcode.New(212, "invalid api key", nil))
	ErrSignatureMismatch = gerror.NewCode(gcode.New(213, "signature mismatch", nil))
	ErrReplayAttack      = gerror.NewCode(gcode.New(214, "replay attack detected", nil))
	ErrIPBanned          = gerror.NewCode(gcode.New(215, "client ip is banned", nil))
)
