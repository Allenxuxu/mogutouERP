package models

// QueryYearFinance 查询年利润，营业额
func QueryYearFinance(year int) (profit []float32, turnover []float32, orderQuantity []int, err error) {
	profit = make([]float32, 12)
	turnover = make([]float32, 12)
	orderQuantity = make([]int, 12)
	for i := 0; i < 12; i++ {
		orderInfo := []struct {
			ID      uint
			Amount  float32
			Freight float32
		}{}
		if err = db.Table("custormer_orders").Select("id, amount, freight").Where("Year(created_at) = ? AND Month(created_at) = ?",
			year, i+1).Scan(&orderInfo).Error; err != nil {
			return
		}

		var totalProfit float32
		for _, v := range orderInfo {
			var cost float32
			if cost, err = queryOrderCost(v.ID); err != nil {
				return
			}
			totalProfit += (v.Amount - v.Freight - cost)
			turnover[i] += v.Amount
		}

		var purchaseFreight float32
		if purchaseFreight, err = queryMPurchaseFreight(year, i+1); err != nil {
			return
		}

		profit[i] = (totalProfit - purchaseFreight)
		orderQuantity[i] = len(orderInfo)
	}
	return
}

// 查询月采购单 运费
func queryMPurchaseFreight(year, month int) (float32, error) {
	var freights []float32
	err := db.Table("purchase_orders").Where("Year(created_at) = ? AND Month(created_at) = ?",
		year, month).Pluck("freight", &freights).Error
	if err != nil {
		return 0, err
	}

	var ret float32
	for _, v := range freights {
		ret += v
	}
	return ret, nil
}

// 查询订单 商品进货总费用
func queryOrderCost(orderID uint) (float32, error) {
	goodsInfo := []struct {
		PurchasePrice float32
		Number        uint
	}{}

	err := db.Raw(`select t2.purchase_price, t1.number  from custormer_goods as t1 left outer join
				commodities as t2 on t1.goods_id = t2.id where t1.custormer_order_id = ?`,
		orderID).Scan(&goodsInfo).Error
	if err != nil {
		return 0, err
	}

	var ret float32
	for _, v := range goodsInfo {
		ret += (float32(v.Number) * v.PurchasePrice)
	}

	return ret, nil
}
