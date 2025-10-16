package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// RegisterProxyRoutes регистрирует защищённые маршруты, которые проксируют запросы на upstream
func RegisterProxyRoutes(router *gin.Engine) {
	// Middleware AuthRequired() защищает все маршруты в этой группе
	group := router.Group("/", AuthRequired())

	for _, upstream := range appConfig.Upstreams {
		target, err := url.Parse("https://" + upstream.Host)
		if err != nil {
			panic("invalid upstream URL: " + upstream.Host)
		}

		// Проксируем все запросы на этот upstream
		group.Any("/*path", ReverseProxyHandler(target))
	}
}

// ReverseProxyHandler создаёт обработчик для проксирования запроса на указанный target
func ReverseProxyHandler(target *url.URL) gin.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
		// Можно подменить Host заголовок, если нужно
		c.Request.Host = target.Host
		c.Request.URL.Scheme = target.Scheme
		c.Request.URL.Host = target.Host

		// Проксируем
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
