package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"payments/infrastructure"
	"payments/utils"
)

type Payment struct {
	Type           string     `json:"type"`
	ID             uuid.UUID  `gorm:"primary_key" json:"id" sql:",type:uuid"`
	Version        uint       `json:"version"`
	OrganisationID uuid.UUID  `json:"organisation_id" sql:",type:uuid"`
	Attributes     Attributes `json:"attributes" gorm:"foreignkey:PaymentRefer"`
}

// GetPaymentByID Get a payment model through an ID
func GetPaymentByID(id uuid.UUID) (Payment, error) {
	payment := Payment{}
	if err := infrastructure.GetDB().Set("gorm:auto_preload", true).Where("ID = ? ", id).First(&payment).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return payment, errors.New(utils.ERROR_RESOURCE_NOT_FOUND)
		}
		return payment, errors.New(utils.ERROR_SERVER)
	}
	return payment, nil
}
