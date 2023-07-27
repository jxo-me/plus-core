package registry

import jwt "github.com/gogf/gf-jwt/v2"

type JwtRegistry struct {
	registry[*jwt.GfJWTMiddleware]
}
