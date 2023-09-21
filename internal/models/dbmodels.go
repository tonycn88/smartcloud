package models

import (
	"log"
	"smartcloud/internal/database"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Id       uint   `gorm:"primary_key;auto_increment"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func init() {
	// 	var d Users
	// init_user_table()
}

func (p *Users) TableName() string {
	return "user"
}

func (p *Users) Init_user_table() {
	// 迁移表结构
	database.Db.AutoMigrate(&Users{})

	// 增加数据
	tx := database.Db.Create(&Users{Username: "user", Password: "123456"})
	if tx == nil {
		log.Printf("insert db error %s", tx.Error)
	}

}

func (p *Users) Exist(username string) bool {
	var d Users
	d.Username = username
	tx := database.Db.First(&d)
	if tx.Error != nil {
		log.Printf("find user error %v,user not exists", tx)
		return false
	}
	log.Println(d.Username)
	if d.Username == "" {
		return false
	}
	return true
}

func (p *Users) Find_password_by_username(username string) string {
	var d Users
	// database.Db.First(&d, 1)                    // find product with integer primary key
	database.Db.First(&d, "username = ?", username) // find product with code D42
	return d.Password
}

// func (p *User)

type Directories struct {
	gorm.Model
	Id       uint   `gorm:"primary_key;auto_increment"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (p *Directories) TableName() string {
	return "user"
}
