package http

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var OAuthConfig *oauth2.Config

func InitOAuth() {
	OAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func LoginHandler(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	c.SetCookie("oauth_state", state, 3600, "/", "localhost", false, true)
	url := OAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func CallbackHandler(c *gin.Context) {
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil || stateCookie != c.Query("state") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	token, err := OAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No ID token received"})
		return
	}

	c.SetCookie("id_token", idToken, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		idToken, err := c.Cookie("id_token")
		if err != nil || idToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token provided"})
			c.Abort()
			return
		}

		if !isTokenValid(idToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isTokenValid(tokenString string) bool {
	payload, err := idtoken.Validate(context.Background(), tokenString, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		fmt.Printf("Token validation error: %v\n", err)
		return false
	}

	if payload.Expires < time.Now().Unix() {
		fmt.Println("Token expired")
		return false
	}

	fmt.Printf("Validated token for user: %s\n", payload.Claims["email"])
	return true
}

func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
