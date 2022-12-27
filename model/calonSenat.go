package model

type CalonSenat struct {
	IDSenat int    `gorm:"primaryKey" json:"idSenat"`
	Nama    string `json:"nama"`
	Foto    string `json:"foto"`
}
