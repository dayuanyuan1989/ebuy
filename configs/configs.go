package configs

import (
	"github.com/astaxie/beego/logs"
)

func init() {
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.Info("configs init!")
}
