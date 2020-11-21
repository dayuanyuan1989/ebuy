package task

import (
	"ebuy/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
)

func init() {
	// TODO: 目前这种定时方式,效率太差,有没有类似触发机制的方式?
	// 使用堆栈的机制,卖方不断的从上往下写数据
	// 买方不断的从下往上写数据,当两边碰面,说明达成交易
	// deal产生
	tk1 := toolbox.NewTask("tk1", "*/1 * * * * *", func() error {
		TryDealByGameID(1)
		return nil
	})
	toolbox.AddTask("tk1", tk1)

	return
}

// DealData ...
type DealData struct {
	ID        int
	UserID    int
	Count     uint64
	DealCount uint64
	LeftCount uint64
	Price     uint64
}

// TryDealByGameID ...
func TryDealByGameID(gameID int) (retErr error) {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		logs.Error("db error, err: %s", err.Error())
		return err
	}
	defer func() {
		if retErr != nil {
			logs.Warn("error occur , will roll back: %s", retErr.Error())
			if err := o.Rollback(); err != nil {
				logs.Error("db rollback error: %s", err.Error())
				return
			}
		}
	}()

	// 1. select from t_sell
	var sellMaps []orm.Params
	selNum, err := o.QueryTable("t_game_seller").
		Filter("PlatformGame__id", gameID).
		RelatedSel("t_user").
		OrderBy("price").
		Limit(1000).
		Values(&sellMaps)
	if err != nil {
		logs.Error("db error, err: %s", err.Error())
		return err
	}

	logs.Debug("query num=%d", selNum)

	if selNum == 0 {
		return nil
	}

	// 2. select from t_buy
	var buyMaps []orm.Params
	buyNum, err := o.QueryTable("t_game_buyer").
		Filter("PlatformGame__id", gameID).
		RelatedSel("t_user").
		OrderBy("-price").
		Limit(1000).
		Values(&buyMaps)
	if err != nil {
		logs.Error("db error, err: %s", err.Error())
		return err
	}

	logs.Debug("query num=%d", buyNum)

	if buyNum == 0 {
		return nil
	}

	// 3. compare
	if buyMaps[0]["Price"].(uint64) < sellMaps[0]["Price"].(uint64) {
		return nil
	}

	// 4. query curPrice
	var tGamePrice models.CurPlatformGamePrice
	if err := o.QueryTable("t_cur_platform_game_price").
		Filter("platform_game_id", gameID).
		One(&tGamePrice); err != nil {
		logs.Error("db error, err: %s", err.Error())
		return err
	}
	curPrice := tGamePrice.Price

	// 4. deal
	buyDeals := []*DealData{}
	sellDeals := []*DealData{}

	var i, j int
