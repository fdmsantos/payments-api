package models

type SponsorParty struct {
	ID                    uint64 `json:"-" gorm:"primary_key"`
	*SponsorPartySkeleton `gorm:"embedded"`
}

type SponsorPartySkeleton struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}
