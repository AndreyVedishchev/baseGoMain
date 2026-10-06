package models

import "time"

type Resume struct {
	UUID        string
	Fio         *string
	Birthdate   *time.Time
	Sex         *string
	City        *string
	Telephon    *string
	Citizenship *string
	Position    *string
	Salary      *float64
	Experience  *int
}
