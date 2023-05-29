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
			"error": "filed to read body request !",
		})
		return
	}

	// hash passwrd
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filed to hash password"})
		return
	}
	// create user
	user := models.User{Email: body.Email, Password: string(hash)}
	res := initializers.DB.Create(&user)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filed to create user"})
		return
	}
	// response
	c.JSON(http.StatusOK, gin.H{"error": "User created success"})
}

func Login(c *gin.Context) {
	// ambil data inputan dari req body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filed to read body request !"})
		return
	}
	// lihat req email
	var user models.User
	initializers.DB.First(&user, "email = ? ", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid read email"})
		return
	}

	// compare pass yang dikirim sama password yang ada di db
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email / password ",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokeString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return
	}

	// kembalikan token
	// set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokeString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login sukses"})
}

/*cek user login dengan validasi*/
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"data": user})
}