label:
	for ; i < len(sellMaps); i++ {
		sellMap := sellMaps[i]
		for ; j < len(buyMaps); j++ {
			buyMap := buyMaps[j]
			// 当买家出价高于卖家，达成成交
			if buyMap["Price"].(uint64) >= sellMap["Price"].(uint64) &&
				sellMap["UndealCount"].(uint64) != 0 &&
				buyMap["UndealCount"].(uint64) != 0 {
				// 实际成交价格
				var dealPrice uint64
				// 买家价格高于当前市价时，卖出低于市价，成交价格以市价结算
				if buyMap["Price"].(uint64) >= curPrice &&
					sellMap["Price"].(uint64) < curPrice {
					// 当购买数量大于卖出数时，购买成交部分，卖出全部成交
					// 当购买数量小于卖出数时，购买全部成交，卖出成交部分
					dealPrice = curPrice
				} else if buyMap["Price"].(uint64) >= curPrice &&
					sellMap["Price"].(uint64) >= curPrice {
					// 买家价格高于当前市价时，卖出高于市价，成交价格以卖出价成交
					dealPrice = sellMap["Price"].(uint64)
				} else if buyMap["Price"].(uint64) < curPrice &&
					sellMap["Price"].(uint64) >= curPrice {
					// 买家价格低于当前市价时，卖出高于市价，不可能
					// ignore, 不可能，因为这样不满足 buyMap["Price"].(uint64) >= sellMap["Price"].(uint64)
				} else {
					// 买家价格低于当前市价时，卖出低于市价，成交价格以买家出价成交
					dealPrice = buyMap["Price"].(uint64)
				}
				if dealPrice != 0 {
					minCount := Min(buyMap["UndealCount"].(uint64), sellMap["UndealCount"].(uint64))
					dealCount := minCount
					sellMap["DealCount"] = dealCount
					sellMap["UndealCount"] = sellMap["UndealCount"].(uint64) - dealCount
					buyMap["DealCount"] = dealCount
					buyMap["UndealCount"] = buyMap["UndealCount"].(uint64) - dealCount
					buyDeals = append(buyDeals, &DealData{
						ID:        int(buyMap["ID"].(int64)),
						UserID:    int(buyMap["User"].(int64)),
						Count:     buyMap["Count"].(uint64),
						DealCount: dealCount,
						LeftCount: buyMap["UndealCount"].(uint64),
						Price:     dealPrice,
					})
					sellDeals = append(sellDeals, &DealData{
						ID:        int(sellMap["ID"].(int64)),
						UserID:    int(sellMap["User"].(int64)),
						Count:     sellMap["Count"].(uint64),
						DealCount: dealCount,
						LeftCount: sellMap["UndealCount"].(uint64),
						Price:     dealPrice,
					})
					curPrice = dealPrice
					// 当是卖方未处理数为0时,让卖方进行循环
					if sellMap["UndealCount"].(uint64) == 0 &&
						buyMap["UndealCount"].(uint64) != 0 {
						break
					} else if sellMap["UndealCount"].(uint64) != 0 &&
						buyMap["UndealCount"].(uint64) == 0 {
						// 当是买方未处理数为0时,让买方进行循环
						continue
					} else if sellMap["UndealCount"].(uint64) == 0 &&
						buyMap["UndealCount"].(uint64) == 0 {
						// 当都为0时,先将买方+1,后跳出当前买方循环,进行卖方下一个循环值
						j++
						continue
					} else {
						// 都不为0不可能的
					}
				}
			} else {
				break label
			}
		}
	}

	for _, buyDeal := range buyDeals {
		logs.Info("buyDeals: %+v", *buyDeal)
	}
	for _, sellDeal := range sellDeals {
		logs.Info("sellDeals: %+v", *sellDeal)
	}
	//return

	// 5. update
	// 5.1 update buyer & seller & balance
	for _, buyDeal := range buyDeals {
		gameSel := models.GameBuyer{
			ID: buyDeal.ID,
			//DealPrice:   buyDeal.Price,
			Count:       buyDeal.Count,
			User:        &models.User{ID: buyDeal.UserID},
			DealCount:   buyDeal.Count - buyDeal.LeftCount,
			UndealCount: buyDeal.LeftCount,
		}
		cols := []string{"Count", "DealCount", "UndealCount"}
		if buyDeal.LeftCount == 0 {
			cols = append(cols, "State")
			gameSel.State = models.StateDone
		}
		if _, err := o.Update(&gameSel, cols...); err != nil {
			logs.Error("db error, err: %s", err.Error())
			retErr = err
			return
		}
		if _, err := o.QueryTable("t_user").
			Filter("ID", buyDeal.UserID).
			Update(orm.Params{
				"Balance": orm.ColValue(orm.ColMinus, buyDeal.Price*buyDeal.DealCount),
			}); err != nil {
			logs.Error("db error, err: %s", err.Error())
			retErr = err
			return
		}
	}

	for _, sellDeal := range sellDeals {
		gameSel := models.GameSeller{
			ID:          sellDeal.ID,
			Count:       sellDeal.Count,
			User:        &models.User{ID: sellDeal.UserID},
			DealCount:   sellDeal.Count - sellDeal.LeftCount,
			UndealCount: sellDeal.LeftCount,
		}
		cols := []string{"Count", "DealCount", "UndealCount"}
		if sellDeal.LeftCount == 0 {
			cols = append(cols, "State")
			gameSel.State = models.StateDone
		}
		if _, err := o.Update(&gameSel, cols...); err != nil {
			logs.Error("db error, err: %s", err.Error())
			retErr = err
			return
		}
		if _, err := o.QueryTable("t_user").
			Filter("ID", sellDeal.UserID).
			Update(orm.Params{
				"Balance": orm.ColValue(orm.ColAdd, sellDeal.Price*sellDeal.DealCount),
			}); err != nil {
			logs.Error("db error, err: %s", err.Error())
			retErr = err
			return
		}
	}

	// update cur price
	if _, err := o.QueryTable("t_cur_platform_game_price").
		Filter("platform_game_id", gameID).
		Update(orm.Params{"price": curPrice}); err != nil {
		logs.Error("db error, err: %s", err.Error())
		retErr = err
		return
	}

	//if err := o.Rollback(); err != nil {
	if err := o.Commit(); err != nil {
		logs.Error("db commit eroor: %s", err.Error())
	}
	return
}

// Min ...
func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
