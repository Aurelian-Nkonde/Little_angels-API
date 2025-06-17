package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	FirstName string
	LastName  string
	Gender    string
	Age       uint8
	Email     *string
	Role      string
}
