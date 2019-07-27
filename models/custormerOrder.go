package models

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
)

// CustormerOrder 客户订单表
type CustormerOrder struct {
	gorm.Model
	Operator        string `gorm:"size:255"`
	Name            string
	Tel             string
	DeliveryAddress string
	DeliveryTime    string
	Amount          float32
	Deposit         float32
	Remarks         string
	State           string `gorm:"default:'未完成'"`
	Freight         float32
}

// CustormerGoods 客户订单详细商品表
type CustormerGoods struct {
	CustormerOrderID uint
	GoodsID          string
	Number           uint
}

// CreateCustormerOrder 创建客户订单
func CreateCustormerOrder(order *CustormerOrder, goods []CustormerGoods) (*CustormerOrderInfo, error) {
	tx := db.Begin()

	err := tx.Table("custormer_orders").Create(order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, v := range goods {
		v.CustormerOrderID = order.ID
		err = tx.Table("custormer_goods").Create(&v).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		number, err := getPresaleNumber(tx, v.GoodsID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = updatePresaleNumber(tx, v.GoodsID, number+v.Number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	orderInfo, err := getCustormerOrder(tx, strconv.Itoa(int(order.ID)))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return orderInfo, nil
}

// DeleteCustormerOrder 删除客户订单
func DeleteCustormerOrder(orderID string) error {
	ok, err := custormerOrderConfirmed(orderID)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("已确认订单无法删除")
	}

	tx := db.Begin()

	var goods []CustormerGoods
	err = tx.Table("custormer_goods").Where("custormer_order_id = ?", orderID).Find(&goods).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range goods {
		presaleNumber, err := getPresaleNumber(tx, v.GoodsID)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = updatePresaleNumber(tx, v.GoodsID, presaleNumber-v.Number)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Table("custormer_goods").Where(
		"custormer_order_id = ?", orderID).Delete(CustormerGoods{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table("custormer_orders").Unscoped().Where(
		"id = ?", orderID).Delete(CustormerOrder{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// ConfirmCustormerOrder 确认订单
func ConfirmCustormerOrder(orderID string, freight float32) (*CustormerOrderInfo, error) {
	ok, err := custormerOrderConfirmed(orderID)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, errors.New("订单请勿重复确认")
	}

	tx := db.Begin()

	err = tx.Table("custormer_orders").Where("id = ?", orderID).UpdateColumns(CustormerOrder{
		State:   "已完成",
		Freight: freight}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var goods []CustormerGoods
	err = tx.Table("custormer_goods").Where("custormer_order_id = ?", orderID).Find(&goods).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, v := range goods {
		presaleNumber, err := getPresaleNumber(tx, v.GoodsID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = updatePresaleNumber(tx, v.GoodsID, presaleNumber-v.Number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		number, err := getNumber(tx, v.GoodsID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if number < v.Number {
			tx.Rollback()
			return nil, errors.New("库存不足")
		}
		err = updateNumber(tx, v.GoodsID, number-v.Number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		salesVolume, err := getSalesVolume(tx, v.GoodsID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = updateSalesVolume(tx, v.GoodsID, salesVolume+v.Number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	orderInfo, err := getCustormerOrder(tx, orderID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return orderInfo, nil
}

// CustormerOrderInfo 与客户端通信的订单信息
type CustormerOrderInfo struct {
	CustormerOrder
	Goods []Commodity
}

// ListCustormerOrders 获取所有客户订单
func ListCustormerOrders() (orders []CustormerOrderInfo, err error) {
	err = db.Table("custormer_orders").Select(
		"id, created_at, operator, name, tel, delivery_address, delivery_time, amount, deposit, remarks, state, freight").Find(&orders).Error
	if err != nil {
		return
	}

	for index, v := range orders {
		err = db.Raw(`select t2.id, t2.name, t2.colour, t2.size, t2.brand, t1.number  from custormer_goods as t1 left outer join 
				commodities as t2 on t1.goods_id = t2.id where t1.custormer_order_id = ?`, v.ID).Scan(&orders[index].Goods).Error
		if err != nil {
			return
		}
	}
	return
}

func getCustormerOrder(tx *gorm.DB, orderID string) (*CustormerOrderInfo, error) {
	var order CustormerOrderInfo
	err := tx.Table("custormer_orders").Select(
		"id, created_at, operator, name, tel, delivery_address, delivery_time, amount, deposit, remarks, state, freight").Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	err = tx.Raw(`select t2.id, t2.name, t2.colour, t2.size, t2.brand, t1.number  from custormer_goods as t1 left outer join 
			commodities as t2 on t1.goods_id = t2.id where t1.custormer_order_id = ?`, orderID).Scan(&order.Goods).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func custormerOrderConfirmed(orderID string) (bool, error) {
	info := struct {
		State string `json:"state"`
	}{}
	err := db.Table("custormer_orders").Select("state").Where("id = ?", orderID).First(&info).Error
	if err != nil {
		return false, err
	}

	if info.State == "已完成" {
		return true, nil
	}

	return false, nil
}
