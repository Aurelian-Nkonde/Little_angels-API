package models

import "gorm.io/gorm"

type Guardian struct {
	gorm.Model
	FirstName   string
	LastName    string
	Email       *string
	Proffession *string
	Gender      string
}
