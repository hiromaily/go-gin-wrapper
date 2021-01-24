package json

// JWTResponse is response of JWT request
type JWTResponse struct {
	Token string `json:"token"`
}

// CreateJWTJson is response of JWT token
func CreateJWTJson(token string) *JWTResponse {
	return &JWTResponse{
		Token: token,
	}
}
