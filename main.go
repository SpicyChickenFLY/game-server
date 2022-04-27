package main

import "github.com/SpicyChickenFLY/game-server/route"

func main() {
	route := route.InitRoute()
	route.Run(":9090")
}
