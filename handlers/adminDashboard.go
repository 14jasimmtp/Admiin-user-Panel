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
	"gorm.io/gorm"
	initialisers "main.go/Initialisers"
	"main.go/models"
	"main.go/utils"
)

var DataMap = make(map[string]any)
var Users []models.User

var result *gorm.DB

func AdminDashboard(c *gin.Context) {
	utils.ClearCache(c)
	//check if user is created
	var err error

	result = initialisers.DB.Order("created_at ASC").Where("is_admin = false").Find(&Users)
	if result.Error != nil {
		fmt.Println("Error fetching users:", result.Error)
		c.HTML(http.StatusOK, "admin-dashboard.html", Users)
		return
	}
	ser := make([]models.Serial, int(result.RowsAffected))
	for i := 0; i < int(result.RowsAffected); i++ {
		ser[i].Ser = i + 1
	}

	DataMap["UsersList"] = Users
	DataMap["serial"] = ser

	// if AdmintokenString != "" {

	// 	_, err := c.Cookie("Admin")
	// 	if err != nil {
	// 		c.Redirect(http.StatusSeeOther, "/")
	// 		return
	// 	}

	// 	c.HTML(http.StatusOK, "admin-dashboard.html", DataMap)

	// } else {
	// 	c.Redirect(http.StatusSeeOther, "/admin")
	// }
	_, err = c.Cookie("User")
	if err == nil {
		c.Redirect(http.StatusSeeOther, "/")
	}

	_, err = c.Cookie("Admin")
	if err == nil {

		c.HTML(http.StatusOK, "admin-dashboard.html", DataMap)
	}
	return

}

func GetCreateUser(c *gin.Context) {
	c.HTML(http.StatusOK, "CreateUser.html", nil)

}

func CreateUser(c *gin.Context) {
	NewUser := models.User{
		UserName: strings.TrimSpace(c.Request.FormValue("username")),
		Email:    strings.TrimSpace(c.Request.FormValue("email")),
		Password: strings.TrimSpace(c.Request.FormValue("password")),
	}
	ConfirmPwd := strings.TrimSpace(c.Request.FormValue("cfmpassword"))
	if ConfirmPwd != NewUser.Password {
		c.HTML(http.StatusNotAcceptable, "CreateUser.html", "Passwords does not match ")
		return
	}

	NewUser.Password = utils.HashPwd(NewUser.Password)

	CreateUser := initialisers.DB.Create(&NewUser)

	if CreateUser.Error != nil {
		c.HTML(http.StatusNotAcceptable, "CreateUser.html", "Something went wrong")
		return
	}
	ser := make([]models.Serial, 100)
	for i := 0; i < int(result.RowsAffected); i++ {
		ser[i].Ser = i + 1
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
	c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")
	// c.HTML(http.StatusOK, "admin-dashboard.html", DataMap)

}

func GetUpdateUser(c *gin.Context) {
	email := c.Request.FormValue("Updateemail")
	fmt.Println(email, "email")

	// Fetch user from the db
	var User models.User
	res := initialisers.DB.Raw(`SELECT * FROM users WHERE email = ?`, email).Scan(&User)
	if res.Error != nil {
		DataMap["Message"] = "Error fetching user"
		fmt.Println("Error is ", res.Error)
		c.HTML(http.StatusOK, "UpdateUser.html", DataMap)
		return
	}

	// Pass the data to show in update page
	DataMap["Users"] = User
	fmt.Println(DataMap["Users"])
	c.HTML(http.StatusOK, "UpdateUser.html", DataMap)
}

func UpdateUser(c *gin.Context) {
	email := c.Request.FormValue("email")
	fmt.Println(email, "email")

	user := models.User{
		UserName: strings.TrimSpace(c.Request.FormValue("Username")),
		Email:    strings.TrimSpace(c.Request.FormValue("Updateemail")),
		Model:    gorm.Model{UpdatedAt: time.Now()},
	}
	fmt.Println(user.UserName, "name", user.Email, "email")

	var currentUserEmail string
	result := initialisers.DB.Raw("SELECT email FROM users WHERE email = ?", user.Email).Scan(&currentUserEmail)
	if result.Error != nil {
		DataMap["Message"] = "Error updating user"
		fmt.Println("error")
		c.HTML(http.StatusOK, "UpdateUser.html", DataMap)
		return
	}

	if result.RowsAffected > 0 && user.Email != currentUserEmail {
		DataMap["Message"] = "Already user exist with this email"
		fmt.Println("user already exist")
		c.HTML(http.StatusOK, "UpdateUser.html", DataMap)
		return
	}

	result = initialisers.DB.Exec(`UPDATE users SET email = ?, user_name = ? WHERE email = ?`, user.Email, user.UserName, email)
	if result.Error != nil {
		DataMap["Message"] = "Error updating user"
		c.HTML(http.StatusBadRequest, "UpdateUser.html", DataMap)
		return
	}

	// Pass success Message
	DataMap["Message"] = "User updated succesfully"
	c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")

}

func DeleteUser(c *gin.Context) {
	email := c.Request.FormValue("email")
	result := initialisers.DB.Exec(`DELETE FROM users WHERE email = ?`, email)
	if result.Error != nil {
		fmt.Println("Error while deleting user", result.Error)
		return
	}
	ser := make([]models.Serial, 1)
	for i := 0; i < int(result.RowsAffected); i++ {
		ser[i].Ser = i + 1
	}

	c.Redirect(http.StatusSeeOther, "/Admin-Dashboard")
	// c.HTML(http.StatusOK, "admin-dashboard.html", DataMap)

}

func SearchUser(c *gin.Context) {

	Search := "%" + c.Request.FormValue("search") + "%"
	fmt.Println(Search)
	result := initialisers.DB.Raw(`SELECT * FROM users WHERE user_name ILIKE ?`, Search).Scan(&Users)
	ser := make([]models.Serial, int(result.RowsAffected))
	for i := 0; i < int(result.RowsAffected); i++ {
		ser[i].Ser = i + 1
	}
	DataMap["UsersList"] = Users
	DataMap["serial"] = ser

	c.HTML(http.StatusOK, "admin-dashboard.html", DataMap)
}
