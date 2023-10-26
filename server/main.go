package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	initialisers "main.go/Initialisers"
	"main.go/handlers"
)

func init() {
	initialisers.LoadEnvVariables()
	initialisers.DbConnect()
	initialisers.SyncDatabase()
}

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("/home/jasim/ecommerce authentication-admin/site/pages/*.html")
	r.Static("/home/jasim/ecommerce authentication-admin/site/Static", "./site/Static")
	r.GET("/", handlers.ValidateUser, handlers.Home)
	r.GET("/signup", handlers.Signup)
	r.POST("/signup", handlers.PostSignup)
	r.GET("/login", handlers.Login)
	r.POST("/login", handlers.PostLogin)
	r.POST("/logout", handlers.Logout)
	r.GET("/admin", handlers.Admin)
	r.POST("/admin", handlers.PostAdmin)
	r.GET("/Admin-Dashboard", handlers.AdminDashboard)
	r.GET("/UpdateUser", handlers.GetUpdateUser)
	r.POST("/UpdateUser", handlers.UpdateUser)
	r.POST("/DeleteUser", handlers.DeleteUser)
	r.POST("/CreateUser", handlers.CreateUser)
	r.GET("/CreateUser", handlers.GetCreateUser)

	fmt.Println("listening and serving on http://localhost:8080")

	r.Run(":8080") // listen and serve localhost:8080

}
