package model

import "guimsmendes/sniptech-company/deals-client-lib/pkg/model/time"

type Expense struct {
	Id     int64
	Name   string
	UserId int
	CityId int
	Type   ExpenseType
	Cost   float64
	Date   time.Time
}

type ExpenseType string

const (
	Transport ExpenseType = "transport"
	Tour      ExpenseType = "tour"
	Hotel     ExpenseType = "hotel"
	Local     ExpenseType = "local"
	Home      ExpenseType = "home"
)
