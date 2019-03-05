package models

type ChargesInformation struct {
	ID                      uint64   `json:"-" gorm:"primary_key"`
	BearerCode              string   `json:"bearer_code"`
	SenderCharges           []Charge `json:"sender_charges"`
	ReceiverChargesAmount   string   `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string   `json:"receiver_charges_currency"`
}
