package auth

import "errors"

var (
	ErrMissingHeader     = errors.New("missing required headers")
	ErrInvalidTimestamp  = errors.New("invalid or expired timestamp")
	ErrInvalidAPIKey     = errors.New("invalid api key")
	ErrSignatureMismatch = errors.New("signature mismatch")
	ErrReplayAttack      = errors.New("replay attack detected")
	ErrIPBanned          = errors.New("client ip is banned")
)
