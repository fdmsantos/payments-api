package models

type DebtorParty struct {
	ID                   uint64 `json:"-" gorm:"primary_key"`
	*DebtorPartySkeleton `gorm:"embedded"`
}

type DebtorPartySkeleton struct {
	*SponsorPartySkeleton `gorm:"embedded"`
	AccountName           string `json:"account_name"`
	AccountNumberCode     string `json:"account_number_code"`
	Address               string `json:"address"`
	Name                  string `json:"name"`
}
