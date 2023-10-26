package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/utils"
)

func Logout(c *gin.Context) {
	utils.ClearCache(c)
	c.SetCookie("User", "", -1, "", "", false, true)
	c.SetCookie("Admin", "", -1, "", "", false, true)
	// tokenString = ""
	c.Redirect(http.StatusSeeOther, "/login")

	// c.HTML(http.StatusOK, "login.html", "logged out")

}
