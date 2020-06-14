package database

import (
	"fmt"
	"os"
	"time"

	"github.com/EnglederLucas/nvs-dood/graph/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func GetDB() (*gorm.DB, error) {
	err := InitDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() error {
	var err error
	connectionString := os.Getenv("DB_CONN_STRING")
	if connectionString == "" {
		connectionString = "root:Pass@(localhost:3306)/" //test_db?charset=utf8&parseTime=True"
	}

	db, err = gorm.Open("mysql", connectionString)

	if err != nil {
		fmt.Println(err)
		return err
	}

	db.LogMode(true)

	// Create the database. This is a one-time step.
	// Comment out if running multiple times - You may see an error otherwise
	db.Exec("DROP DATABASE test_db;")
	db.Exec("CREATE DATABASE test_db;")
	db.Exec("USE test_db;")

	// Migration to create tables for Order and Item schema
	db.AutoMigrate(&models.Student{}, &models.Shift{})
	SeedDatabase()

	return nil
}

func SeedDatabase() {
	students := []models.Student{
		{
			ID:    "seed-niki",
			Class: "2CHIF",
			Name:  "Niki Kudrijaschow",
			Role:  models.RoleBuffet,
			Shifts: []*models.Shift{
				{
					StudentID: "seed-niki",
					Start:     &[]time.Time{(time.Now())}[0],
					End:       &[]time.Time{(time.Now().Add(time.Hour))}[0],
				},
			},
		},
		{
			ID:    "seed-beck",
			Class: "2CHIF",
			Name:  "Daniel Beck",
			Role:  models.RoleEmpfang,
			Shifts: []*models.Shift{
				{
					StudentID: "seed-beck",
					Start:     &[]time.Time{(time.Now().Add(1 * time.Hour))}[0],
					End:       &[]time.Time{(time.Now().Add(2 * time.Hour))}[0],
				},
				{
					StudentID: "seed-beck",
					Start:     &[]time.Time{(time.Now().Add(4 * time.Hour))}[0],
					End:       &[]time.Time{(time.Now().Add(6 * time.Hour))}[0],
				},
			},
		},
		{Class: "2CHIF", Name: "Spasenovic Bozidar", Role: models.RoleGUIDE, ID: "seed-bozidar"},
		{Class: "2CHIF", Name: "Supper Marco", Role: models.RoleGUIDE, ID: "seed-supper"},
		{Class: "2CHIF", Name: "Tanzer Rafael", Role: models.RoleGUIDE, ID: "seed-tanzer"},
		{Class: "2CHIF", Name: "Weidinger	Daniel", Role: models.RoleGUIDE, ID: "seed-weidi"},
		{Class: "2CHIF", Name: "Wirth	Lukas", Role: models.RoleGUIDE, ID: "seed-wirth"},
	}

	for _, s := range students {
		shifts := s.Shifts
		err := db.Save(s).Error
		if err != nil {
			panic("could not insert seed student")
		}

		for _, shift := range shifts {
			err := db.Save(shift).Error
			if err != nil {
				panic("could not insert seed student's shift")
			}
		}
	}
}
