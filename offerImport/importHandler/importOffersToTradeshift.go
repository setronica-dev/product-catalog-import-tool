package importHandler

import (
	"fmt"
	"log"
	"ts/externalAPI/tradeshiftAPI"
	"ts/offerImport/offerReader"
)

type ImportOfferHandler struct {
	transport *tradeshiftAPI.TradeshiftAPI
}

func NewImportOfferHandler(deps Deps) ImportOfferInterface {
	return &ImportOfferHandler{
		transport: deps.Transport,
	}
}

func (i *ImportOfferHandler) ImportOffers(offers []offerReader.RawOffer) {
	for _, offer := range offers {
		if err := validateOffer(offer); err != nil {
			log.Printf("failed to import offer \"%v\". Reason:  %v", offer, err)
			break
		}
		_, err := i.ImportOffer(offer.Offer, offer.Receiver)
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

func (i *ImportOfferHandler) ImportOffer(offerName string, buyerID string) (Status, error) {
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
		log.Print("offer \"%v\" was found", offerName)
		return OfferFound, nil
	}

	code, err := i.createOffer(offerName, buyerID)
	if err != nil {
		return Failed, fmt.Errorf("offer was not created: %v", err)
	}

	log.Printf("new offer was created: \"%v\"", code)
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

func (i *ImportOfferHandler) createOffer(offerName string, buyerID string) (string, error) {
	offerCode, err := i.transport.CreateOffer(offerName, buyerID)
	if err != nil {
		return "", err
	}
	return offerCode, nil
}
