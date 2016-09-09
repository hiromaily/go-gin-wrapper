package jsons

// JWTResponse is response of JWT request
type JWTResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

// CreateJWTJson is for response format of JWT
func CreateJWTJson(token string) *JWTResponse {
	jwtJSON := JWTResponse{
		Code:  0,
		Token: token,
	}
	return &jwtJSON
}
