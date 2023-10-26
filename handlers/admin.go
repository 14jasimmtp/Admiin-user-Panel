package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	initialisers "main.go/Initialisers"
	"main.go/models"
	"main.go/utils"
)

var AdmintokenString string

func Admin(c *gin.Context) {
	utils.ClearCache(c)
	// AdmintokenString, _ = c.Cookie("Admin")
	// if AdmintokenString == "" {
	// 	c.HTML(http.StatusUnauthorized, "admin-login.html", nil)
	// 	return
	// }
	// c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")
	// return

	var err error

	_, err = c.Cookie("User")
	if err == nil {
		c.Redirect(http.StatusSeeOther, "/")
	}

	_, err = c.Cookie("Admin")
	if err == nil {
		c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")
	}
	c.HTML(http.StatusOK, "admin-login.html", nil)
	return
}

func PostAdmin(c *gin.Context) {

	Admin := models.User{
		Email:    strings.TrimSpace(c.Request.FormValue("admin-email")),
		Password: strings.TrimSpace(c.Request.FormValue("admin-password")),
	}

	var DBUsers models.User
	result := initialisers.DB.Where("email = ?", Admin.Email).First(&DBUsers)
	if result.Error != nil {
		c.HTML(http.StatusBadRequest, "admin-login.html", "Admin not found")
		return
	}

	if !DBUsers.Is_Admin {
		c.HTML(http.StatusBadRequest, "admin-login.html", "Admin not found")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(DBUsers.Password), []byte(Admin.Password))

	if err != nil {
		c.HTML(http.StatusBadRequest, "admin-login.html", "wrong password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  Admin.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err = token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		log.Fatal("failed to create token", err)
	}
	c.SetCookie("Admin", tokenString, 3600, "", "", false, true)

	c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")
}
