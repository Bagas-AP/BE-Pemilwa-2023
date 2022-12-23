package model

type CalonSenat struct {
	IDSenat      int    `json:"idKepala"`
	Nama         string `json:"nama"`
	Foto         string `json:"foto"`
	CalonSenatID int
	User         Users `gorm:"foreignKey:CalonSenatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
