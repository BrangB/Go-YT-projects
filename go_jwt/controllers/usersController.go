package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/brangb/go_jwt/initializers"
	"github.com/brangb/go_jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Signup handles user registration
func Signup(c *gin.Context) {
	var Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind request body to the struct
	if err := c.Bind(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read Body",
		})
		return
	}

	// Hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create a new user record
	user := models.User{Email: Body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	// Handle errors during user creation
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login handles user login and token generation
func Login(c *gin.Context) {
	var Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind request body to the struct
	if err := c.Bind(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the data",
		})
		return
	}

	// Find the user by email
	var user models.User
	initializers.DB.First(&user, "email = ?", Body.Email)

	// If user is not found
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// Compare the provided password with the stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
		return
	}

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Token expires in 30 days
	})

	// Sign the token with the secret key (converted to byte slice)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
