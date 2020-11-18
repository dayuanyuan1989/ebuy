package routers

import (
	"ebuy/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// index
	beego.Router("/", &controllers.MainController{})
	// login
	beego.Router("/login/index", &controllers.LoginController{}, "*:Index")
	// register
	beego.Router("/login/register", &controllers.LoginController{}, "post:Register")
	// login
	beego.Router("/login/login", &controllers.LoginController{}, "post:Login")
	// logout
	beego.Router("/login/logout", &controllers.LoginController{}, "post:Logout")
}
