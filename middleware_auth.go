package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Middleware проверки авторизации
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Читаем куку
		sessionKey, err := c.Cookie(appConfig.Sessions.CookieName)
		if err != nil || sessionKey == "" {
			redirectToAuth(c)
			c.Abort()
			return
		}

		// 2. Проверяем наличие в Valkey
		username, err := ValkeyGet("session_" + sessionKey)
		if err != nil || username == "" {
			redirectToAuth(c)
			c.Abort()
			return
		}

		// 3. Обновляем TTL
		_ = ValkeyExpire("session_"+sessionKey, appConfig.Sessions.TTLSeconds)

		// 4. Добавляем имя пользователя в контекст
		c.Set("username", username)

		// Продолжаем обработку
		c.Next()
	}
}

// func redirectToAuth(c *gin.Context) {
// 	authenticatedRedirectUrl := "https://" + c.Request.Host + c.Request.RequestURI
// 	authUrl := "https://" + appConfig.AuthDomain + "/?redirectUrl=" + url.QueryEscape(authenticatedRedirectUrl)
// 	c.Redirect(http.StatusFound, authUrl)
// }

func redirectToAuth(c *gin.Context) {
	target := "https://auth.secure-proxy.lan:8443/login.html?redirectUrl=" +
		url.QueryEscape("https://site1.secure-proxy.lan:9443"+c.Request.URL.Path)
	c.Redirect(http.StatusFound, target)
}

// вспомогательная функция редиректа
//
//	func redirectToAuth(c *gin.Context) {
//		redirectURL := "https://auth.secure-proxy.lan:8443/?redirectUrl=" +
//			url.QueryEscape(c.Request.URL.String())
//		c.Redirect(http.StatusFound, redirectURL)
//	}
