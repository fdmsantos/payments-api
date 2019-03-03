package models

import (
	"github.com/satori/go.uuid"
)

type Payment struct {
	Type           string     `json:"type"`
	ID             uuid.UUID  `gorm:"primary_key" json:"id" sql:",type:uuid"`
	Version        uint       `json:"version"`
	OrganisationID uuid.UUID  `json:"organisation_id" sql:",type:uuid"`
	Attributes     Attributes `json:"attributes" gorm:"foreignkey:PaymentRefer"`
}
