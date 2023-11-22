package Controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const SecretKey = "secret"
const SecretKey2 = "secret"

func (h *Handler) RegisterClients(c *gin.Context) {
	var data map[string]string
	errr := c.ShouldBindJSON(&data)
	if errr != nil {
		panic(errr)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	randomUserID := uint(rand.Intn(10000) + 1)
	// Create new user
	user := &Client{
		UserID:   uint64(randomUserID),
		Username: data["username"],
		Email:    data["email"],
		Password: password,
	}
	result := UDB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) LoginClients(c *gin.Context) {
	var data map[string]string
	err := c.ShouldBindJSON(&data)
	if err != nil {
		panic(err)
	}

	var user Client
	RDB.Where("email = ?", data["email"]).First(&user)
	if user.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"Message": "User Not Found"})
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Password"})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error ID": "token error"})
		return
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.FormatUint(user.UserID, 10),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	tokens, errs := claim.SignedString([]byte(SecretKey2))
	if errs != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error userID": "token error"})
		return
	}

	// Set token in response header
	c.Set("Authorization", "Bearer "+token)
	c.Set("Authorizations", "Bearers "+tokens)
	responseData := gin.H{
		"message":  "successfull",
		"token":    token,
		"jwt-user": tokens,
		"data":     user,
		"userid":   user.UserID, // Mengirim data pengguna lengkap
	}
	c.JSON(http.StatusOK, responseData)
	return
}
