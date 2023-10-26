package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	initialisers "main.go/Initialisers"
	"main.go/models"
	"main.go/utils"
)

func Login(c *gin.Context) {
	utils.ClearCache(c)

	var err error

	_, err = c.Cookie("User")
	if err == nil {
		c.Redirect(http.StatusSeeOther, "/")
	}

	_, err = c.Cookie("Admin")
	if err == nil {
		c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")
	}
	c.HTML(http.StatusOK, "login.html", nil)
	return

}

func PostLogin(c *gin.Context) {

	User := models.User{
		Email:    c.Request.FormValue("email"),
		Password: c.Request.FormValue("password"),
	}
	fmt.Println(User.Email)
	var DBUsers models.User
	result := initialisers.DB.Where("email = ?", User.Email).First(&DBUsers)

	// result := initialisers.DB.Raw("SELECT * FROM customes WHERE email = ?", User.Email).Scan(&)
	// result := initialisers.DB.First(&DBUsers, "email = ?", User.Email)
	fmt.Println(DBUsers)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "login.html", "User doesn't exist")
		fmt.Print("some problem")
		return
	}
	if result.RowsAffected != 1 {
		c.HTML(http.StatusBadRequest, "login.html", "User doesn't exist")
		fmt.Println("uswe don't exist")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(DBUsers.Password), []byte(User.Password))
	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", "Wrong password")
		fmt.Print("wrong password")
		return
	}
	fmt.Println("user created with id : ", User.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email": User.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		log.Fatal("failed to create token", err)
	}
	c.SetCookie("User", tokenString, 3600, "", "", false, true)

	c.Redirect(http.StatusSeeOther, "/")
	// c.HTML(http.StatusOK, "index.html", "succefully logged in")
	fmt.Println("redirected to home page")
}
