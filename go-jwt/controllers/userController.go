package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"strings"

	helper "github.com/Nikhils-179/go-jwt/helpers"
	"golang.org/x/crypto/bcrypt"

	database "github.com/Nikhils-179/go-jwt/database"
	"github.com/Nikhils-179/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	hasedpassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hasedpassword)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	check := true
	msg := ""
	if err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword)); err != nil {
		msg = "Incorrect password"
		check = false
	}
	return check, msg
}
//COMPLETED
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if validateErr := validate.Struct(user); validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking the email"})
			return
		}
		password := HashPassword(*user.Password)

		user.Password = &password
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking the phone"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exist"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User Item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		defer cancel()

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		}

		passwordIsValid, _ := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid {
			c.JSON(http.StatusOK, gin.H{"success": "User logged in"})
		}
		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, &foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helper.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		//pipeline stage function
		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

		groupStage := bson.D{{Key: "$group",Value: bson.D{
			{Key: "_id",Value: bson.D{{Key: "_id",Value: "null"}}},
			{Key: "total_count",Value: bson.D{{Key: "$sum",Value: 1}}},
			{Key: "data",Value: bson.D{{Key: "$push",Value: "$$ROOT"}}},
		}}}

		projectStage:= bson.D{{Key: "$project", Value: bson.D{
				{Key:"_id",Value: 0},
				{Key:"total_count",Value:1},
				{Key:"user_items",Value: bson.D{{Key: "slice",Value: []interface{}{"$data",startIndex,recordPerPage}}}},
			}},
		}

		result , err := userCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage , groupStage , projectStage,
		})
		defer cancel()
		if err!= nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error": "error occured while listing user items"})
		}

		var allUsers []bson.M
		if err = result.All(ctx , &allUsers);err!= nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK , allUsers)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "err.Error()"})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Assuming the token is sent in the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			return
		}

		// Invalidate the token
		err := helper.InvalidateToken(ctx, tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "User logged out"})
	}
}
