package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	FirstName string
	LastName  string
	Gender    string
	Age       uint8
	Guardian  *int
}
