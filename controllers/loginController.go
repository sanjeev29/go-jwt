package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/sanjeev29/go-jwt/helpers"
	"github.com/sanjeev29/go-jwt/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func verifyPassword(userPassword string, providedPassword string) (bool, string) {
	isValid := true
	msg := ""

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		isValid = false
		msg = "Email or password is incorrect."
	}

	return isValid, msg
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Find user from database
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Email or password is incorrect.",
			})
			return
		}

		isValid, msg := verifyPassword(*user.Password, *foundUser.Password)
		if !isValid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "User not found.",
			})
			return
		}

		token, refreshToken, err := helper.GenerateAllTokens(
			*foundUser.Email,
			*foundUser.First_name,
			*foundUser.Last_name,
			*foundUser.User_type,
			*&foundUser.User_id,
		)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		// Update all user tokens
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		// Find user from database
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}
