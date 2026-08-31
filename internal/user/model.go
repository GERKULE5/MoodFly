package user

import "time"

type User struct {
	ID          int        `json:"id"`
	Username    string     `json:"username"`
	Password    string     `json:"-"`
	PhoneNumber string     `json:"phone_number"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
