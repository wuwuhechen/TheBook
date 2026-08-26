package bank

//TODO

import "gorm.io/gorm"

type PracticeBankSQLite struct {
	practices *gorm.DB
}

func NewPracticeBankSQLite(db *gorm.DB) *PracticeBankSQLite {
	return &PracticeBankSQLite{
		practices: db,
	}
}
