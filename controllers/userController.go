package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Arueljust/initializers"
	"github.com/Arueljust/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// ambil data inputan user "email dan password"
	var body struct {
		Email    string
		Password string
	}
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "filed to read body request !",
		})
		return
	}

	// hash passwrd
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "filed to hash password"})
		return
	}
	// create user
	user := models.User{Email: body.Email, Password: string(hash)}
	res := initializers.DB.Create(&user)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "filed to create user"})
		return
	}
	// response
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	// ambil data inputan dari req body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "filed to read body request !"})
		return
	}
	// lihat req email
	var user models.User
	initializers.DB.First(&user, "email = ? ", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid read email"})
		return
	}
	// compare pass yang dikirim sama password yang ada di db
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "kode email / password "})
		return
	}
	// generate jwt token
	// var createdAt []uint8
	// errTime := initializers.DB.Row().Scan(&createdAt)
	// if errTime != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "filed convert time"})
	// 	return
	// }
	// createdAtStr := string(createdAt)
	// createdTime, err := time.Parse("2023-05-25 04:26:51.405000000", createdAtStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "filed to parse time"})
	// 	return
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokeString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		return
	}
	// kembalikan token
	c.JSON(http.StatusOK, gin.H{"token": tokeString})
}
