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

	for _, item := range offers {
		res, err := i.ImportOffer(item.Offer, item.Receiver)
		fmt.Sprintf("This is res %v, and this is err: %v", res, err)
	}
}

func (i *ImportOfferHandler) ImportOffer(offerName string, buyerID string) (Status, error) {

	err := i.findBuyer(offerName, buyerID)
	if err != nil {
		return BuyerNotFound, err
	}

	isFound, err := i.findOffer(offerName, buyerID)
	if err != nil {
		return Failed, err
	}
	if isFound {
		//todo update offer (mvp 2)
		return OfferFound, nil
	}

	code, err := i.createOffer(offerName, buyerID)
	if err != nil {
		return Failed, fmt.Errorf("failed to create offer: %v", err)
	}

	log.Println("new offer created: %v", code)
	return OfferCreated, nil
}

func (i *ImportOfferHandler) findBuyer(offerName string, buyerID string) error {
	res, err := i.transport.GetBuyer(buyerID)
	if err != nil {
		return fmt.Errorf("buyer %v not found: %v", buyerID, err)
	}
	if fmt.Sprintf("%s", res["Connected"]) != "true" {
		return fmt.Errorf("buyer %v for offer %v not connected with supplier", buyerID, offerName)
	}
	return nil
}

func (i *ImportOfferHandler) findOffer(offerName string, buyerID string) (bool, error) {
	res, err := i.transport.SearchOffer(offerName)
	if err != nil {
		return false, fmt.Errorf("failed to find offer: %v", err)
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
		return "", fmt.Errorf("failed to create offer: %v", err)
	}
	return offerCode, nil
}
