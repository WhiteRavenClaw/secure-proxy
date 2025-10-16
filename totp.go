package main

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

var appConfig *AppConfig

type TOTPResponse struct {
	Secret    string `json:"secret"`
	QRCodeURL string `json:"qr_code_url"`
	QRCodePNG string `json:"qr_code_png"`
}

// Генерация нового TOTP-секрета и QR
func GenerateTOTPHandler(c *gin.Context) {
	issuer := c.DefaultQuery("issuer", "Secure Proxy")
	accountName := c.DefaultQuery("account_name", "auth@secure-proxy.lan")

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации TOTP ключа"})
		return
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации QR изображения"})
		return
	}

	if err := png.Encode(&buf, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при кодировании PNG"})
		return
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	response := TOTPResponse{
		Secret:    key.Secret(),
		QRCodeURL: key.URL(),
		QRCodePNG: qrCodeBase64,
	}

	c.JSON(http.StatusOK, response)
}

// Проверка TOTP и установка сессии в Valkey
func ValidateTOTPHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// 1. Найти пользователя в конфиге
	var user *UserConfig
	for _, u := range appConfig.Users {
		if u.Username == req.Username {
			user = &u
			break
		}
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unknown user"})
		return
	}

	// 2. Проверить TOTP
	if !totp.Validate(req.Code, user.TOTPSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP"})
		return
	}

	// 3. Сгенерировать session ID
	sessionKey := uuid.NewString()

	// 4. Сохранить сессию в Valkey
	err := ValkeySet("session_"+sessionKey, req.Username, appConfig.Sessions.TTLSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	// 5. Установить cookie
	c.SetCookie(
		appConfig.Sessions.CookieName,
		sessionKey,
		appConfig.Sessions.TTLSeconds,
		//"/",
		".secure-proxy.lan",
		appConfig.Sessions.CookieDomain,
		true, // secure
		true, // httpOnly
	)

	// 6. Редиректить обратно
	// redirectURL := c.DefaultQuery("redirectUrl", "https://"+appConfig.Upstreams[0].Host)
	// c.Redirect(http.StatusFound, redirectURL)
	redirectURL := c.PostForm("redirectUrl")
	if redirectURL == "" {
		redirectURL = "https://site1.secure-proxy.lan:9443/"
	}
	c.Redirect(http.StatusFound, redirectURL)

}

// Общая функция для проверки кода и пользователя
func ValidateUserTOTP(username, code string) bool {
	for _, u := range appConfig.Users {
		if u.Username == username {
			return totp.Validate(code, u.TOTPSecret)
		}
	}
	return false
}
