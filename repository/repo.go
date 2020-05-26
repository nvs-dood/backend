package repository

import (
	"github.com/EnglederLucas/nvs-dood/graph/models"
	"github.com/jinzhu/gorm"
)

func GetStudentByID(db *gorm.DB, studentID int) (*models.Student, error) {
	var student models.Student
	err := db.Preload("Shifts").Where("id = ?", studentID).Find(&student).Error

	if err != nil {
		return nil, err
	}

	return &student, nil
}
