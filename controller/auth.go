package controller

import (
	"fmt"
	"net/http"

	"github.com/SpicyChickenFLY/game-server/service"
	"github.com/SpicyChickenFLY/game-server/utils"
	"github.com/gin-gonic/gin"
)

// ShowIndex page for user
func ShowIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// ShowLogin page for user
func ShowLogin(c *gin.Context) {

}

// Login with user nickname(no password for now)
func Login(c *gin.Context) {
	nickname := c.Param("nickname")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "nickname should not be empty"})
	}

	// sign jwt token for user
	token, err := utils.Sign(nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "server encounter error when signing token: " + err.Error()})
		return
	}

	err = service.Login(nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "server encounter error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "msg": fmt.Sprintf("Welcome, %s", nickname)})
}

// Register new user
func Register(c *gin.Context) {}
