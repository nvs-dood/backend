package models

type Student struct {
	ID     string   `json:"id" gorm:"primary_key;not null;"`
	Name   string   `json:"name"`
	Shifts []*Shift `json:"shifts" gorm:"foreignkey:StudentID"`
	Role   Role     `json:"role"`
	Class  string   `json:"class"`
}
