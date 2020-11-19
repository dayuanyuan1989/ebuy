package main

import (
	"encoding/json"
	"strings"

	_ "ebuy/configs"
	_ "ebuy/models"
	_ "ebuy/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

// filterUser ...
func filterUser(ctx *context.Context) {
	sessID := ctx.Input.Session("UserId")
	isLoginURL := strings.Contains(ctx.Request.RequestURI, "/login")
	if sessID == nil && !isLoginURL {
		ctx.Redirect(302, "/login")
	}
}

// 添加日志拦截器
func filterLog(ctx *context.Context) {
	url, _ := json.Marshal(ctx.Input.URL())
	params, _ := json.Marshal(ctx.Request.Form)
	outputBytes, _ := json.Marshal(ctx.Input.Data()["json"])
	divider := " - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -"
	topDivider := "┌" + divider
	middleDivider := "├" + divider
	bottomDivider := "└" + divider
	outputStr := "\n" + topDivider + "\n│ 请求地址:" + string(url) + "\n" + middleDivider + "\n│ 请求参数:" + string(params) + "\n│ 返回数据:" + string(outputBytes) + "\n" + bottomDivider
	logs.Info(outputStr)
}

func main() {

	//注册过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, filterUser, false)
	beego.InsertFilter("/*", beego.FinishRouter, filterLog, false)

	beego.Run()
}
