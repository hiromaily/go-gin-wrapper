// Package jwts is authentication by JWT
package jwts

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	// HMAC signing algorithm
	HMAC uint8 = 1
	// RSA signing algorithm
	RSA uint8 = 2
)

// CustomClaims is jwt claim
type CustomClaims struct {
	Option string `json:"option"`
	jwt.StandardClaims
}

// JWTer interface
type JWTer interface {
	CreateBasicToken(t int64, clientID, userName string) (string, error)
	CreateCustomToken(t int64, clientID, userName, option string) (string, error)
	ValidateToken(tokenString string) error
	ValidateTokenWithClaim(tokenString, clientID, userName string) error
	ValidateTokenWithCustomClaim(tokenString, clientID, userName, option string) error
}

// ----------------------------------------------------------------------------
// jwtee
// ----------------------------------------------------------------------------

type jwtee struct {
	audience string
	SigAlgoer
}

// NewJWT returns JWTer interface
func NewJWT(audience string, sigAlgoer SigAlgoer) JWTer {
	return &jwtee{
		audience:  audience,
		SigAlgoer: sigAlgoer,
	}
}

// CreateBasicToken returns basic claim
// - encode Header, Payload, Signature by Base64 and concatenate these with dot
func (j *jwtee) CreateBasicToken(t int64, clientID, userName string) (string, error) {
	token := jwt.NewWithClaims(
		j.SigAlgoer.GetMethod(),
		j.getClaims(t, clientID, userName),
	)
	return j.SigAlgoer.SignedString(token)
}

// CreateToken returns user defined claim
// - encode Header, Payload, Signature by Base64 and concatenate these with dot
func (j *jwtee) CreateCustomToken(t int64, clientID, userName, option string) (string, error) {
	claims := &CustomClaims{
		option,
		j.getClaims(t, clientID, userName),
	}
	token := jwt.NewWithClaims(j.SigAlgoer.GetMethod(), claims)
	return j.SigAlgoer.SignedString(token)
}

// Payload
func (j *jwtee) getClaims(t int64, clientID, userName string) jwt.StandardClaims {
	// Audience  string `json:"aud,omitempty"` // https://login.hiromaily.com
	// ExpiresAt int64  `json:"exp,omitempty"`
	// Id        string `json:"jti,omitempty"`
	// IssuedAt  int64  `json:"iat,omitempty"`
	// Issuer    string `json:"iss,omitempty"` // OAuth client_id
	// NotBefore int64  `json:"nbf,omitempty"`
	// Subject   string `json:"sub,omitempty"` // user name or email
	claims := jwt.StandardClaims{
		Audience: j.audience,
		// ExpiresAt: time.Now().Add(time.Second * 2).Unix(),
		ExpiresAt: t,
		Issuer:    clientID,
		Subject:   userName,
	}
	return claims
}

func (j *jwtee) keyFunc(token *jwt.Token) (interface{}, error) {
	if isValid := j.SigAlgoer.ValidateMethod(token); !isValid {
		return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return j.SigAlgoer.GetKey(), nil
}

// ValidateToken validates token string, it may be too strict
func (j *jwtee) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, j.keyFunc)
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}
	return nil
}

// ValidateTokenWithClaim validates token by clientID and userName
// may be too strict to check
func (j *jwtee) ValidateTokenWithClaim(tokenString, clientID, userName string) error {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, j.keyFunc)
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}

	// validate claim
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return errors.New("fail to validate claim: claims can't be retrieved")
	}
	if claims.Issuer != clientID {
		return errors.New("fail to validate claim: issuer is invalid")
	}
	if claims.Subject != userName {
		return errors.New("fail to validate claim: subject is invalid")
	}
	return nil
}

// ValidateTokenWithCustomClaim validates token by clientID and userName and option
// may be too strict to check
func (j *jwtee) ValidateTokenWithCustomClaim(tokenString, clientID, userName, option string) error {
	// token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, j.keyFunc)
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token is invalid")
	}

	// validate claim
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return errors.New("fail to validate claim: claims data can't be retrieved")
	} else if claims.Issuer != clientID {
		return errors.New("fail to validate claim: issuer is invalid")
	} else if claims.Subject != userName {
		return errors.New("fail to validate claim: subject is invalid")
	} else if claims.Option != option {
		return errors.New("fail to validate claim: option is invalid")
	}

	return nil
}

// ----------------------------------------------------------------------------
// SigAlgoer
// ----------------------------------------------------------------------------

// SigAlgoer interface
type SigAlgoer interface {
	GetMethod() jwt.SigningMethod
	SignedString(token *jwt.Token) (string, error)
	ValidateMethod(token *jwt.Token) bool
	GetKey() interface{}
}

// ----------------------------------------------------------------------------
// HMAC
// ----------------------------------------------------------------------------
type algoHMAC struct {
	encrypted uint8
	method    jwt.SigningMethod
	secret    string
}

// NewHMAC returns SigAlgoer
func NewHMAC(secret string) SigAlgoer {
	return &algoHMAC{
		encrypted: HMAC,
		method:    jwt.SigningMethodHS256,
		secret:    secret,
	}
}

// GetMethod returns method
func (a *algoHMAC) GetMethod() jwt.SigningMethod {
	return a.method
}

// SignedString returns signed string from toke
func (a *algoHMAC) SignedString(token *jwt.Token) (string, error) {
	return token.SignedString([]byte(a.secret))
}

// ValidateMethod validates method of token
func (a *algoHMAC) ValidateMethod(token *jwt.Token) bool {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	return ok
}

// GetKey returns key
func (a *algoHMAC) GetKey() interface{} {
	return []byte(a.secret)
}

// ----------------------------------------------------------------------------
// RSA
// ----------------------------------------------------------------------------
type algoRSA struct {
	encrypted  uint8
	method     jwt.SigningMethod
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSA returns SigAlgoer interface and error
func NewRSA(privKey, pubKey string) (SigAlgoer, error) {
	privKeyParsed, err := lookupPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	pubKeyParsed, err := lookupPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	return &algoRSA{
		encrypted:  RSA,
		method:     jwt.SigningMethodRS256,
		privateKey: privKeyParsed,
		publicKey:  pubKeyParsed,
	}, nil
}

func (a *algoRSA) GetMethod() jwt.SigningMethod {
	return a.method
}

func (a *algoRSA) SignedString(token *jwt.Token) (string, error) {
	return token.SignedString(a.privateKey)
}

func (a *algoRSA) ValidateMethod(token *jwt.Token) bool {
	_, ok := token.Method.(*jwt.SigningMethodRSA)
	return ok
}

func (a *algoRSA) GetKey() interface{} {
	return a.publicKey
}

// public key using ParseRSAPublicKeyFromPEM()
func lookupPublicKey(keyPath string) (*rsa.PublicKey, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	return parsedKey, err
}

// private key using ParseRSAPrivateKeyFromPEM()
func lookupPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	return parsedKey, err
}
