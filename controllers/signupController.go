package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sanjeev29/go-jwt/database"
	helper "github.com/sanjeev29/go-jwt/helpers"
	"github.com/sanjeev29/go-jwt/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(hashedPassword)
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validate user struct
		validationErr := validate.Struct(&user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
			})
			return
		}

		// Check for duplicate email and phone number
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occured while checking for the email.",
			})
			return
		}

		// Hash user entered password
		hashedPassword := hashPassword(*user.Password)
		user.Password = &hashedPassword

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occured while checking for the email.",
			})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "This email / phone number provided already exists.",
			})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// Generate tokens
		token, refreshToken, _ := helper.GenerateAllTokens(
			*user.Email,
			*user.First_name,
			*user.Last_name,
			*user.User_type,
			*&user.User_id,
		)
		user.Token = &token
		user.Refresh_token = &refreshToken

		// Insert user to database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error during user signup.",
			})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}
