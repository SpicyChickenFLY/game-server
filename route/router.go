package route

import (
	"net/http"

	"github.com/SpicyChickenFLY/kamisado/controller"
	"github.com/SpicyChickenFLY/kamisado/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.StaticFS("/static", http.Dir("static"))
	root := r.Group("/")
	{
		root.GET("", controller.ShowIndex)
		userGroup := root.Group("user/")
		{
			userGroup.GET("", controller.ShowLogin)
			userGroup.GET(":nickname", controller.Login)
			userGroup.POST("", controller.Register)
		}
		roomGroup := root.Group("room/", middleware.AuthJWT())
		{
			roomGroup.GET("", controller.ListRooms)
			roomGroup.GET(":room-id", controller.JoinRoom)
			roomGroup.POST("", controller.CreateRoom)
		}
	}
	return r
}
