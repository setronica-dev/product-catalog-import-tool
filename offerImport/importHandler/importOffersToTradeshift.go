package importHandler

import (
	"fmt"
	"log"
	"time"
	"ts/config/configModels"
	"ts/externalAPI/tradeshiftAPI"
	"ts/offerImport/offerReader"
)

type ImportOfferHandler struct {
	transport  *tradeshiftAPI.TradeshiftAPI
	Recipients *configModels.Recipients
}

func NewImportOfferHandler(deps Deps) ImportOfferInterface {
	return &ImportOfferHandler{
		transport:  deps.Transport,
		Recipients: deps.Config.TradeshiftAPI.Recipients,
	}
}

func (i *ImportOfferHandler) ImportOffers(offers []offerReader.RawOffer) {
	log.Printf("IMPORT OFFERS TO TRADESHIFT WAS STARTED")
	for _, offer := range offers {
		if err := validateOffer(offer); err != nil {
			log.Printf("failed to import offer \"%v\". Reason:  %v", offer, err)
			break
		}
		_, err := i.ImportOffer(
			offer.Offer,
			offer.Receiver,
			offer.ValidFrom,
			offer.ExpiresAt,
			offer.Countries)
		if err != nil {
			log.Printf("failed to import offer \"%v\". Reason:  %v", offer.Offer, err)
		}
	}
}

func validateOffer(offer offerReader.RawOffer) error {
	if offer.Offer == "" {
		return fmt.Errorf("offer name should be defined")
	}
	if offer.Receiver == "" {
		return fmt.Errorf("offer receiver should be defined")
	}
	return nil
}

func (i *ImportOfferHandler) ImportOffer(
	offerName string,
	recipientName string,
	startDate *time.Time,
	endDate *time.Time,
	countries []string) (Status, error) {
	buyerID := i.Recipients.GetRecipientIDByName(recipientName)
	if buyerID == "" {
		return Failed, fmt.Errorf("failed to find buyer %v in config", recipientName)
	}
	isFound, err := i.isBuyerExists(buyerID)
	if err != nil {
		return Failed, err
	}
	if !isFound {
		return BuyerNotFound, nil
	}

	isFound, err = i.isOfferExists(offerName, buyerID)
	if err != nil {
		return Failed, err
	}
	if isFound {
		log.Printf("offer \"%v\" was found", offerName)
		return OfferFound, nil
	}

	_, err = i.createOffer(offerName, buyerID, startDate, endDate, countries)
	if err != nil {
		return Failed, err
	}

	log.Printf("New offer with name %v has been created", offerName)
	return OfferCreated, nil
}

func (i *ImportOfferHandler) isBuyerExists(buyerID string) (bool, error) {
	res, err := i.transport.GetBuyer(buyerID)
	if err != nil {
		return false, fmt.Errorf("buyer \"%v\" not found: %v", buyerID, err)
	}
	if fmt.Sprintf("%v", res["Connected"]) != "true" {
		return false, nil
	}
	return true, nil
}

func (i *ImportOfferHandler) isOfferExists(offerName string, buyerID string) (bool, error) {
	res, err := i.transport.SearchOffer(offerName)
	if err != nil {
		return false, fmt.Errorf("failed to find offer \"%v\": %v", offerName, err)
	}
	if res["total"].(float64) > 0.0 {
		data := res["data"].([]interface{})
		for _, item := range data {
			itemBuyerId := item.(map[string]interface{})["buyerId"]
			if fmt.Sprintf("%s", itemBuyerId) == buyerID {
				return true, fmt.Errorf("offer %v for buyer %v is allready exists", offerName, buyerID)
			}
		}
	}
	return false, nil
}

func (i *ImportOfferHandler) createOffer(
	offerName string,
	buyerID string,
	startDate *time.Time,
	endDate *time.Time,
	countries []string) (string, error) {
	offerID, err := i.transport.CreateOffer(offerName, buyerID)
	if err != nil {
		return "", fmt.Errorf("offer %v was not created: %v", offerName, err)
	}
	err = i.transport.UpdateOffer(offerID, offerName, startDate, endDate, countries)
	if err != nil {
		return offerID, fmt.Errorf("offer %v was not updated: %v", offerID, err)
	}
	return offerID, nil
}
