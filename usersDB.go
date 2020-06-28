package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Users struct {
	gorm.Model
	Name     string //頭文字を大文字にしないと、DBにマイグレーションできない
	Password string
	Score    int
}

// DB接続
func dbInit_users() {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Init失敗")
	}
	defer db.Close()
	db.AutoMigrate(&Users{})
}

//サインアップ
func dbSignup(name string, password string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Signup失敗")
	}
	defer db.Close()
	var users Users
	if err := db.Where("name = ?", name).First(&users).Error; err == nil {
		panic("すでに同じ名前が使われています")
	} else {
		db.Create(&Users{Name: name, Password: password})
	}
}

//ログイン
func dblogin(name string, password string) Users {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("login失敗")
	}
	defer db.Close()
	var users Users
	if err := db.Where("name = ? AND password = ?", name, password).First(&users).Error; err != nil {
		panic("存在しないアカウント")
	}
	return users
}

func dbDelete(id int) Users {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Delete失敗")
	}
	defer db.Close()
	var users Users
	db.First(&users, id)
	db.Delete(&users)
	return users
}

func dbGetOne(id int) Users {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("GetAll失敗")
	}
	defer db.Close()
	var users Users
	db.First(&users, id)
	return users
}
