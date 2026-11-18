package token

import (
	"log"
	"os"
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

// weakDefaultKey is the insecure signing key that was historically committed to
// conf/jwt.json. The service refuses to start if it is still in use.
const weakDefaultKey = "asdfasf"

// InitConfig 初始化
func InitConfig(filePath string, path ...string) {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		privateKey = []byte(secret)
	} else {
		fileSource := file.NewSource(
			file.WithPath(filePath),
		)
		conf := config.NewConfig()
		if err := conf.Load(fileSource); err != nil {
			log.Fatal(err)
		}
		privateKey = conf.Get(path...).Bytes()
	}

	if len(privateKey) == 0 || string(privateKey) == weakDefaultKey {
		log.Fatal("jwt signing key is empty or uses the known insecure default; set a strong JWT_SECRET environment variable")
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
