package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Инициализируем Valkey
	InitValkey()

	// 2. Загружаем конфиг
	appConfig = LoadConfig()

	// 3. Auth сервер
	go func() {
		auth := gin.Default()
		auth.SetTrustedProxies([]string{"127.0.0.1"})
		RegisterAuthHandlers(auth)
		// Простой тестовый эндпоинт
		auth.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "auth.secure-proxy.lan")
		})

		// Генерация и проверка TOTP
		auth.GET("/totp/generate", GenerateTOTPHandler)
		auth.POST("/totp/validate", ValidateTOTPHandler)
		auth.StaticFile("/login.html", "./login.html")

		// HTTPS сервер для auth
		auth.RunTLS(":8443", "_.secure-proxy.lan.crt", "_.secure-proxy.lan.pem")
	}()

	// 4. REST сайт (демо)
	rest := gin.Default()
	rest.SetTrustedProxies([]string{"127.0.0.1"})
	rest.Use(AuthRequired())

	rest.GET("/set-cookie", func(c *gin.Context) {
		c.SetCookie("test-cookie", "hello-world", 3600, ".secure-proxy.lan", "", false, true)
		c.String(http.StatusOK, "Cookie установлена!")
	})

	rest.GET("/get-cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("test-cookie")
		if err != nil {
			c.String(http.StatusOK, "Cookie не найдена")
			return
		}
		c.String(http.StatusOK, "Cookie значение: %s", cookie)
	})

	// теперь все маршруты защищены

	rest.GET("/", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.String(http.StatusOK, "site1.secure-proxy.lan\n Добро пожаловать, %v", username)
	})
	//RegisterProxyRoutes(rest)
	// HTTPS сервер для demo сайта
	rest.RunTLS(":9443", "_.secure-proxy.lan.crt", "_.secure-proxy.lan.pem")
	proxy := gin.Default()
	proxy.SetTrustedProxies([]string{"127.0.0.1"})
	RegisterProxyRoutes(proxy)
	proxy.RunTLS(":10443", "_.secure-proxy.lan.crt", "_.secure-proxy.lan.pem")

}

// rest.GET("/", func(c *gin.Context) {
// 	c.String(http.StatusOK, "site1.secure-proxy.lan")
// })
