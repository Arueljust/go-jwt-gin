package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Arueljust/initializers"
	"github.com/Arueljust/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(c *gin.Context) {
	// ambil data cookie dari request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// parse token string dan buat fungsi untuk mencari kunci
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validasi token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// ambil secret key type data []byte
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check expired token , dengan cara cek "exp"
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// check user dengan token "sub"
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// masukkan di request
		c.Set("user", user)
		// lanjutkan
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
