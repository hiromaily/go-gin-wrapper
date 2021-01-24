package ginsession

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-gin-wrapper/pkg/token"
)

// Sessioner interface
type Sessioner interface {
	IsLogin(ctx *gin.Context) (bool, int)
	Logout(ctx *gin.Context)
	SetUserID(ctx *gin.Context, userID int)
	GenerateToken() string
	SetToken(ctx *gin.Context, token string)
	GetToken(ctx *gin.Context) string
	DeleteToken(ctx *gin.Context)
	ValidateToken(ctx *gin.Context, token string) error
}

type sessioner struct {
	logger *zap.Logger
	token  token.Generator
}

// NewSessioner returns Sessioner
func NewSessioner(logger *zap.Logger, token token.Generator) Sessioner {
	return &sessioner{
		logger: logger,
		token:  token,
	}
}

// IsLogin returns login status of boolean and uid
func (s *sessioner) IsLogin(ctx *gin.Context) (bool, int) {
	session := sessions.Default(ctx)
	v := session.Get("uid")
	if v == nil {
		return false, 0
	}
	return true, v.(int)
}

// Logout clears session
func (s *sessioner) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}

// SetUserID sets userID as session
// - skip if session has uid
func (s *sessioner) SetUserID(ctx *gin.Context, userID int) {
	session := sessions.Default(ctx)
	v := session.Get("uid")
	if v == nil {
		session.Set("uid", userID)
		session.Save()
	}
	s.logger.Debug("SetUserID has already been set")
}

// GenerateToken generates token
func (s *sessioner) GenerateToken() string {
	return s.token.Generate()
}

// SetToken sets token as session
func (s *sessioner) SetToken(ctx *gin.Context, token string) {
	session := sessions.Default(ctx)
	session.Set("token", token)
	session.Save()
}

// GetToken returns session token
func (s *sessioner) GetToken(ctx *gin.Context) string {
	session := sessions.Default(ctx)
	v := session.Get("token")
	if v == nil {
		return ""
	}
	return v.(string)
}

// DeleteToken deletes session token
func (s *sessioner) DeleteToken(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete("token")
	session.Save()
}

// ValidateToken validates token
func (s *sessioner) ValidateToken(ctx *gin.Context, token string) error {
	sesToken := s.GetToken(ctx)
	s.logger.Info("IsTokenValid",
		zap.String("GetToken()", sesToken),
		zap.String("given_token", token),
	)

	if sesToken != "" && token != "" {
		return nil
	}

	var err error
	if token == "" {
		err = errors.New("given token is empty")
	} else if sesToken == "" {
		err = errors.New("session token is missing, might have been expired")
	} else if sesToken != token {
		err = errors.New("token is invalid")
	}
	s.DeleteToken(ctx)
	return err
}
