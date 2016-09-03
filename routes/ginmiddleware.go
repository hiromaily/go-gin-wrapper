package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	conf "github.com/hiromaily/go-gin-wrapper/configs"
	sess "github.com/hiromaily/go-gin-wrapper/libs/ginsession"
	hh "github.com/hiromaily/go-gin-wrapper/libs/httpheader"
	"github.com/hiromaily/golibs/auth/jwt"
	lg "github.com/hiromaily/golibs/log"
	"net/http"
	"strings"
)

//TODO: it's not finished yet.
func CheckHttpRefererAndCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckHttpRefererAndCSRF]")
		//Referer
		url := hh.GetUrl(c)
		//get referer data mapping table using url (map[string])
		if refUrl, ok := RefererUrls[url]; ok {
			//Check Referer
			if !hh.IsRefererHostValid(c, refUrl) {
				c.Next()
				return
			}
		}

		//CSRF
		sess.IsTokenSessionValid(c, c.PostForm("gintoken"))
		c.Next()
	}
}

//Check Http Referer.
//TODO: it's not checked yet if it work well or not.
func CheckHttpReferer() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[heckHttpReferer]")
		url := hh.GetUrl(c)
		//get referer data mapping table using url (map[string])
		if refUrl, ok := RefererUrls[url]; ok {
			//Check Referer
			hh.IsRefererHostValid(c, refUrl)
		}
		c.Next()
	}
}

//TODO: it's not finished yet.
func CheckCSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckCSRF]")
		sess.IsTokenSessionValid(c, c.PostForm("gintoken"))
		c.Next()
	}
}

//Check Http Header for Ajax request. (For REST)
func CheckHttpHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckHttpHeader]")

		apiConf := conf.GetConf().Auth.Api

		lg.Debugf("[Request Header]\n%#v\n", c.Request.Header)
		lg.Debugf("[Request Form]\n%#v\n", c.Request.Form)
		lg.Debugf("[Request Body]\n%#v\n", c.Request.Body)
		lg.Debugf("c.ClientIP() %s", c.ClientIP())

		//IsAjax := c.Request.Header.Get("X-Requested-With")
		//lg.Debugf("[X-Requested-With] %s", IsAjax)

		//IsKey := c.Request.Header.Get("X-Custom-Header-Gin")
		//lg.Debugf("[X-Custom-Header-Gin] %s", IsKey)

		IsKey := c.Request.Header.Get(apiConf.Header)
		lg.Debugf("[%s] %s", apiConf.Header, IsKey)

		IsContentType := c.Request.Header.Get("Content-Type")
		lg.Debugf("[Content-Type] %s", IsContentType)

		//TODO:if no Content-Type, how sould be handled.

		//check
		//if IsXHR(c) || IsKey != "key" || IsContentType != "application/json" {
		//if IsXHR(c) && IsKey != "key" {
		if (apiConf.Ajax && !IsXHR(c)) || IsKey != apiConf.Key {
			//error
			c.AbortWithStatus(400)
			return
		}

		//Context Meta Data
		//SetMetaData(c)

		c.Next()
	}
}

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[CheckJWT]")

		var err error

		IsAuth := c.Request.Header.Get("Authorization")
		if IsAuth != "" {
			aAry := strings.Split(IsAuth, " ")
			if len(aAry) != 2 {
				err = errors.New("Authorization header is invalid")
			} else {
				if aAry[0] != "Bearer" {
					err = errors.New("Authorization header is invalid")
				} else {
					token := aAry[1]
					err = jwt.JudgeJWT(token)
				}
			}
		} else {
			err = errors.New("Authorization header was missed.")
		}

		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		c.Next()
	}
}

//Update user session expire.
//TODO:When session has already started, update session expired
func UpdateUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		if bRet, uid := sess.IsLogin(c); bRet {
			sess.SetUserSession(c, uid)
		}
		c.Next()
	}
}

//Set Meta Data
func SetMetaData() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[SetMetaData]")

		//Context Meta Data
		//http.Header{
		// "Referer":[]string{"http://localhost:9999/"},
		// "Accept-Language":[]string{"ja,en-US;q=0.8,en;q=0.6,de;q=0.4,nl;q=0.2"},
		// "X-Frame-Options":[]string{"deny"},
		// "Content-Security-Policy":[]string{"default-src 'none'"},
		// "X-Xss-Protection":[]string{"1, mode=block"},
		// "Connection":[]string{"keep-alive"},
		// "Accept":[]string{"application/json, text/javascript, */*; q=0.01"},
		// "User-Agent":[]string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36"},
		// "X-Content-Type-Options":[]string{"nosniff"},
		// "X-Requested-With":[]string{"XMLHttpRequest"},
		// "X-Custom-Header-Gin":[]string{"key"},
		// "Content-Type":[]string{"application/x-www-form-urlencoded"},
		// "Accept-Encoding":[]string{"gzip, deflate, sdch"}}

		//Ajax
		if IsXHR(c) {
			c.Set("ajax", "1")
		} else {
			c.Set("ajax", "0")
		}

		//Response Data
		if IsAcceptHeaderJson(c) {
			c.Set("responseData", "json")
		} else {
			c.Set("responseData", "html")
		}

		//Requested Data
		if IsContentTypeJson(c) {
			c.Set("requestData", "json")
		} else {
			c.Set("requestData", "data")
		}

		//User Agent
		c.Set("userAgent", GetUserAgent(c))

		//Language
		c.Set("language", GetLanguage(c))

		c.Next()
	}
}

