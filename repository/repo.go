package repository

import (
	"github.com/EnglederLucas/nvs-dood/graph/models"
	"github.com/jinzhu/gorm"
)

func GetStudentByID (db *gorm.DB, studentID string) (*models.Student, error){
	var student models.Student
	err := db.Where("studentID = ?", studentID).Find(&student).Error

	if err != nil {
		return nil, err
	}

	return student, nil
}
