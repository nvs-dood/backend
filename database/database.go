package database

import (
	"fmt"

	"github.com/EnglederLucas/nvs-dood/graph/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	InitDB()
	return db
}

func InitDB() {
	var err error
	//dataSourceName := "root:@tcp(172.17.0.2:3306)/?parseTime=True"
	connectionString := "root:Pass@(localhost:3306)/" //test_db?charset=utf8&parseTime=True"

	db, err = gorm.Open("mysql", connectionString)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.LogMode(true)

	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	db.Exec("DROP DATABASE test_db;")
	db.Exec("CREATE DATABASE test_db;")
	db.Exec("USE test_db;")

	// Migration to create tables for Order and Item schema
	db.AutoMigrate(&models.Student{}, &models.Shift{})
}
