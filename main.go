package main

import "github.com/SpicyChickenFLY/kamisado/route"

func main() {
	route := route.InitRoute()
	route.Run(":9090")
}
