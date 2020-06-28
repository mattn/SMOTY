package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Problem_router struct {
	gorm.Model
	Question string
	Hint     string
	Anser    string
}

// DB接続
func dbInit_router() {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("dbInit_router失敗")
	}
	defer db.Close()
	db.AutoMigrate(&Problem_router{})
}

func check_router(id int, anser string) (Problem_router, string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("router_check失敗")
	}
	defer db.Close()
	var result string
	var router Problem_router
	if err := db.Where("id = ? AND anser = ?", id, anser).First(&router).Error; err != nil {
		result = "不正解"
	} else {
		result = "正解"
	}
	return router, result
}

func routerGetAll() []Problem_router {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("データベース開けず(dbGetAll)")
	}
	defer db.Close()
	var router []Problem_router
	db.Order("created_at desc").Find(&router)
	return router
}

func routerGetOne(id int) Problem_router {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("データベース開けず(dbGetOne)")
	}
	defer db.Close()
	var router Problem_router
	db.First(&router, id)
	return router
}

func routerInsert(question string, anser string, hint string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("routerInsert失敗")
	}
	defer db.Close()
	db.Create(&Problem_router{Question: question, Anser: anser, Hint: hint})
}

func routerUpdate(id int, question string, hint string, anser string) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("routerUpdate失敗")
	}
	defer db.Close()
	var router Problem_router
	db.First(&router, id)
	router.Question = question
	router.Hint = hint
	router.Anser = anser
	db.Save(&router)
}

func routerDelete(id int) {
	db, err := gorm.Open("mysql", "root:password@/database_name?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("routerDelete失敗")
	}
	defer db.Close()
	var router Problem_router
	db.Where("id = ?", id).Delete(&router)
}
