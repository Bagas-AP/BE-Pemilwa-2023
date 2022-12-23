package model

import "time"

type Users struct {
	NIM           string `gorm:"primaryKey" json:"nim"`
	Nama          string `json:"nama"`
	Foto          string `json:"foto"`
	IsVote        bool   `gorm:"default:false" json:"isVote"`
	ISAdmin       bool   `gorm:"default:false" json:"isAdmin"`
	CalonKepalaID int    `gorm:"default:0" json:"calonKepala"`
	CalonSenatID  int    `gorm:"default:0" json:"calonSenat"`
	waktuVote     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
