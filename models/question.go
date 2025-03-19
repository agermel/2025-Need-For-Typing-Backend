package models

type Question struct {
	Title   string `json:"title" gorm:"unique;not null"`
	Context string `json:"context" gorm:"not null"`
}
