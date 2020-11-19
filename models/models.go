package models

import (
	"time"

	"github.com/astaxie/beego/logs"
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
	CreateTime    time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdateTime    time.Time     `orm:"auto_now;type(datetime)"`
	GameBuyers    []*GameBuyer  `orm:"reverse(many)"`
	GameSellers   []*GameSeller `orm:"reverse(many)"`
	// price ref
	CurPlatformGamePrice  *CurPlatformGamePrice   `orm:"reverse(one)"`
	HisPlatformGamePrices []*HisPlatformGamePrice `orm:"reverse(many)"`
}

// User ...
type User struct {
	ID          int    `orm:"column(id);pk;auto"`
	UserName    string `orm:"unique"`
	Password    string
	Phone       string `orm:"unique"`
	Balance     uint64
	State       int
	CreateTime  time.Time     `orm:"auto_now_add;type(datetime)"`
	UpdateTime  time.Time     `orm:"auto_now;type(datetime)"`
	GameBuyers  []*GameBuyer  `orm:"reverse(many)"`
	GameSellers []*GameSeller `orm:"reverse(many)"`
}

// GameBuyer ...
type GameBuyer struct {
	ID           int           `orm:"column(id);pk;auto"`
	User         *User         `orm:"rel(fk)"`
	PlatformGame *PlatformGame `orm:"rel(fk)"`
	Price        uint64
	Count        uint64
	DealCount    uint64
	UndealCount  uint64
	State        int
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime   time.Time `orm:"auto_now;type(datetime)"`
}

// GameSeller ...
type GameSeller struct {
	ID           int           `orm:"column(id);pk;auto"`
	User         *User         `orm:"rel(fk)"`
	PlatformGame *PlatformGame `orm:"rel(fk)"`
	Price        uint64
	Count        uint64
	DealCount    uint64
	UndealCount  uint64
	State        int
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime   time.Time `orm:"auto_now;type(datetime)"`
}

// CurPlatformGamePrice  当前游戏成交价
type CurPlatformGamePrice struct {
	ID           int           `orm:"column(id);pk"`
	PlatformGame *PlatformGame `orm:"rel(one)"`
	Price        uint64
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime   time.Time `orm:"auto_now;type(datetime)"`
}

// HisPlatformGamePrice  历史游戏成交价
type HisPlatformGamePrice struct {
	ID           int           `orm:"column(id);pk;auto"`
	PlatformGame *PlatformGame `orm:"rel(fk)"`
	Price        uint64
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	logs.Info("models init!")

	//orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:ucloudcn12@tcp(117.50.20.250:3306)/ebuy?charset=utf8")
	// 需要在init中注册定义的model
	orm.RegisterModelWithPrefix("t_",
		new(Game), new(BigPlatform), new(SmallPlatform), new(PlatformGame),
		new(User), new(GameBuyer), new(GameSeller),
		new(CurPlatformGamePrice), new(HisPlatformGamePrice))

	orm.SetMaxIdleConns("default", 30)
	orm.SetMaxOpenConns("default", 30)

	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
	// 同步
	orm.RunSyncdb("default", false, true)
	// debug
	orm.Debug = true
}
