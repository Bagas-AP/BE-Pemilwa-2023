package model

import "time"

type Users struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	NIM           string      `gorm:"uniqueIndex" json:"nim"`
	Nama          string      `json:"nama"`
	Foto          string      `json:"foto"`
	IsVote        bool        `gorm:"default:false" json:"isVote"`
	ISAdmin       bool        `gorm:"default:false" json:"isAdmin"`
	CalonKepalaID int         `gorm:"default:null" json:"calonKepala"`
	CalonSenatID  int         `gorm:"default:null" json:"calonSenat"`
	WaktuVote     time.Time   `json:"waktuVote"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
	DeletedAt     time.Time   `json:"deletedAt"`
	CalonSenat    CalonSenat  `gorm:"foreignKey:CalonSenatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CalonKepala   CalonKepala `gorm:"foreignKey:CalonKepalaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
