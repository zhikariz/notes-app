package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"notes-app/account"
	"notes-app/auth"
	"notes-app/handler"
	"notes-app/helper"
	"notes-app/note"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	router.Use(location.New(location.Config{
		Scheme:  os.Getenv("SCHEME"),
		Host:    os.Getenv("BASE_URL"),
		Base:    os.Getenv("DEFAULT_ROUTE"),
		Headers: location.Headers{Scheme: "X-Forwarded-Proto", Host: "X-Forwarded-For"},
	}))
	db := helper.SetupDb()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Repository
	authService := auth.NewService()
	accountRepository := account.NewRepository(db)
	noteRepository := note.NewRepository(db)

	// Service
	accountService := account.NewService(accountRepository)
	noteService := note.NewService(noteRepository)

	// Handler
	accountHandler := handler.NewAccountHandler(authService, accountService)
	noteHandler := handler.NewNoteHandler(authService, noteService)

	api := router.Group("api/v1/")

	api.POST("/register", accountHandler.RegisterAccount)
	api.POST("/login", accountHandler.Login)
	api.GET("/verify-account/:token", accountHandler.VerifyEmail)
	api.POST("/reset-password", accountHandler.ResetPassword)
	api.GET("/verify-password/:token", accountHandler.VerifyReset)

	api.GET("/notes", noteHandler.GetAllNote)
	api.GET("/notes/:id", noteHandler.GetNoteById)
	api.GET("/notes/account", authMiddleware(authService, accountService), noteHandler.GetNoteByAccountId)
	api.POST("/notes", authMiddleware(authService, accountService), noteHandler.CreateNote)
	api.PUT("/notes/:id", authMiddleware(authService, accountService), noteHandler.UpdateNote)
	api.DELETE("/notes/:id", authMiddleware(authService, accountService), noteHandler.DeleteNote)

	api.POST("/accounts", adminMiddleware(authService, accountService), accountHandler.CreateAccount)

	router.Run()
}

func errorMiddleware(somethingAwsm interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.APIResponse("Email checking failed !", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
}

func adminMiddleware(authService auth.Service, accountService account.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			helper.AuthorizationHandling(c)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			helper.AuthorizationHandling(c)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			helper.AuthorizationHandling(c)
			return
		}

		accountId := int(claim["account_id"].(float64))

		account, err := accountService.GetAccountById(accountId)

		if err != nil {
			helper.AuthorizationHandling(c)
			return
		}

		if claim["role"] != "Admin" {
			helper.AuthorizationHandling(c)
			return
		}

		c.Set("currentUser", account)
	}
}

func authMiddleware(authService auth.Service, accountService account.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			helper.AuthorizationHandling(c)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			helper.AuthorizationHandling(c)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			helper.AuthorizationHandling(c)
			return
		}

		accountId := int(claim["account_id"].(float64))

		account, err := accountService.GetAccountById(accountId)

		if err != nil {
			helper.AuthorizationHandling(c)
			return
		}

		c.Set("currentUser", account)
	}
}
