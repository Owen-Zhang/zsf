package xjwt

import (
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	Name string
	ID   int
	Age  int
}

func Testjwt(t *testing.T) {
	signStr, err := Encrypt("123456",
		&User{
			Name: "test",
			ID:   123,
			Age:  20,
		}, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(signStr)
	fmt.Println("---------------------------------")
	var user User
	if err := Decrypt(signStr, "123456", &user); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(user)
	dataUser, err := json.Marshal(&user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(dataUser))
}

//go test -v -bench=. -run=none
func BenchmarkEncrypt(t *testing.B) {
	for i := 0; i < t.N; i++ {
		_, err := Encrypt("123456",
			&User{
				Name: "test",
				ID:   123,
				Age:  20,
			}, 10)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

/*
type JWTAuth struct {
	Name       string   `json:"name" valid:"ascii,required"`
	ExpireAt   int64    `json:"expireAt" valid:"required"`
	Mode       string   `json:"mode" valid:"in(HS256|HS384|HS512|RS256|ES256|ES384|ES512|RS384|RS512|PS256|PS384|PS512),required"`
	Secret     string   `json:"secret" valid:"ascii,required"`
	Exclude    []string `json:"exclude,omitempty"`
}*/
