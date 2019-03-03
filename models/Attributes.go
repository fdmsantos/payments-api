package models

import "github.com/satori/go.uuid"

type Attributes struct {
	ID                   uint64             `json:"-" gorm:"primary_key"`
	PaymentRefer         uuid.UUID          `json:"-" sql:",type:uuid"`
	Amount               string             `json:"amount"`
	BeneficiaryParty     BeneficiaryParty   `json:"beneficiary_party"`
	BeneficiaryPartyID   uint64             `json:"-"`
	ChargesInformation   ChargesInformation `json:"charges_information"`
	ChargesInformationID uint64             `json:"-"`
	Currency             string             `json:"currency"`
	DebtorParty          DebtorParty        `json:"debtor_party"`
	DebtorPartyID        uint64             `json:"-"`
	EndToEndReference    string             `json:"end_to_end_reference"`
	FX                   FX                 `json:"fx"`
	FXID                 uint64             `json:"-"`
	NumericReference     string             `json:"numeric_reference"`
	PaymentID            string             `json:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme"`
	PaymentType          string             `json:"payment_type"`
	ProcessingDate       string             `json:"processing_date"`
	Reference            string             `json:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type"`
	SponsorParty         SponsorParty       `json:"sponsor_party"`
	SponsorPartyID       uint64             `json:"-"`
}
