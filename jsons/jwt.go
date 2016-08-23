package jsons

type JWTResponse struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

func CreateJWTJson(token string) *JWTResponse {
	jwtJson := JWTResponse{
		Code:  0,
		Token: token,
	}
	return &jwtJson
}
