package models

type BeneficiaryParty struct {
	ID                   uint64 `json:"-" gorm:"primary_key"`
	*DebtorPartySkeleton `gorm:"embedded"`
	AccountType          int `json:"account_type"`
}
