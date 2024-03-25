package domain

import "time"

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	LastName    string    `json:"last_name"`
	Surname     string    `json:"surname"`
	Gender      string    `json:"gender"`
	Status      string    `json:"status"`
	DateOfBirth time.Time `json:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at"`
}

type UsersParam struct {
	Gender   string
	Status   string
	OrderBy  string
	OrderDir string
	FullName string
	Limit    int64
	Offset   int64
}
