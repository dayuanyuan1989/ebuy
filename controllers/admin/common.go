package admin

import (
	"ebuy/models"

	resp "ebuy/controllers/response"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	captcha2 "github.com/astaxie/beego/utils/captcha"
)

// 全局验证码结构体
var captcha *captcha2.Captcha

// init函数初始化captcha
func init() {
	// 验证码功能
	// 使用Beego缓存存储验证码数据
	store := cache.NewMemoryCache()
	// 创建验证码
	captcha = captcha2.NewWithFilter("/captcha", store)
	// 设置验证码长度
	captcha.ChallengeNums = 4
	// 设置验证码模板高度
	captcha.StdHeight = 50
	// 设置验证码模板宽度
	captcha.StdWidth = 120
}

// CommonController ...
type CommonController struct {
	beego.Controller
}

// Index ...
func (c *CommonController) Index() {
	c.TplName = "admin/common/login.html"
}

// Register ...
func (c *CommonController) Register() {
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
		c.Ctx.Output.JSON(resp.ResponseBase{Code: -1, Error: err.Error()}, true, false)
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
		c.Ctx.Output.JSON(resp.ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	//注册session值
	c.SetSession("UserId", tUser.ID)
	c.Data["UserId"] = c.GetSession("UserId")
	c.TplName = "login.html"
}

// Login ...
func (c *CommonController) Login() {

	// 验证码验证
	if !captcha.VerifyReq(c.Ctx.Request) {
		//c.Data["Error"] = "验证码错误"
		logs.Error("验证码错误")
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(resp.ResponseBase{Code: -1, Error: "验证码错误"}, true, false)
		return
	}

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
		c.Ctx.Output.JSON(resp.ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	//注册session值
	c.SetSession("UserId", tUser.ID)
	// 重定向到首页
	c.Redirect("/", 302)
}

// Logout ...
func (c *CommonController) Logout() {
	c.DelSession("UserId")
	c.Redirect("/login", 302)
}
