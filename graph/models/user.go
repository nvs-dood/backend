package models

type User struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      *string `json:"name"`
	StudentID *string `json:"studentID"`
	Admin     bool    `json:"admin"`
}
