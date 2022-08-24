package config

import "context"

type Initialize interface {
	String() string
	Init(ctx context.Context, s *Settings) error
}
