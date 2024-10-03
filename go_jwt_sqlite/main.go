package main

import (
	"fmt"

	"github.com/brangb/go_jwt_sqlite/config"
	"github.com/brangb/go_jwt_sqlite/controllers"
	"github.com/brangb/go_jwt_sqlite/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/brangb/go_jwt_sqlite/docs" // Import generated docs
)

// @title Gin Voting API
// @version 1.0
// @description This is a simple API for voting system.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email youremail@provider.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:4000
// @BasePath /

func init() {
	config.LoadEnvVariables()
	config.ConnectToDB()
	config.SyncDatabase()
}

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		fmt.Println("Hello, you'r in the right path")
	})
	r.POST("/signup", controllers.CreateUser)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/createRoom", middleware.CreateVotingRoom, controllers.CreateVotingRoom)
	r.POST("/getRoom", controllers.GetRoomsById)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
	fmt.Println("Server Started")
}
