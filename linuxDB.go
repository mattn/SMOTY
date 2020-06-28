package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Problem_linux struct {
	gorm.Model
	Question string
	Hint     string
	Anser    string
}

// DB接続
func dbInit_linux() {
	db, err := gorm.Open("mysql", "root:@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbInit_linux失敗")
	}
	defer db.Close()
	db.AutoMigrate(&Problem_linux{})
}

func check_linux(id int, anser string) (Problem_linux, string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("linux_check失敗")
	}
	defer db.Close()
	var result string
	var linux Problem_linux
	if err := db.Where("id = ? AND anser = ?", id, anser).First(&linux).Error; err != nil {
		result = "不正解"
	} else {
		result = "正解"
	}
	return linux, result
}

func linuxGetAll() []Problem_linux {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("データベース開けず(dbGetAll)")
	}
	defer db.Close()
	var linux []Problem_linux
	db.Order("created_at desc").Find(&linux)
	return linux
}

func linuxGetOne(id int) Problem_linux {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("データベース開けず(dbGetOne)")
	}
	defer db.Close()
	var linux Problem_linux
	db.First(&linux, id)
	return linux
}

func linuxInsert(question string, anser string, hint string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("linuxInsert失敗")
	}
	defer db.Close()
	db.Create(&Problem_linux{Question: question, Anser: anser, Hint: hint})
}

func linuxUpdate(id int, question string, hint string, anser string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("linuxUpdate失敗")
	}
	defer db.Close()
	var linux Problem_linux
	db.First(&linux, id)
	linux.Question = question
	linux.Anser = anser
	linux.Hint = hint
	db.Save(&linux)
}

func linuxDelete(id int) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("linuxDelete失敗")
	}
	defer db.Close()
	var linux Problem_linux
	db.Where("id = ?", id).Delete(&linux)
}
