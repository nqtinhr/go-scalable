package tokenprovider

import (
	"todololist/common"

	"github.com/dgrijalva/jwt-go"
)

type jwtProvider struct {
	prefix    string
	secretKey string
}

func NewTokenJWTProvider(prefix string) *jwtProvider {
	return &jwtProvider{prefix: prefix}
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.StandardClaims
}
