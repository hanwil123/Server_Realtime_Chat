package Controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-chat/Databases"
	"go-chat/Models"
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
	rand.NewSource(time.Now().UnixNano())
	randomUserID := uint(rand.Intn(10000) + 1)
	// Create new user
	user := &Models.UserClients{
		UserID:   uint64(randomUserID),
		Username: data["username"],
		Email:    data["email"],
		Password: password,
	}
	result := Databases.CDB.Create(&user)
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

	var user Models.UserClients
	Databases.RDB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Message": "User Not Found"})
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Password"})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
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
