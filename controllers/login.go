package controllers

import (
	"ebuy/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// LoginController ...
type LoginController struct {
	beego.Controller
}

// Index ...
func (c *LoginController) Index() {
	//获取session值
	c.Data["UserId"] = c.GetSession("UserId")
	c.TplName = "login.html"
}

// Register ...
func (c *LoginController) Register() {
	var (
		userName = c.GetString("UserName")
		password = c.GetString("Password")
		//Phone    = c.GetString("Phone")
	)

	// 也可以直接使用对象作为表名
	o := orm.NewOrm()
	tUser := new(models.User)
	tUser.UserName = userName
	tUser.Password = password
	_, err := o.Insert(tUser)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	// query
	qs := o.QueryTable(tUser)
	// where
	qs = qs.Filter("user_name", userName)
	qs = qs.Filter("password", password)
	if err := qs.One(tUser); err != nil {
		errMsg := "db error"
		if err == orm.ErrMultiRows {
			// 多条的时候报错
			errMsg = "Returned Multi Rows Not One"
		} else if err == orm.ErrNoRows {
			// 没有找到记录
			errMsg = "Not row found"
		}

		logs.Error(errMsg+", err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	//注册session值
	c.SetSession("UserId", tUser.ID)
	c.Data["UserId"] = c.GetSession("UserId")
	c.TplName = "login.html"
}

// Login ...
func (c *LoginController) Login() {
	// input param
	userName := c.GetString("UserName")
	password := c.GetString("Password")

	o := orm.NewOrm()
	tUser := new(models.User)
	qs := o.QueryTable(tUser)

	// where
	qs = qs.Filter("user_name", userName)
	qs = qs.Filter("password", password)
	if err := qs.One(tUser); err != nil {
		errMsg := "db error"
		if err == orm.ErrMultiRows {
			// 多条的时候报错
			errMsg = "Returned Multi Rows Not One"
		} else if err == orm.ErrNoRows {
			// 没有找到记录
			errMsg = "Not row found"
		}

		logs.Error(errMsg+", err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	//注册session值
	c.SetSession("UserId", tUser.ID)
	// 重定向到首页
	c.Redirect("/index", 302)

}

// Logout ...
func (c *LoginController) Logout() {
	c.DelSession("UserId")
	c.Redirect("/login/login", 302)
}