//After request, handle aborted code or 500 error.
//When 404 or 405 error occurred, response already been set in controllers/errors/errors.go
func GlobalRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(c *gin.Context) {
			lg.Info("[GlobalRecover] defer func()")
			var errMsg string

			refUrl := "/"
			if c.Request.Header.Get("Referer") != "" {
				refUrl = c.Request.Header.Get("Referer")
			}

			if c.IsAborted() {
				lg.Debug("[GlobalRecover] c.IsAborted() is true")
				if c.Errors != nil {
					if c.Errors.Last() != nil {
						errMsg = c.Errors.Last().Err.Error()
					}
				}
				if errMsg == "" {
					switch c.Writer.Status() {
					case 400:
						errMsg = http.StatusText(http.StatusBadRequest)
					case 401:
						errMsg = http.StatusText(http.StatusUnauthorized)
					case 403:
						errMsg = http.StatusText(http.StatusForbidden)
					case 404:
						errMsg = http.StatusText(http.StatusNotFound)
					case 405:
						errMsg = http.StatusText(http.StatusMethodNotAllowed)
					case 406:
						errMsg = http.StatusText(http.StatusNotAcceptable)
					case 407:
						errMsg = http.StatusText(http.StatusProxyAuthRequired)
					case 408:
						errMsg = http.StatusText(http.StatusRequestTimeout)
					case 500:
						errMsg = http.StatusText(http.StatusInternalServerError)
					default:
						errMsg = "something error is happend."
					}
				}

				if IsXHR(c) {
					c.JSON(c.Writer.Status(), gin.H{
						"code": fmt.Sprintf("%d", c.Writer.Status()),
						//"error": c.Errors.Last().Err.Error(), //it caused error because of nil of objct
						"error": errMsg,
					})
				} else {
					c.HTML(c.Writer.Status(), "pages/errors/error.tmpl", gin.H{
						"code":    fmt.Sprintf("%d", c.Writer.Status()),
						"message": errMsg,
						"url":     refUrl,
					})
				}
				return
			}

			if conf.GetConf().Develop.RecoverEnable {
				if rec := recover(); rec != nil {
					lg.Debugf("[GlobalRecover] recover() is not nil:\n %v", rec)
					//TODO:How should response data be decided whether html or json?
					//TODO:Ajax or not doesn't matter to response. HTTP header of Accept may be better.
					//TODO:How precise should I follow specifications of HTTP header.

					if IsXHR(c) {
						c.JSON(http.StatusInternalServerError, gin.H{
							"code":  "500",
							"error": rec,
						})
					} else {
						c.HTML(http.StatusInternalServerError, "pages/errors/error.tmpl", gin.H{
							"code":    "500",
							"message": rec,
							"url":     refUrl,
						})
					}
					return
				}
			}
		}(c)

		c.Next()
		//Next is [Main gin Recovery]
	}
}

//Reject specific IP.
//TODO: working in progress yet.
//TODO: it reject all without reverse　proxy ip address.
//TODO:check registered IP address to reject
func RejectSpecificIp() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[RejectSpecificIp]")
		ip := c.ClientIP()
		//proxy
		if conf.GetConf().Proxy.Mode != 0 {
			if conf.GetConf().Proxy.Server.Host != ip {
				c.AbortWithStatus(403)
				return
			}
		}
		//if ip != "127.0.0.1" {
		//	c.AbortWithStatus(403)
		//}
		c.Next()
	}
}

//Reject non HTTPS access.
//TODO: it's not fixed yet.
func RejectNonHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg.Info("[RejectNonHTTPS]")

		//TODO:check protocol of url
		//if strings.Index(c.url, "https://") == -1 {
		//	c.AbortWithStatus(403)
		//}
		c.Next()
	}
}

//Is this request Ajax or not
func IsXHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}

//Is this request to require JSON or not
func IsAcceptHeaderJson(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Accept"))
	return strings.Index(accept, "application/json") != -1
}

//Is this request including JSON as parameter or not
func IsContentTypeJson(c *gin.Context) bool {
	accept := strings.ToLower(c.Request.Header.Get("Content-Type"))
	return strings.Index(accept, "application/json") != -1
}

//TODO:Get User Agent
func GetUserAgent(c *gin.Context) string {
	// "User-Agent":[]string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36
	// (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36"},
	tmpUserAgent := c.Request.Header.Get("User-Agent")
	return tmpUserAgent
}

//Get Language As Top priority
func GetLanguage(c *gin.Context) string {
	// "Accept-Language":[]string{"ja,en-US;q=0.8,en;q=0.6,de;q=0.4,nl;q=0.2"},
	tmpLanguage := c.Request.Header.Get("Accept-Language")
	if tmpLanguage != "" {
		return strings.Split(tmpLanguage, ",")[0]
	}
	return ""
}
