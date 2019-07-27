package token

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	config "github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
)

// CustomClaims 自定义的 metadata在加密后作为 JWT 的第二部分返回给客户端
type CustomClaims struct {
	UserName string   `json:"user_name"`
	UserID   string   `json:"user_id"`
	PerAddr  string   `json:"per_addr"`
	Roles    []string `json:"roles"`

	jwt.StandardClaims
}

// Token jwt服务
var privateKey []byte

// InitConfig 初始化
func InitConfig(filePath string, path ...string) {
	fileSource := file.NewSource(
		file.WithPath(filePath),
	)
	conf := config.NewConfig()
	err := conf.Load(fileSource)
	if err != nil {
		log.Fatal(err)
	}

	privateKey = conf.Get(path...).Bytes()
	if err != nil {
		log.Fatal(err)
	}
}

//Decode 解码
func Decode(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})

	if err != nil {
		return nil, err
	}
	// 解密转换类型并返回
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, err
}

// Encode 将 User 用户信息加密为 JWT 字符串
// expireTime := time.Now().Add(time.Hour * 24 * 3).Unix() 三天后过期
func Encode(userName, userID, perAddr string, roles []string, issuer string, expireTime int64) (string, error) {
	claims := CustomClaims{
		userName,
		userID,
		perAddr,
		roles,
		jwt.StandardClaims{
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(privateKey)
}
