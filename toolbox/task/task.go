package task

import (
	"ebuy/models"
)

func init() {
	// tk1 := toolbox.NewTask("tk1", "*/1 * * * * *", func() error {
	// 	o := orm.NewOrm()
	// 	o.QueryTable(tPGame)

	// 	// 1. select from t_sell

	// 	// 2. select from t_buy

	// 	// 3. compare
	// 	fmt.Println("tk1")
	// 	return nil
	// })
	// toolbox.AddTask("tk1", tk1)

	// toolbox.StartTask()
	// defer toolbox.StopTask()
}

func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// tryDeal 将满足成交的交易单完成成交
func tryDeal() (retErr error) {

	type DealData struct {
		ID        int
		Count     uint64
		DealCount uint64
		LeftCount uint64
		Price     uint64
	}

	buyDeals := []*DealData{}
	sellDeals := []*DealData{}

	// 当前成交价，市价
	var curPrice uint64
	buyers := []*models.GameBuyer{}
	sellers := []*models.GameSeller{}
label:
	for _, seller := range sellers {
		for _, buyer := range buyers {
			// 当买家出价高于卖家，达成成交
			if buyer.Price >= seller.Price {
				// 实际成交价格
				var dealPrice uint64
				// 买家价格高于当前市价时，卖出低于市价，成交价格以市价结算
				if buyer.Price >= curPrice && seller.Price < curPrice {
					// 当购买数量大于卖出数时，购买成交部分，卖出全部成交
					// 当购买数量小于卖出数时，购买全部成交，卖出成交部分
					dealPrice = curPrice
				} else if buyer.Price >= curPrice && seller.Price >= curPrice {
					// 买家价格高于当前市价时，卖出高于市价，成交价格以卖出价成交
					dealPrice = seller.Price
				} else if buyer.Price < curPrice && seller.Price >= curPrice {
					// 买家价格低于当前市价时，卖出高于市价，不可能
					// ignore, 不可能，因为这样不满足 buyer.Price >= seller.Price
				} else {
					// 买家价格低于当前市价时，卖出低于市价，成交价格以买家出价成交
					dealPrice = buyer.Price
				}
				if dealPrice != 0 {
					minCount := Min(buyer.Count, seller.Count)
					dealCount := minCount
					seller.Count -= dealCount
					buyer.Count -= dealCount
					buyDeals = append(buyDeals, &DealData{
						ID:        buyer.ID,
						Count:     buyer.Count,
						DealCount: dealCount,
						LeftCount: buyer.Count - dealCount,
						Price:     dealPrice,
					})
					sellDeals = append(sellDeals, &DealData{
						ID:        seller.ID,
						Count:     seller.Count,
						DealCount: dealCount,
						LeftCount: seller.Count - dealCount,
						Price:     dealPrice,
					})
					curPrice = dealPrice
				}
			} else {
				break label
			}
		}
	}

	return
}
