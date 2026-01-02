package pkg

import (
	"5/work/Newyear/user-srv/basic/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func GetJwtToken(userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	iat := time.Now().Unix()
	claims["exp"] = iat + 86400
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	secretKey := config.AppConf.Jwt.SecretKey
	return token.SignedString([]byte(secretKey))
}
func ParseJwtToken(tokenString string) (float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		secretKey := config.AppConf.Jwt.SecretKey
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims["userId"].(float64), err
	} else {
		fmt.Println(err)
	}
	return 0, nil
}
