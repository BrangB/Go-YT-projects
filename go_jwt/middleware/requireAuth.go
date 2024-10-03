package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/brangb/go_jwt/initializers"
	"github.com/brangb/go_jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {

	// Retrieve the token from the "Authorization" cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what you expect (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for validation
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	// Handle errors in token parsing
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Verify the token claims (type assertion to jwt.MapClaims)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Check if the token has expired
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	} else {
		// If "exp" claim is not present or invalid
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Fetch the user from the database using the "sub" (subject) claim
	var user models.User
	initializers.DB.First(&user, claims["sub"])

	// If the user does not exist, return unauthorized
	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Set the user object to the context, so it's accessible in subsequent handlers
	c.Set("user", user)

	// Allow the request to continue
	c.Next()
}
