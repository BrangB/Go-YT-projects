package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateVotingRoom(c *gin.Context) {
	fmt.Println("Hello You in the middleware")
	c.Next()
}
