package task

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

func init() {
	tk1 := toolbox.NewTask("tk1", "*/1 * * * * *", func() error {
		o := orm.NewOrm()
		o.QueryTable(tPGame)

		// 1. select from t_sell

		// 2. select from t_buy

		// 3. compare
		fmt.Println("tk1")
		return nil
	})
	toolbox.AddTask("tk1", tk1)

	toolbox.StartTask()
	defer toolbox.StopTask()
}
