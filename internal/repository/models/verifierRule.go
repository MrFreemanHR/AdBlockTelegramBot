package models

import (
	"time"
)

type VerifierRule struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" gorm:"uniqueIndex"`
	Author      string    `json:"author" gorm:"index"`
	AudioNote   bool      `json:"audio_note"`
	VideoNote   bool      `json:"video_note"`
	Photo       bool      `json:"photo"`
	Text        string    `json:"text_regexp"`
	LocaleGroup string    `json:"locale_group"`
	LocaleKey   string    `json:"locale_key"`
}
