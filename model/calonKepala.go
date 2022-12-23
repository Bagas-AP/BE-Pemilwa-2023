package model

type CalonKepala struct {
	IDKepala uint   `gorm:"primaryKey" json:"idKepala"`
	Nama     string `json:"nama"`
	Foto     string `json:"foto"`
}
