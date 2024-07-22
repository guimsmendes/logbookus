package model

import "time"

type UserCity struct {
	UserId    int
	CityId    int
	Status    Status
	StartDate time.Time
	EndDate   time.Time
}
