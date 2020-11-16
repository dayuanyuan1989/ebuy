package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// Game ...
type Game struct {
	ID            int    `orm:"column(id);pk;auto"`
	Name          string `orm:"unique"`
	State         int
	PubDate       int
	CreateTime    time.Time       `orm:"auto_now_add;type(datetime)"`
	UpdateTime    time.Time       `orm:"auto_now;type(datetime)"`
	PlatformGames []*PlatformGame `orm:"reverse(many)"`
}

// BigPlatform ...
type BigPlatform struct {
	ID             int    `orm:"column(id);pk;auto"`
	Name           string `orm:"unique"`
	State          int
	CreateTime     time.Time        `orm:"auto_now_add;type(datetime)"`
	UpdateTime     time.Time        `orm:"auto_now;type(datetime)"`
	SmallPlatforms []*SmallPlatform `orm:"reverse(many)"`
}

// SmallPlatform ...
type SmallPlatform struct {
	ID            int          `orm:"column(id);pk;auto"`
	BigPlatform   *BigPlatform `orm:"rel(fk)"`
	Name          string       `orm:"unique"`
	State         int
	CreateTime    time.Time       `orm:"auto_now_add;type(datetime)"`
	UpdateTime    time.Time       `orm:"auto_now;type(datetime)"`
	PlatformGames []*PlatformGame `orm:"reverse(many)"`
}

// PlatformGame ...
type PlatformGame struct {
	ID            int            `orm:"column(id);pk;auto"`
	SmallPlatform *SmallPlatform `orm:"rel(fk)"`
	Game          *Game          `orm:"rel(fk)"`
	State         int
	PubDate       int
	CreateTime    time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime    time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	fmt.Println("test model init")
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:ucloudcn12@tcp(117.50.20.250:3306)/ebuy?charset=utf8")
	// 需要在init中注册定义的model
	orm.RegisterModelWithPrefix("t_", new(Game), new(BigPlatform), new(SmallPlatform), new(PlatformGame))

	orm.SetMaxIdleConns("default", 30)
	orm.SetMaxOpenConns("default", 30)

	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
	// 同步
	orm.RunSyncdb("default", false, true)
}
