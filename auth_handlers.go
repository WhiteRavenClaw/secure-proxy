// файл auth_handlers.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterAuthHandlers(auth *gin.Engine) {
	auth.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		code := c.PostForm("totp")
		redirectUrl := c.PostForm("redirectUrl")

		if !ValidateUserTOTP(username, code) {
			c.String(http.StatusUnauthorized, "Invalid username or TOTP code")
			return
		}

		sessionKey := uuid.New().String()
		ValkeySet(sessionKey, username, appConfig.Sessions.TTLSeconds)

		c.SetCookie(
			appConfig.Sessions.CookieName,
			sessionKey,
			appConfig.Sessions.TTLSeconds,
			"/",
			appConfig.Sessions.CookieDomain,
			true,
			true,
		)

		c.Redirect(http.StatusFound, redirectUrl)
	})
}
