package authmodels

type User struct {
	ID        string
	Email     string
	Name      *string
	StudentID *string
	Admin     bool
}
