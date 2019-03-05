package models

type Charge struct {
	ID                   uint64 `json:"-" gorm:"primary_key"`
	ChargesInformationID uint64 `json:"-"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
}
