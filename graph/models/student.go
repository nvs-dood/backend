package models

import "time"

type Student struct {
	ID 	   string   `json:"id" gorm:"primary_key"`
	Name   string   `json:"name"`
	Shifts []*Shift `json:"shifts" gorm:"foreignkey:StudentID"`
	Role   Role     `json:"role"`
	Class  string   `json:"class"`
}

type Shift struct {
	StudentID string     `json:"-"`
	Start     *time.Time `json:"start"`
	End       *time.Time `json:"end"`
}