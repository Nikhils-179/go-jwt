package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type Response struct{
	Message string `json:"message"`
}

type Jwks struct{
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct{
	Kty string `json:"kty"`
	Kid string `json:"kif"`
	Use string `json:"use"`
	N   string `json:"keys"`
	E   string `json:"keys"`
	X5c []string `json:"keys"`
}

type Joke struct{
	ID int `json:"id" binding: "required" `
	Likes int `json:"likes" `
	Joke string `json:"joke" binding: "required" `
}

var jokes = []Joke{
	Joke{1,0,"Did you hear about the restaurant on the moon?Great food,no atmosphere"},
	Joke{2,0,"What do you call a fake noodle?An Impasta"},
	Joke{3, 0, "How many apples grow on a tree? All of them."},
	Joke{4, 0, "Want to hear a joke about paper? Nevermind it's tearable."},
	Joke{5, 0, "I just watched a program about beavers. It was the best dam program I've ever seen."},
	Joke{6, 0, "Why did the coffee file a police report? It got mugged."},
	Joke{7, 0, "How does a penguin build it's house? Igloos it together."},
}

var jwtMiddleWare *jwtmiddleware.JWTMiddleware


func main() {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token)(interface{},error){
			aud := os.Getenv("AUTH0_API_AUDIENCE")
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud,false)
			if !checkAudience{
				return token, errors.New("Invalid audience")
			}
			iss := os.Getenv("AUTH0_DOMAIN")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss,false)
			if !checkIss{
				return token , errors.New("Invalid Issuer")
			}

			cert,err:= getPemCert(token)
			if err!=nil{
				log.Fatalf("could not get cert : %+v",err)
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result,nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	jwtMiddleWare = jwtMiddleware
	router := gin.Default()

	// Serve the index.html file as the root URL
	router.GET("/", func(c *gin.Context) {
		c.File("./views/js/index.html")
	})

	// Serve static files from the "/js" URL prefix
	router.Use(static.Serve("/js", static.LocalFile("./views/js", true)))

	// Define your API routes
	api := router.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/jokes", authMiddleware() ,JokeHandler)
		api.POST("/jokes/like/:jokeID", authMiddleware() , LikeJoke)
	}

	// Start the Gin router
	router.Run(":3000")
}

func getPemCert(token *jwt.Token)(string,error){
	cert := ""
	resp,err := http.Get(os.Getenv("AUTH0_DOMAIN") + ".well-known/jwks.json")

	if err != nil{
		return cert,err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err!=nil{
		return cert,err
	}

	x5c := jwks.keys[0].X5c

	for k , v := range x5c{
		if token.Header["kid"] == jwks.keys[k].Kid{
			cert = " ------BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == ""{
		return cert , errors.New("unable to find appropriate key.")
	}
	return cert, nil
}
func authMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		err := jwtMiddleWare.CheckJWT(c.Writer,c.Request)
		if err !=nil {
			fmt.Println(err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}

func JokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, jokes)
}


func LikeJoke(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	if jokeid,err:= strconv.Atoi(c.Param("jokeID"));err==nil{
		for i := 0 ; i < len(jokes) ; i++{
			if jokes[i].ID == jokeid{
				jokes[i].Likes++
			}
		}
		c.JSON(http.StatusOK, &jokes)
	}else{
		c.AbortWithStatus(http.StatusNotFound)
	}
}
