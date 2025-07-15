package auth

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	ErrMissingHeader     = gerror.NewCode(gcode.New(210, "msg_missing_required_headers", nil))
	ErrInvalidTimestamp  = gerror.NewCode(gcode.New(211, "msg_invalid_or_expired_timestamp", nil))
	ErrInvalidAPIKey     = gerror.NewCode(gcode.New(212, "msg_authentication_failed", nil))
	ErrSignatureMismatch = gerror.NewCode(gcode.New(213, "msg_invalid_signature", nil))
	ErrReplayAttack      = gerror.NewCode(gcode.New(214, "msg_replay_attack_detected", nil))
	ErrIPBanned          = gerror.NewCode(gcode.New(215, "msg_client_ip_is_banned", nil))
)
