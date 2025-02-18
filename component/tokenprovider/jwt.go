package tokenprovider

import (
	"fmt"
	"time"
	"todololist/common"

	"github.com/dgrijalva/jwt-go"
)

type jwtProvider struct {
	prefix    string
	secretKey string
}

func NewTokenJWTProvider(prefix string, secret string) *jwtProvider {
	return &jwtProvider{prefix: prefix, secretKey: secret}
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.StandardClaims
}

type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}

func (j *jwtProvider) SecretKey() string {
	return j.secretKey
}

// Generate tạo JWT token
func (j *jwtProvider) Generate(data TokenPayload, expiry int) (Token, error) {
	// Thời gian hiện tại
	now := time.Now()

	// Tạo JWT với payload
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		Payload: common.TokenPayload{
			UId:   data.UserId(),
			URole: data.Role(),
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt:  now.Unix(),
			Id:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})

	// Ký token với secret
	myToken, err := t.SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, err
	}

	// return token
	return &token{
		Token:   myToken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (TokenPayload, error) {
	// Parse token với custom claims
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil || !res.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Trả về payload hợp lệ
	return claims.Payload, nil
}
