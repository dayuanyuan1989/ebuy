package routers

import (
	"ebuy/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// index
	beego.Router("/", &controllers.MainController{})

	// login
	beego.Router("/login/", &controllers.LoginController{}, "*:Index")
	// register
	beego.Router("/login/register", &controllers.LoginController{}, "post:Register")
	// login
	beego.Router("/login/login", &controllers.LoginController{}, "post:Login")
	// logout
	beego.Router("/login/logout", &controllers.LoginController{}, "post:Logout")

	// game
	beego.Router("/game/", &controllers.GameController{}, "*:Index")
	// lower price
	beego.Router("/game/sell/top", &controllers.GameController{}, "get:LowerSellPrice")
	// high price
	beego.Router("/game/buy/top", &controllers.GameController{}, "get:HighBuyPrice")
	// sell
	beego.Router("/game/sell", &controllers.GameController{}, "post:Sell")
	// buy
	beego.Router("/game/buy", &controllers.GameController{}, "post:Buy")
}
