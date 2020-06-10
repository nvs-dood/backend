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

func GetUserByID(db *gorm.DB, userId string) (*models.User, error) {
	var user models.User
	err := db.Where("id = ?", userId).Find(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func InsertUser(db *gorm.DB, user *models.User) (*models.User, error) {
	err := (*db).Create(&user).Error
	if err != nil {
		return nil, err
	}

	db.Commit()

	dbUser, err := GetUserByID(db, user.ID)
	if err != nil {
		return nil, err
	}

	return dbUser, nil
}
