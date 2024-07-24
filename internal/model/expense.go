package model

import "time"

type Expense struct {
	ID     int64
	Name   string
	UserID int
	CityID int
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
