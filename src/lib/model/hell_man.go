package model

import (
	"time"
)

type (
	HellMan struct {
		TelegramId int       `json:"telegram_id"`
		Username   string    `json:"username"`
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		EnrollAt   time.Time `json:"enroll_at"`
	}
)
