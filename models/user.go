package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/rs/xid"
)

// User 用户信息
type User struct {
	UserID   string `gorm:"primary_key"`
	Name     string `gorm:"size:255;not null"`
	Tel      string `gorm:"not null;unique_index"`
	Password string `gorm:"not null"`
	Position string `gorm:"not null"`
}

// Role 用户角色权限
type Role struct {
	UserID string `gorm:"primary_key"`
	Admin  bool
	CM     bool
	PM     bool
}

// RoleType 角色权限枚举
type RoleType int

const (
	// Admin 管理员
	Admin RoleType = iota
	// CM 保留管理权限
	CM
	// PM 保留管理权限
	PM
)

// VerifyUser 验证用户账号密码
func VerifyUser(tel string, password string) (*User, error) {
	var user User
	if db.Table("users").Where("tel = ?", tel).First(&user).RecordNotFound() {
		return nil, errors.New("账号不存在")
	}
	h := md5.New()
	h.Write([]byte(user.Password))

	if hex.EncodeToString(h.Sum(nil)) != password {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

// EnableUserRole 更新角色用户权限
func EnableUserRole(userID string, role RoleType) error {
	return enableUserRole(db, userID, role)
}

func enableUserRole(tx *gorm.DB, userID string, role RoleType) error {
	var err error
	switch role {
	case Admin:
		err = tx.Table("roles").Where("user_id = ?", userID).UpdateColumn("admin", true).Error
	case CM:
		err = tx.Table("roles").Where("user_id = ?", userID).UpdateColumn("cm", true).Error
	case PM:
		err = tx.Table("roles").Where("user_id = ?", userID).UpdateColumn("pm", true).Error
	}
	return err
}

// GetUserRole 获取用户角色权限
func GetUserRole(userID string) ([]string, error) {
	var role Role
	var roles []string
	err := db.Table("roles").Where("user_id = ?", userID).First(&role).Error
	if err != nil {
		return roles, err
	}

	if role.Admin {
		roles = append(roles, "admin")
	}
	if role.CM {
		roles = append(roles, "cm")
	}
	if role.PM {
		roles = append(roles, "pm")
	}

	return roles, nil
}

// DisableUserRole 更新角色用户权限
func DisableUserRole(userID string, role RoleType) error {
	var err error
	switch role {
	case Admin:
		err = db.Table("roles").Where("user_id = ?", userID).UpdateColumn("admin", false).Error
	case CM:
		err = db.Table("roles").Where("user_id = ?", userID).UpdateColumn("cm", false).Error
	case PM:
		err = db.Table("roles").Where("user_id = ?", userID).UpdateColumn("pm", false).Error
	}
	return err
}

// CreateUser 创建公司普通员工账号
func CreateUser(user *User) error {
	user.UserID = xid.New().String()
	user.Password = user.Tel //默认密码

	tx := db.Begin()

	err := tx.Table("users").Create(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table("roles").Create(&Role{
		UserID: user.UserID,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = enableUserRole(tx, user.UserID, CM)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// DeleteUser 删除公司普通员工账号
func DeleteUser(userID string) error {
	tx := db.Begin()
	err := tx.Table("roles").Where("user_id = ?", userID).Delete(Role{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Table("users").Where("user_id = ?", userID).Delete(User{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// UpdateUser 更新用户信息
func UpdateUser(userID string, newName, newTel, newPassword, newPosition string) error {
	user := User{
		UserID:   userID,
		Name:     newName,
		Tel:      newTel,
		Password: newPassword,
		Position: newPosition,
	}
	return db.Table("users").Table("users").Where("user_id = ?", userID).Updates(user).Error
}

// UpdatePassword 更新密码
func UpdatePassword(userID string, newPassword string) error {
	return db.Table("users").Table("users").Where("user_id = ?", userID).UpdateColumn("password", newPassword).Error
}

// ListUsers 所有用户
func ListUsers() (users []User, err error) {
	err = db.Table("users").Table("users").Select("user_id, name, tel, position").Find(&users).Error
	return
}

// GetUser 获取用户信息
func GetUser(userID string) *User {
	var user User
	db.Table("users").Where("user_id = ?", userID).Select("user_id, name, tel, position").First(&user)
	return &user
}

// HaveUser 查询是否有此用户
func HaveUser(userID string) bool {
	return !db.Table("users").Where("user_id = ?", userID).First(&User{}).RecordNotFound()
}

// HaveTel 查询是否有此账号
func HaveTel(tel string) bool {
	return !db.Table("users").Where("tel = ?", tel).First(&User{}).RecordNotFound()
}

func createAdminUser() error {
	var count int
	if err := db.Table("users").Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		user := User{
			UserID:   xid.New().String(),
			Name:     "admin",
			Tel:      "11223344556",
			Password: "11223344556",
			Position: "管理员",
		}

		tx := db.Begin()

		err := tx.Table("users").Create(user).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Table("roles").Create(&Role{
			UserID: user.UserID,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = enableUserRole(tx, user.UserID, Admin)
		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	}
	return nil
}
