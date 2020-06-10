package repository

import (
	"github.com/EnglederLucas/nvs-dood/graph/models"
	"github.com/jinzhu/gorm"
)

func GetStudentByID(db *gorm.DB, studentID string) (*models.Student, error) {
	var student models.Student
	err := db.Preload("Shifts").Where("id = ?", studentID).Find(&student).Error

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func InsertStudent(db *gorm.DB, student *models.Student) (*models.Student, error) {

	err := (*db).Create(&student).Error
	if err != nil {
		return nil, err
	}

	dbStudent, err := GetStudentByID(db, student.ID)
	if err != nil {
		return nil, err
	}

	return dbStudent, nil
}
