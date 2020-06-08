package models

import "time"

type Student struct {
	ID     int      `json:"id" gorm:"auto_increment;not null;primary_key"`
	Name   string   `json:"name"`
	Shifts []*Shift `json:"shifts" gorm:"foreignkey:StudentID"`
	Role   Role     `json:"role"`
	Class  string   `json:"class"`
}

type Shift struct {
	StudentID int        `json:"-"` //Is student_id in database
	Start     *time.Time `json:"start"`
	End       *time.Time `json:"end"`
}
