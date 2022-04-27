package controller

import (
	"fmt"
	"net/http"

	"github.com/SpicyChickenFLY/kamisado/service"
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
	service.UserLogined(nickname)

	c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("Welcome, %s", nickname)})
}

// Register new user
func Register(c *gin.Context) {}
