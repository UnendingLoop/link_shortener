// Package model provides data-models for shortener-service
package model

import "time"

type Link struct {
	LID      uint      `gorm:"primaryKey;autoIncrement" json:"-"`              //праймари, автоинкремент, OnDelete:CASCADE - при удалении Link должны удаляться все Referal из привязанного массива
	ShortKey string    `gorm:"not null;uniqueIndex" json:"shortkey,omitempty"` //ключ сокращенной (внутренней) ссылки; используется значение клиента или генерируется
	OutLink  string    `gorm:"not null" json:"outlink"`                        //КУДА перенаправить клиента
	Referals []Referal `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:LID;references:LID" json:"referals,omitempty"`
}

type Referal struct {
	LID       uint      `gorm:"index;not null" json:"-"`           //внешний ключ на Link.LID
	RID       uint      `gorm:"primaryKey;autoIncrement" json:"-"` //праймари, автоинкремент
	Timestamp time.Time `json:"time"`
	UserAgent string    `json:"user-agent"`
}
