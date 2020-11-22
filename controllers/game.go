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

// GameIndexResp ...
type GameIndexResp struct {
	ResponseBase
	PlatformGames []*models.PlatformGame `json:"Data,omitempty"`
	Total         int
}

// Index ...
func (c *GameController) Index() {

	var (
		limit, _ = c.GetInt("Limit")
	)
	if limit == 0 {
		limit = 5
	}

	// query
	o := orm.NewOrm()
	var maps []orm.Params
	// query top platform game
	num, err := o.QueryTable("t_top_platform_game").
		Limit(limit).Values(&maps)
	if err != nil {
		logs.Error("db query error: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	logs.Debug("query num: %d", num)
	if num == 0 {
		c.Ctx.Output.SetStatus(200)
		c.Ctx.Output.JSON(GameIndexResp{}, true, false)
		return
	}

	platGameIds := []int{}
	for _, m := range maps {
		platGameIds = append(platGameIds, int(m["PlatformGameID"].(int64)))
	}

	tPGame := new(models.PlatformGame)
	// query DB
	qs := o.QueryTable(tPGame)
	num, err = qs.Filter("ID__in", platGameIds).Values(&maps)
	if err != nil {
		logs.Error("query table error, err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	logs.Info("query num: %d", num)

	resp := GameIndexResp{}
	for _, m := range maps {
		logs.Info(m)
		pfGame := &models.PlatformGame{
			ID:      int(m["ID"].(int64)),
			State:   int(m["State"].(int64)),
			PubDate: int(m["PubDate"].(int64)),
		}
		tsmallPf := models.SmallPlatform{ID: int(m["SmallPlatform"].(int64))}
		if err := o.Read(&tsmallPf); err != nil {
			logs.Error("read db error: %s", err.Error())
		} else {
			pfGame.SmallPlatform = &tsmallPf
		}
		curPrice := models.CurPlatformGamePrice{PlatformGame: &models.PlatformGame{ID: int(m["ID"].(int64))}}
		if err := o.Read(&curPrice, "PlatformGame"); err != nil {
			logs.Error("read db error: %s", err.Error())
		} else {
			pfGame.CurPlatformGamePrice = &curPrice
		}
		resp.PlatformGames = append(resp.PlatformGames, pfGame)
	}
	resp.Total = len(resp.PlatformGames)

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	c.Ctx.Output.JSON(resp, true, false)
	return
}

// List ...
func (c *GameController) List() {

	var (
		limit, _  = c.GetInt("Limit")
		offset, _ = c.GetInt("Offset")
	)
	if limit == 0 {
		limit = 20
	}

	// query
	o := orm.NewOrm()
	var maps []orm.Params
	tPGame := new(models.PlatformGame)
	// query DB
	qs := o.QueryTable(tPGame)
	num, err := qs.Limit(limit).Offset(offset).Values(&maps)
	if err != nil {
		logs.Error("query table error, err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}
	total, err := qs.Count()
	if err != nil {
		logs.Error("query table error, err: %s", err.Error())
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(ResponseBase{Code: -1, Error: err.Error()}, true, false)
		return
	}

	logs.Info("query num: %d", num)

	resp := GameIndexResp{}
	for _, m := range maps {
		logs.Info(m)
		pfGame := &models.PlatformGame{
			ID:      int(m["ID"].(int64)),
			State:   int(m["State"].(int64)),
			PubDate: int(m["PubDate"].(int64)),
		}
		tsmallPf := models.SmallPlatform{ID: int(m["SmallPlatform"].(int64))}
		if err := o.Read(&tsmallPf); err != nil {
			logs.Error("read db error: %s", err.Error())
		} else {
			pfGame.SmallPlatform = &tsmallPf
		}
		curPrice := models.CurPlatformGamePrice{PlatformGame: &models.PlatformGame{ID: int(m["ID"].(int64))}}
		if err := o.Read(&curPrice, "PlatformGame"); err != nil {
			logs.Error("read db error: %s", err.Error())
		} else {
			pfGame.CurPlatformGamePrice = &curPrice
		}
		resp.PlatformGames = append(resp.PlatformGames, pfGame)
	}
	resp.Total = int(total)

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	c.Ctx.Output.JSON(resp, true, false)
	return
}

type GetGameResp struct {
	ResponseBase
	Data *models.PlatformGame `json:"Data,omitempty"`
}

// GetGame ...
func (c *GameController) GetGame() {

	var (
		gameID, _ = c.GetInt("GameId")
	)

	// query
	o := orm.NewOrm()

	tPGame := new(models.PlatformGame)
	// query DB
	qs := o.QueryTable(tPGame)
	if err := qs.Filter("id", gameID).One(tPGame); err != nil {
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
	// platform game
	tsmallPf := models.SmallPlatform{ID: tPGame.SmallPlatform.ID}
	if err := o.Read(&tsmallPf); err != nil {
		logs.Error("read db error: %s", err.Error())
	} else {
		tPGame.SmallPlatform = &tsmallPf
	}
	// cur price
	curPrice := models.CurPlatformGamePrice{PlatformGame: &models.PlatformGame{ID: tPGame.SmallPlatform.ID}}
	if err := o.Read(&curPrice, "PlatformGame"); err != nil {
		logs.Error("read db error: %s", err.Error())
	} else {
		tPGame.CurPlatformGamePrice = &curPrice
	}

	resp := GetGameResp{
		Data: tPGame,
	}

	// 暂时先返回json,后续返回view
	c.Ctx.Output.SetStatus(200)
	c.Ctx.Output.JSON(resp, true, false)
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
