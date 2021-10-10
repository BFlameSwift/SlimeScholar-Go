package model

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserID   uint64 `json:"user_id"`
	Password string `json:"password"`
	Username string `json:"username"`
}
