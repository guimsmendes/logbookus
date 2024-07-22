package model

func GetModels() []interface{} {
	return []interface{}{
		&City{},
		&Country{},
		&Expense{},
		&User{},
		&UserCity{},
	}
}
