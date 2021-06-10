package offerReader

import (
	"fmt"
	"ts/utils"
)

const (
	defaultOfferID    = "Offer"
	defaultReceiver   = "Receiver"
	defaultContractID = "Contract ID"
	defaultValidFrom  = "Valid From"
	defaultExpiresAt  = "ExpiresAt"
	defaultCountries  = "Countries"
)

type RawHeader struct {
	OfferID    string
	Receiver   string
	ContractID string
	ValidFrom  string
	ExpiresAt  string
	Countries  string
}

func NewHeader(input []string) *RawHeader {
	var newHeader RawHeader
	for _, columnLabel := range input {
		//required
		trimmedColumnLabel := utils.TrimAll(columnLabel)
		switch trimmedColumnLabel {
		case utils.TrimAll(defaultOfferID):
			newHeader.OfferID = columnLabel
		case utils.TrimAll(defaultReceiver):
			newHeader.Receiver = columnLabel
		//unrequired
		case utils.TrimAll(defaultContractID):
			newHeader.ContractID = columnLabel
		case utils.TrimAll(defaultValidFrom):
			newHeader.ValidFrom = columnLabel
		case utils.TrimAll(defaultExpiresAt):
			newHeader.ExpiresAt = columnLabel
		case utils.TrimAll(defaultCountries):
			newHeader.Countries = columnLabel
		}
	}
	return &newHeader
}

func (rh *RawHeader) ValidateHeader() error {
	if rh.OfferID == "" || rh.Receiver == "" {
		return fmt.Errorf("offers file does not contains all requiered fields: actual [OfferID: %v, Receiver: %v]", rh.OfferID, rh.Receiver)
	}
	return nil
}
