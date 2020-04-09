package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

// JWT deals with jwt token generation and validation
type JWT struct {
	expire time.Duration
	method jwt.SigningMethod
	secret []byte
}

// NewJWT is ctor for JWT
func NewJWT(expire time.Duration, signingAlgorithm string, secret []byte) (j *JWT, err error) {
	method := jwt.GetSigningMethod(signingAlgorithm)
	if method == nil {
		err = fmt.Errorf("invalid signingAlgorithm:%s", method)
		return
	}
	j = &JWT{expire: expire, method: method, secret: secret}
	return
}

const (
	// ExpireATKey for expire_at
	ExpireATKey = "expire_at"
	// CreatedKey for created
	CreatedKey = "created"
)

// Generate a token
func (j *JWT) Generate(values map[string]interface{}) (tokenString string, err error) {

	claims := jwt.MapClaims{
		ExpireATKey: time.Now().Add(j.expire).Unix(),
		CreatedKey:  time.Now().Unix(),
	}
	for k, v := range values {
		if _, ok := claims[k]; ok {
			err = fmt.Errorf("%s is reserved for claims", k)
			return
		}
		claims[k] = v
	}
	token := jwt.NewWithClaims(j.method, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(j.secret)

	return
}

var (
	// ErrInvalidSigningAlgorithm for invalid signing algorithm
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")
)

// Validate tokenString
func (j *JWT) Validate(tokenString string) (ok bool, claims jwt.MapClaims) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (key interface{}, err error) {
		if j.method != token.Method {
			err = ErrInvalidSigningAlgorithm
			return
		}
		key = j.secret
		return
	})
	if err != nil || !token.Valid {
		return
	}
	claims, ok = token.Claims.(jwt.MapClaims)

	expireAt, exists := claims[ExpireATKey]
	if !exists {
		ok = false
		return
	}

	ok = j.verifyExp(j.GetClaimInt64(expireAt))
	return
}

// GetClaimInt64 for retrieve claim value as int64
func (j *JWT) GetClaimInt64(value interface{}) (int64Value int64) {
	switch exp := value.(type) {
	case float64:
		int64Value = int64(exp)
	case json.Number:
		int64Value, _ = exp.Int64()
	}
	return
}

func (j *JWT) verifyExp(exp int64) (ok bool) {
	nowSecond := time.Now().Unix()
	ok = exp > nowSecond
	return
}
