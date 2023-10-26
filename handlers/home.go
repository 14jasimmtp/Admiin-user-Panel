package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	initialisers "main.go/Initialisers"
	"main.go/models"
	"main.go/utils"
)

func ValidateUser(c *gin.Context) {
	tokenString, err := c.Cookie("User")

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRETKEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var User models.User
		initialisers.DB.Raw("SELECT * FROM users WHERE email = ? ", claims["Email"]).Scan(&User)

		c.Set("User", User.UserName)
		fmt.Println(User.UserName)
		c.Next()

	} else {
		fmt.Println(err)
	}

}

func Home(c *gin.Context) {
	utils.ClearCache(c)
	//check if user is created
	tokenString, err := c.Cookie("User")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	if tokenString != "" {

		_, err := c.Cookie("User")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")
			return
		}
		fmt.Println("name", c.GetString("User"))
		c.HTML(http.StatusOK, "index.html", c.GetString("User"))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/login")
	return
}
