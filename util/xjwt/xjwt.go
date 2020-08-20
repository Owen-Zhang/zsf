package xjwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//Encrypt 加密生成jwt字符串
func Encrypt(secret string, data interface{}, timeout int64) (string, error) {
	expireAt := time.Now().Unix() + timeout
	if timeout == 0 {
		expireAt = 0
	}
	claims := &jwt.MapClaims{
		"exp":  expireAt,
		"data": data,
	}
	method := jwt.GetSigningMethod("HS256")
	token := jwt.NewWithClaims(method, claims)
	ss, err := token.SignedString([]byte(secret))
	return ss, err
}

//Decrypt JWT解密
func Decrypt(signStr string, secret string) (interface{}, error) {
	token, err := jwt.Parse(signStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// if err := mapstructure.Decode(claims["data"], &out); err != nil {
		// 	return fmt.Errorf("mapstructure.Decode 出错: %v", err)
		// }
		return claims["data"], nil
	}
	return nil, fmt.Errorf("输入参数签名错误")
}
