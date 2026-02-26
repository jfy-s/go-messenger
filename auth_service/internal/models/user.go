package models

import "time"

type User struct {
	Id         uint64
	Username   string `validate:"required,min=5,max=20"`
	Password   string `validate:"required"`
	Created_at time.Time
	Updated_at time.Time
}
