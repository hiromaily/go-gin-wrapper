package jwt

import (
	"crypto/rsa"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	lg "github.com/hiromaily/golibs/log"
	"io/ioutil"
	//"time"
)

type CustomClaims struct {
	Option string `json:"option"`
	jwt.StandardClaims
}

const (
	HMAC uint8 = 1
	RSA  uint8 = 2
)

var (
	privateKeyParsed *rsa.PrivateKey
	publicKeyParsed  *rsa.PublicKey
	audience         string = "hiromaily.com"
	encrypted        uint8  = 2 //1:HMAC, 2:RSA
	secret           string = "default-secret-key"
)

func init() {
	//log
	lg.InitializeLog(lg.DEBUG_STATUS, lg.LOG_OFF_COUNT, 0, "[JWT]", "")
}

func InitEncrypted(mode uint8) {
	encrypted = mode
}

func InitSecretKey(str string) {
	secret = str
}

func InitAudience(str string) {
	audience = str
}

func InitKeys(priKey, pubKey string) (err error) {
	privateKeyParsed, err = lookupPrivateKey(priKey)
	if err != nil {
		return err
	}

	publicKeyParsed, err = lookupPublicKey(pubKey)

	return
}

func getMethod() jwt.SigningMethod {
	if encrypted == HMAC {
		return jwt.SigningMethodHS256
	} else {
		//RSA
		return jwt.SigningMethodRS256
	}
}

// Payload
func getClaims(t int64, clientId, userName string) jwt.StandardClaims {
	//Audience  string `json:"aud,omitempty"` // https://login.hiromaily.com
	//ExpiresAt int64  `json:"exp,omitempty"`
	//Id        string `json:"jti,omitempty"`
	//IssuedAt  int64  `json:"iat,omitempty"`
	//Issuer    string `json:"iss,omitempty"` // OAuth client_id
	//NotBefore int64  `json:"nbf,omitempty"`
	//Subject   string `json:"sub,omitempty"` // user name or email
	claims := jwt.StandardClaims{
		Audience: audience,
		//ExpiresAt: time.Now().Add(time.Second * 2).Unix(),
		ExpiresAt: t,
		Issuer:    clientId,
		Subject:   userName,
	}
	return claims
}

//Header
//Payload
//encode Header,Payload,Signature by Base64 and concatenate these by dot.
func CreateBasicToken(t int64, clientId, userName string) (string, error) {
	lg.Info("CreateBasicToken()")

	// payload
	claims := getClaims(t, clientId, userName)

	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token := jwt.NewWithClaims(getMethod(), &claims)

	//return token.SignedString([]byte("secret")) //OK
	if encrypted == HMAC {
		return token.SignedString([]byte(secret))
	} else {
		//RSA
		return token.SignedString(privateKeyParsed) //use private key
	}
}

func CreateToken(t int64, clientId, userName, option string) (string, error) {
	lg.Info("CreateToken()")

	// Create the Claims
	// payload
	claims := getClaims(t, clientId, userName)

	cClaims := &CustomClaims{
		option,
		claims,
	}

	//SigningMethodRS256
	token := jwt.NewWithClaims(getMethod(), cClaims)
	if encrypted == HMAC {
		return token.SignedString([]byte(secret))
	} else {
		//RSA
		return token.SignedString(privateKeyParsed) //use private key
	}
}

// judge parse
func judgeParse(token *jwt.Token) (interface{}, error) {
	lg.Info("judgeParse()")

	var ok bool = false
	if encrypted == HMAC {
		_, ok = token.Method.(*jwt.SigningMethodHMAC)
	} else if encrypted == RSA {
		_, ok = token.Method.(*jwt.SigningMethodRSA)
	}

	if !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	if encrypted == HMAC {
		return []byte(secret), nil
		//} else if encrypted == RSA {
	} else {
		//RSA
		return publicKeyParsed, nil
		//return privateKeyParsed, nil //key is of invalid type
	}
}

// Check Token (it may be too strict to check)
func JudgeJWT(tokenString string) error {
	lg.Info("JudgeJWT()")

	//token
	token, err := jwt.Parse(tokenString, judgeParse)

	if err != nil {
		return err
	} else if !token.Valid {
		return fmt.Errorf("token is invalid")
	}

	return nil
}

// Check Token (it may be too strict to check)
func JudgeJWTWithClaim(tokenString, clientId, userName string) error {
	lg.Info("JudgeJWTWithClaim()")

	//token
	//token, err := jwt.Parse(tokenString, judgeParse)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, judgeParse)

	if err != nil {
		return err
	} else if !token.Valid {
		return fmt.Errorf("token is invalid")
	}

	//check claim
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return fmt.Errorf("Claims data can't be retrieved")
	} else if claims.Issuer != clientId {
		return fmt.Errorf("Issuer is invalid")
	} else if claims.Subject != userName {
		return fmt.Errorf("Subject is invalid")
	}

	return nil
}

// Check Token (it may be too strict to check)
func JudgeJWTWithCustomClaim(tokenString, clientId, userName, option string) error {
	lg.Info("JudgeJWTWithCustomClaim()")

	//token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, judgeParse)

	if err != nil {
		return err
	} else if !token.Valid {
		return fmt.Errorf("token is invalid")
	}

	//check claim
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return fmt.Errorf("Claims data can't be retrieved")
	} else if claims.Issuer != clientId {
		return fmt.Errorf("Issuer is invalid")
	} else if claims.Subject != userName {
		return fmt.Errorf("Subject is invalid")
	} else if claims.Option != option {
		return fmt.Errorf("Option is invalid")
	}

	return nil
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
