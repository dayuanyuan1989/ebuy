package main

import (
	_ "ebuy/models"
	_ "ebuy/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
