package controllers

import (
	"ebuy/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// GameController ...
type GameController struct {
	beego.Controller
}

// Index ...
func (c *GameController) Index() {
	// query
	o := orm.NewOrm()
	tPGame := new(models.PlatformGame)
	// query DB
	var maps []orm.Params
	qs := o.QueryTable(tPGame)
	num, err := qs.Limit(3).Values(&maps)
	if err != nil {
		logs.Error("query table error, err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	logs.Info("query num: %d", num)

	for _, m := range maps {
		var curPrice models.CurPlatformGamePrice
		if err := o.QueryTable("t_cur_platform_game_price").
			Filter("PlatformGame__Id", m["ID"]).
			RelatedSel().
			One(&curPrice); err != nil {
			logs.Error("query table failed, err: %s", err.Error())
		}
		logs.Debug("cur price: %+v", curPrice)
	}

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	c.Ctx.Output.JSON(maps, true, false)
	return
}

// LowerSellPrice ...
func (c *GameController) LowerSellPrice() {
	var (
		gameID = c.GetString("GameId")
	)
	// query
	o := orm.NewOrm()
	// query DB
	type QueryPrice struct {
		Price uint64
		Count uint64
	}

	var qprices []QueryPrice
	num, err := o.Raw(`
SELECT
	price,
	sum(count) AS count
FROM
	t_game_seller
WHERE
	platform_game_id = ?
GROUP BY
	price
ORDER BY
	price ASC
LIMIT 10;`, gameID).QueryRows(&qprices)
	if err != nil {
		logs.Error("query table error, err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	logs.Info("query num: %d", num)

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	c.Ctx.Output.JSON(qprices, true, false)
	return
}

// HighBuyPrice ...
func (c *GameController) HighBuyPrice() {
	var (
		gameID, _ = c.GetInt("GameId")
	)
	// query
	o := orm.NewOrm()
	// query DB
	type QueryPrice struct {
		Price uint64
		Count uint64
	}

	var qprices []QueryPrice
	num, err := o.Raw(`
SELECT
	price,
	sum(count) AS count
FROM
	t_game_buyer
WHERE
	platform_game_id = ?
GROUP BY
	price
ORDER BY
	price DESC
LIMIT 10;`, gameID).QueryRows(&qprices)
	if err != nil {
		logs.Error("query table error, err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	logs.Info("query num: %d", num)

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	c.Ctx.Output.JSON(qprices, true, false)
	return
}

// Buy ...
func (c *GameController) Buy() {
	var (
		gameID, _ = c.GetInt("GameId")
		price, _  = c.GetInt("Price")
		count, _  = c.GetInt("Count")
	)
	// query
	// 也可以直接使用对象作为表名
	o := orm.NewOrm()
	tGBuyer := new(models.GameBuyer)
	tGBuyer.User.ID = c.GetSession("UserId").(int)
	tGBuyer.PlatformGame.ID = gameID
	tGBuyer.Price = uint64(price)
	tGBuyer.Count = uint64(count)
	tGBuyer.UndealCount = uint64(count)
	_, err := o.Insert(tGBuyer)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	return
}

// Sell ...
func (c *GameController) Sell() {
	var (
		gameID, _ = c.GetInt("GameId")
		price, _  = c.GetInt("Price")
		count, _  = c.GetInt("Count")
	)
	// query
	// 也可以直接使用对象作为表名
	o := orm.NewOrm()
	tGSeller := new(models.GameSeller)
	tGSeller.User.ID = c.GetSession("UserId").(int)
	tGSeller.PlatformGame.ID = gameID
	tGSeller.Price = uint64(price)
	tGSeller.Count = uint64(count)
	tGSeller.UndealCount = uint64(count)
	_, err := o.Insert(tGSeller)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	return
}
