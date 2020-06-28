package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Problem_server struct {
	gorm.Model
	Question string
	Hint     string
	Anser    string
}

// DB接続
func dbInit_server() {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbInit_server失敗")
	}
	defer db.Close()
	db.AutoMigrate(&Problem_server{})
}

func check_server(id int, anser string) (Problem_server, string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("server_check失敗")
	}
	defer db.Close()
	var result string
	var server Problem_server
	if err := db.Where("id = ? AND anser = ?", id, anser).First(&server).Error; err != nil {
		result = "不正解"
	} else {
		result = "正解"
	}
	return server, result
}

func serverGetAll() []Problem_server {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("データベース開けず(dbGetAll)")
	}
	defer db.Close()
	var server []Problem_server
	db.Order("created_at desc").Find(&server)
	return server
}

func serverGetOne(id int) Problem_server {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("データベース開けず(dbGetOne)")
	}
	defer db.Close()
	var server Problem_server
	db.First(&server, id)
	return server
}

func serverInsert(question string, anser string, hint string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("serverInsert失敗")
	}
	defer db.Close()
	db.Create(&Problem_server{Question: question, Anser: anser, Hint: hint})
}

func serverUpdate(id int, question string, hint string, anser string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("serverUpdate失敗")
	}
	defer db.Close()
	var server Problem_server
	db.First(&server, id)
	server.Question = question
	server.Hint = hint
	server.Anser = anser
	db.Save(&server)
}

func serverDelete(id int) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("serverDelete失敗")
	}
	defer db.Close()
	var server Problem_server
	db.Where("id = ?", id).Delete(&server)
}
