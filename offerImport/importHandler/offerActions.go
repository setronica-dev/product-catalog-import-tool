package importHandler

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"
)

func (i *ImportOfferHandler) findOfferByNameAndBuyer(offerName string, buyerID string) (*Offer, error) {
	res, err := i.transport.SearchOffer(offerName)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer \"%v\": %v", offerName, err)
	}
	if res["total"].(float64) > 0.0 {
		data := res["data"].([]interface{})

		for _, item := range data {
			offer := newOffer(item.(map[string]interface{}))
			if offer.Receiver == buyerID {
				return offer, nil
			}
		}
	}
	return nil, nil
}

func (i *ImportOfferHandler) createOffer(
	offerName string,
	buyerID string,
	startDate time.Time,
	endDate time.Time,
	countries []string) (string, error) {
	offerID, err := i.transport.CreateOffer(offerName, buyerID)
	if err != nil {
		return "", fmt.Errorf("offer %v was not created: %v", offerName, err)
	}
	err = i.transport.UpdateOffer(offerID, offerName, &startDate, &endDate, countries)
	if err != nil {
		return offerID, fmt.Errorf("offer %v was not updated: %v", offerID, err)
	}
	return offerID, nil
}

func (i *ImportOfferHandler) updateOffer(
	offer *Offer,
	newStartDate time.Time,
	newEndDate time.Time,
	newCountries []string) error {

	sort.Strings(offer.Countries)
	sort.Strings(newCountries)

	if offer.ValidFrom.Unix() == newStartDate.Unix() &&
		offer.ExpiresAt.Unix() == newEndDate.Unix() &&
		reflect.DeepEqual(offer.Countries, newCountries) {
		log.Printf("Offer '%v' exists already with the same data and we skip it", offer.Name)
	} else {
		err := i.transport.UpdateOffer(offer.ID, offer.Name, &newStartDate, &newEndDate, newCountries)

		log.Printf("Offer '%v' exists already but it has changes and has been updated", offer.Name)
		if err != nil {
			return fmt.Errorf("offer %v was not updated: %v", offer.ID, err)
		}
	}
	return nil
}
