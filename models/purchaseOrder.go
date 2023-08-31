package models

import (
	"errors"
	"strconv"

	"gorm.io/gorm"
)

// PurchaseOrder 采购订单表
type PurchaseOrder struct {
	gorm.Model
	Operator string `gorm:"size:255"`
	Remarks  string
	Amount   float32
	Freight  float32
	State    string `gorm:"default:'未完成'"`
}

// PurchaseGoods 采购订单详细产品表
type PurchaseGoods struct {
	PurchaseOrderID uint
	GoodsID         string
	Number          uint
}

// CreatePurchaseOrder 创建采购订单
func CreatePurchaseOrder(order *PurchaseOrder, goods []PurchaseGoods) (*PurchaseOrderInfo, error) {
	tx := db.Begin()

	err := tx.Table("purchase_orders").Create(order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, v := range goods {
		v.PurchaseOrderID = order.ID
		err = tx.Table("purchase_goods").Create(&v).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	orderInfo, err := getPurchaseOrder(tx, strconv.Itoa(int(order.ID)))
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return orderInfo, nil
}

// DeletePurchaseOrder 删除采购订单
func DeletePurchaseOrder(orderID string) error {
	ok, err := purchaseOrderConfirmed(orderID)
	if err != nil {
		return err
	}
	if ok {
		return errors.New("无法删除已经确认订单")
	}

	tx := db.Begin()

	err = tx.Table("purchase_goods").Where(
		"purchase_order_id = ?", orderID).Delete(PurchaseGoods{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table("purchase_orders").Unscoped().Where(
		"id = ?", orderID).Delete(PurchaseOrder{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// ConfirmPurchaseOrder 确认订单
func ConfirmPurchaseOrder(orderID string, freight float32) (*PurchaseOrderInfo, error) {
	ok, err := purchaseOrderConfirmed(orderID)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, errors.New("订单请勿重复确认")
	}

	tx := db.Begin()

	err = tx.Table("purchase_orders").Where("id = ?", orderID).UpdateColumns(PurchaseOrder{
		State:   "已完成",
		Freight: freight}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var goods []PurchaseGoods
	err = tx.Table("purchase_goods").Where("purchase_order_id = ?", orderID).Find(&goods).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, v := range goods {
		number, err := getNumber(tx, v.GoodsID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = updateNumber(tx, v.GoodsID, number+v.Number)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	orderInfo, err := getPurchaseOrder(tx, orderID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return orderInfo, nil
}

// PurchaseOrderInfo 订单信息
type PurchaseOrderInfo struct {
	PurchaseOrder
	Goods []Commodity `gorm:"-"`
}

// ListPurchaseOrders 所有订单
func ListPurchaseOrders() (orders []PurchaseOrderInfo, err error) {
	err = db.Table("purchase_orders").Select(
		"id, created_at, operator, amount, remarks, state, freight").Find(&orders).Error
	if err != nil {
		return
	}

	for index, v := range orders {
		err = db.Raw(`select t2.id, t2.name, t2.colour, t2.size, t2.brand, t1.number from purchase_goods as t1 left outer join 
				commodities as t2 on t1.goods_id = t2.id where t1.purchase_order_id = ?`, v.ID).Scan(&orders[index].Goods).Error
		if err != nil {
			return
		}
	}
	return
}

// GetPurchaseOrder 获取订单信息
func getPurchaseOrder(tx *gorm.DB, orderID string) (*PurchaseOrderInfo, error) {
	var order PurchaseOrderInfo
	err := tx.Table("purchase_orders").Select(
		"id, created_at, operator, amount, remarks, state, freight").Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}

	err = tx.Raw(`select t2.id, t2.name, t2.colour, t2.size, t2.brand, t1.number from purchase_goods as t1 left outer join
			commodities as t2 on t1.goods_id = t2.id where t1.purchase_order_id = ?`, orderID).Scan(&order.Goods).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func purchaseOrderConfirmed(orderID string) (bool, error) {
	info := struct {
		State string `json:"state"`
	}{}
	err := db.Table("purchase_orders").Select("state").Where("id = ?", orderID).First(&info).Error
	if err != nil {
		return false, err
	}

	if info.State == "已完成" {
		return true, nil
	}

	return false, nil
}
