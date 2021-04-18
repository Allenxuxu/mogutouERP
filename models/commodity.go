package models

import (
	"gorm.io/gorm"
)

// Commodity 商品表
type Commodity struct {
	ID            string `gorm:"primary_key"`
	Name          string
	Colour        string
	Size          string
	Brand         string
	Number        uint
	PresaleNumber uint
	SalesVolume   uint
	Price         float32
	PurchasePrice float32
}

// CreateCommodity 创建新商品
func CreateCommodity(commodity *Commodity) error {
	return db.Table("commodities").Create(commodity).Error
}

// DeleteCommodity 删除商品
func DeleteCommodity(commodityID string) error {
	return db.Table("commodities").Where("id = ?", commodityID).Delete(Commodity{}).Error
}

// UpdateCommodity 更新商品信息
func UpdateCommodity(commodityID string, commodity *Commodity) error {
	return db.Table("commodities").Where("id = ?", commodityID).Updates(*commodity).Error
}

// GetCommoditiesWithPurchasePrice 获取所有商品信息（包含进价）
func GetCommoditiesWithPurchasePrice() (commodities []Commodity, err error) {
	err = db.Table("commodities").Find(&commodities).Error
	return
}

// GetCommodities 获取所有商品信息
func GetCommodities() (commodities []Commodity, err error) {
	err = db.Table("commodities").Select(
		"id, name, colour, size, brand, number, presale_number, sales_volume, price").Find(&commodities).Error
	return
}

// HaveCommodity 查询是否有此商品
func HaveCommodity(id string) bool {
	return !(db.Table("commodities").Where("id = ?", id).First(&Commodity{}).Error == gorm.ErrRecordNotFound)
}

// GetCommodity 获取商品信息
func GetCommodity(id string) *Commodity {
	var commodity Commodity
	db.Table("commodities").Where("id = ?", id).Select(
		"id, name, colour, size, brand, number, presale_number, sales_volume, price").First(&commodity)
	return &commodity
}

func getPresaleNumber(tx *gorm.DB, goodsID string) (number uint, err error) {
	info := struct {
		PresaleNumber uint
	}{}
	err = tx.Table("commodities").Select("presale_number").Where("id = ?", goodsID).Scan(&info).Error
	number = info.PresaleNumber
	return
}

func updatePresaleNumber(tx *gorm.DB, goodsID string, number uint) error {
	return tx.Table("commodities").Where("id = ?", goodsID).Update("presale_number", number).Error
}

func getNumber(tx *gorm.DB, goodsID string) (number uint, err error) {
	info := struct {
		Number uint
	}{}
	err = tx.Table("commodities").Select("number").Where("id = ?", goodsID).Scan(&info).Error
	number = info.Number
	return
}

func updateNumber(tx *gorm.DB, goodsID string, number uint) error {
	return tx.Table("commodities").Where("id = ?", goodsID).Update("number", number).Error
}

func getSalesVolume(tx *gorm.DB, goodsID string) (number uint, err error) {
	info := struct {
		SalesVolume uint
	}{}
	err = tx.Table("commodities").Select("sales_volume").Where("id = ?", goodsID).Scan(&info).Error
	number = info.SalesVolume
	return
}

func updateSalesVolume(tx *gorm.DB, goodsID string, number uint) error {
	return tx.Table("commodities").Where("id = ?", goodsID).Update("sales_volume", number).Error
}
