package model

type CalonKepala struct {
	IDKepala      int    `json:"idKepala"`
	Nama          string `json:"nama"`
	Foto          string `json:"foto"`
	CalonKepalaID int
	User          Users `gorm:"foreignKey:CalonKepalaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
