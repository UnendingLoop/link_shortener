// Package model provides data-models for shortener-service
package model

import "time"

const Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Link struct {
	LID       int       `json:"-"`                  //
	ShortKey  string    `json:"shortkey,omitempty"` // ключ сокращенной (внутренней) ссылки; используется значение клиента или генерируется
	Redirect  string    `json:"redirect"`           // КУДА перенаправить клиента, mandatory
	CreatedAt time.Time `json:"created_at"`
	// Referrals []Referral `json:"referrals,omitempty"` //
}

type Referral struct {
	LID       int       `json:"-"`          // внешний ключ на Link.LID
	RID       int       `json:"-"`          //
	CreatedAt time.Time `json:"created_at"` // когда был переход по редиректу
	UserAgent string    `json:"user-agent"` // какой был юзерагент у клиента в формате "Chrome/Windows"
}
