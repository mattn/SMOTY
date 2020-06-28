package main

import (
	"fmt"

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
func dbInit_server() error {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return fmt.Errorf("dbInit_server失敗: %w", err)
	}
	defer db.Close()
	return db.AutoMigrate(&Problem_server{}).Error
}

func check_server(id int, anser string) (Problem_server, string, error) {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return Problem_server{}, "", fmt.Errorf("check_server失敗: %w", err)
	}
	defer db.Close()
	var result string
	var server Problem_server
	if err := db.Where("id = ? AND anser = ?", id, anser).First(&server).Error; err != nil {
		result = "不正解"
	} else {
		result = "正解"
	}
	return server, result, nil
}

func serverGetAll() ([]Problem_server, error) {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return nil, fmt.Errorf("データベース開けず(dbGetAll): %w", err)
	}
	defer db.Close()
	var server []Problem_server
	err = db.Order("created_at desc").Find(&server).Error
	if err != nil {
		return nil, fmt.Errorf("serverGetAll失敗: %w", err)
	}
	return server, nil
}

func serverGetOne(id int) (Problem_server, error) {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return Problem_server{}, fmt.Errorf("データベース開けず(dbGetOne): %w", err)
	}
	defer db.Close()
	var server Problem_server
	err = db.First(&server, id).Error
	if err != nil {
		return Problem_server{}, fmt.Errorf("serverGetOne失敗: %w", err)
	}
	return server, nil
}

func serverInsert(question string, anser string, hint string) error {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return fmt.Errorf("serverInsert失敗: %w", err)
	}
	defer db.Close()
	return db.Create(&Problem_server{Question: question, Anser: anser, Hint: hint}).Error
}

func serverUpdate(id int, question string, hint string, anser string) error {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return fmt.Errorf("serverUpdate失敗: %w", err)
	}
	defer db.Close()
	var server Problem_server
	err = db.First(&server, id).Error
	if err != nil {
		return fmt.Errorf("serverUpdate失敗: %w", err)
	}
	server.Question = question
	server.Hint = hint
	server.Anser = anser
	return db.Save(&server).Error
}

func serverDelete(id int) error {
	db, err := gorm.Open("mysql", connectString)
	if err != nil {
		return fmt.Errorf("serverDelete失敗: %w", err)
	}
	defer db.Close()
	var server Problem_server
	return db.Where("id = ?", id).Delete(&server).Error
}
