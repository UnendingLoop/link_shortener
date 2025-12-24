// Package model provides data-models for shortener-service
package model

import "time"

type Link struct {
	LID       uint       `json:"-"`                   //
	ShortKey  string     `json:"shortkey,omitempty"`  // ключ сокращенной (внутренней) ссылки; используется значение клиента или генерируется
	Redirect  string     `json:"redirect"`            // КУДА перенаправить клиента, mandatory
	Referrals []Referral `json:"referrals,omitempty"` //
	RentEnd   *time.Time `json:"rentend"`             // конец аренды редиректа
	IsActive  bool       `json:"isactive"`            // флаг активности редиректа, false при наступлении конца аренды
}

type Referral struct {
	LID       uint      `json:"-"`          // внешний ключ на Link.LID
	RID       uint      `json:"-"`          //
	Created   time.Time `json:"time"`       // когда был переход по редиректу
	UserAgent string    `json:"user-agent"` // какой был юзерагент у клиента
}
