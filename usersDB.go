package main

import (
	"fmt"

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
func dbInit_users() error {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return fmt.Errorf("dbInit_users失敗: %w", err)
	}
	defer db.Close()
	return db.AutoMigrate(&Users{}).Error
}

// サインアップ
func dbSignup(name string, password string) error {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return fmt.Errorf("dbSignup失敗: %w", err)
	}
	defer db.Close()
	var users Users
	if err := db.Where("name = ?", name).First(&users).Error; err == nil {
		return fmt.Errorf("すでに同じ名前が使われています: %w", err)
	}
	return db.Create(&Users{Name: name, Password: password}).Error
}

// ログイン
func dblogin(name string, password string) (Users, error) {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return Users{}, fmt.Errorf("login失敗: %w", err)
	}
	defer db.Close()
	var users Users
	if err := db.Where("name = ? AND password = ?", name, password).First(&users).Error; err != nil {
		return Users{}, fmt.Errorf("存在しないアカウント: %w", err)
	}
	return users, nil
}

func dbDelete(id int) (Users, error) {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return Users{}, fmt.Errorf("dbDelete失敗: %w", err)
	}
	defer db.Close()
	var users Users
	err = db.First(&users, id).Error
	if err != nil {
		return Users{}, fmt.Errorf("dbDelete失敗: %w", err)
	}
	return users, nil
}

func dbGetOne(id int) (Users, error) {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return Users{}, fmt.Errorf("dbGetOne失敗: %w", err)
	}
	defer db.Close()
	var users Users
	err = db.First(&users, id).Error
	if err != nil {
		return Users{}, fmt.Errorf("dbGetOne失敗: %w", err)
	}
	return users, nil
}
