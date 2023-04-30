package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("thesecret")

const TokenExpireDuration = time.Minute * 5
const RefreshTokenExpireDuration = time.Hour * 8

var ErrorInvalidToken = errors.New("invalid token")

type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

var keyFunc = func(token *jwt.Token) (interface{}, error) {
	return secret, nil
}

// GenToken 生成Token
func GenToken(userid int64) (aToken, rToken string, err error) {
	c := MyClaims{
		userid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt64("auth.jwt_expire")) * time.Hour).Unix(),
			Issuer:    "bluebell",
		},
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(secret)
	// refresh Token不需要存储任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(),
		Issuer:    "bluebell",
	}).SignedString(secret)

	return
}

func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	claims = new(MyClaims)
	var token *jwt.Token
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid {
		err = ErrorInvalidToken
	}
	return
}

func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}
	// 从旧access token中解析出claims数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)
	// 当access token已经过期并且refresh token没有过期的时候， 创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID)
	}
	return
}
