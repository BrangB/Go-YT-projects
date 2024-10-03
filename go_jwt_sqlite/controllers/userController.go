package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/brangb/go_jwt_sqlite/config"
	"github.com/brangb/go_jwt_sqlite/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "User Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /signup [post]
func CreateUser(c *gin.Context) {
	var Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read data",
		})
		return
	}

	if Body.Email == "" || Body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email or password cannot be empty",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{Email: Body.Email, Password: string(hash)}
	result := config.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login godoc
// @Summary Log in a user
// @Description Log in with email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body models.User true "User Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the data",
		})
		return
	}

	var user models.User
	config.DB.First(&user, "email = ?", Body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Token expires in 30 days
	})

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

// Validate godoc
// @Summary Validate a user token
// @Description Validate the JWT of a logged-in user
// @Tags Auth
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router /validate [get]
func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

// CreateVotingRoom godoc
// @Summary Create a new voting room
// @Description Create a room for voting with title and description
// @Tags VotingRoom
// @Accept  json
// @Produce  json
// @Param room body models.Room true "Room Data"
// @Success 200 {object} models.Room
// @Failure 400 {object} map[string]string
// @Router /createRoom [post]
func CreateVotingRoom(c *gin.Context) {
	var input struct {
		OwnerID     uint   `json:"owner_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create room instance
	room := models.Room{
		OwnerID:     input.OwnerID,
		Title:       input.Title,
		Description: input.Description,
	}

	// Save to database
	if err := config.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create voting room"})
		return
	}

	// Respond with the created room
	c.JSON(http.StatusOK, gin.H{"room": room})
}

// GetRoomsById godoc
// @Summary Get voting rooms by owner ID
// @Description Retrieve all voting rooms created by a specific owner
// @Tags VotingRoom
// @Accept  json
// @Produce  json
// @Param owner_id body uint true "Owner ID"
// @Success 200 {object} []models.Room
// @Failure 500 {object} map[string]string
// @Router /getRoom [post]
func GetRoomsById(c *gin.Context) {
	var Body struct {
		OwnerID uint `json:"owner_id"`
	}

	// Bind the request JSON to Body
	if err := c.ShouldBindJSON(&Body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var rooms []models.Room
	// Fetch all rooms where the OwnerID matches
	if err := config.DB.Where("owner_id = ?", Body.OwnerID).Find(&rooms).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	// Return the rooms as JSON
	c.JSON(200, gin.H{
		"data": rooms,
	})
}
