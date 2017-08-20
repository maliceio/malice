package http

import (
	"os"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/maliceio/engine/api"
)

var (
	defaultUser string
	defaultPass string
	port        string
)

// Token is an API key used for auth
type Token struct {
	Key string `form:"key" json:"key" binding:"required"`
}

func getOpt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

func init() {
	defaultUser = getOpt("USER", "admin")
	defaultPass = getOpt("PASS", "admin")
	port = getOpt("PORT", "8080")
}

// StartHTTP start http server
func StartHTTP() {

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "malice-server",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			if userId == defaultUser && password == defaultPass {
				return userId, true
			}

			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == defaultUser {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		TokenLookup: "query:token",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", api.HelloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	endless.ListenAndServe(":"+port, r) // listen and serve on 0.0.0.0:8080
}
