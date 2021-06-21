package entity

import "gorm.io/gorm"

type NoteType string

const (
	Idea       = "Idea"
	Info       = "Info"
	Credential = "Credential"
	Reminder   = "Reminder"
	Plan       = "Plan"
	Jurnal     = "Jurnal"
)

type Note struct {
	gorm.Model
	Title     string
	Body      string
	Secret    string
	Type      NoteType
	AccountID uint
	Account   Account
}
