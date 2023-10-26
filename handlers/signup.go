package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	initialisers "main.go/Initialisers"
	"main.go/models"
	"main.go/utils"
)

var tokenString string
var Name string

func Signup(c *gin.Context) {
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
	c.HTML(http.StatusOK, "signup.html", nil)
	return

}

func PostSignup(c *gin.Context) {
	NewUser := models.User{
		UserName: strings.TrimSpace(c.Request.FormValue("name")),
		Email:    strings.TrimSpace(c.Request.FormValue("email")),
		Password: strings.TrimSpace(c.Request.FormValue("password")),
	}
	ConfirmPwd := strings.TrimSpace(c.Request.FormValue("confirmpwd"))
	if ConfirmPwd != NewUser.Password {
		c.HTML(http.StatusNotAcceptable, "signup.html", "Passwords does not match ")
		return
	}

	NewUser.Password = utils.HashPwd(NewUser.Password)

	CreateUser := initialisers.DB.Create(&NewUser)

	if CreateUser.Error != nil {
		c.HTML(http.StatusNotAcceptable, "signup.html", "Something went wrong")
		return
	}

	fmt.Println("user created with id : ", NewUser.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  NewUser.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	var err error
	tokenString, err = token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		log.Fatal("failed to create token", err)
	}
	c.SetCookie("User", tokenString, 3600, "", "", false, true)

	c.Set("name", NewUser.UserName)
	Name = c.GetString("name")
	fmt.Println(NewUser.UserName)
	c.Redirect(http.StatusSeeOther, "/")
	// c.Redirect(http.StatusSeeOther, "/")

}
