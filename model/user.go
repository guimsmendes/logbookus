package model

type User struct {
	Id                int
	Name              string
	Email             string
	VisitedCountryIds []int
}
