package jwt

import (
	"errors"
	"github.com/bitbeliever/creve-grpc/configs"
	"github.com/bitbeliever/creve-grpc/model"
	"github.com/golang-jwt/jwt"
)

type userToken struct {
	model.User
	jwt.StandardClaims
}

func NewToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userToken{User: user})

	return token.SignedString([]byte(configs.Cfg.JWTKey))
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims
	t, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Cfg.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil

	//if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
	//	return claims, nil
	//} else {
	//	return nil, err
	//}
}
