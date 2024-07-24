package model

import "time"

type UserCity struct {
	UserID    int
	CityID    int
	Status    Status
	StartDate time.Time
	EndDate   time.Time
}
