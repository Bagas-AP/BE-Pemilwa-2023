package model

type CalonSenat struct {
	IDSenat   int    `gorm:"primaryKey" json:"idSenat"`
	NamaSenat string `json:"namaSenat"`
	Foto      string `json:"foto"`
}
