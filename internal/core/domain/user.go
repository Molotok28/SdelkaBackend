package domain

type User struct {
	ID          int
	Version     int
	Name        string
	Surname     string
	PhoneNumber *string
}
