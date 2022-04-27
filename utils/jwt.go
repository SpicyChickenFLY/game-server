package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	uuid "github.com/satori/go.uuid"
)

const (
	expirationTime   = 7 * 24 * time.Hour
	expiredNotBefore = 30 * time.Minute
)

// 签名算法, 随机, 不保存密钥, 每次都是随机的
var privateKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
var publicKey = &privateKey.PublicKey
var hs = jwt.NewES256(
	jwt.ECDSAPublicKey(publicKey),
	jwt.ECDSAPrivateKey(privateKey),
)

// 记录登录信息的 JWT
type LoginToken struct {
	jwt.Payload
	Nickname string `json:"nickname"`
}

// Sign jwt token with id
func Sign(nickname string) ([]byte, error) {
	now := time.Now()
	pl := LoginToken{
		Payload: jwt.Payload{
			Issuer:         "coolcat",
			Subject:        "login",
			Audience:       jwt.Audience{},
			ExpirationTime: jwt.NumericDate(now.Add(expirationTime)),
			NotBefore:      jwt.NumericDate(now.Add(expiredNotBefore)),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          uuid.NewV4().String(),
		},
		Nickname: nickname,
	}
	token, err := jwt.Sign(pl, hs)
	return token, err
}

// Verify the jwt token
func Verify(token []byte) (*LoginToken, error) {
	pl := &LoginToken{}
	_, err := jwt.Verify(token, hs, pl)
	return pl, err
}
