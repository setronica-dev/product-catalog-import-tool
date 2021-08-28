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
	transport        *tradeshiftAPI.TradeshiftAPI
	recipientsConfig *configModels.Recipients
}

func NewImportOfferHandler(deps Deps) ImportOfferInterface {
	return &ImportOfferHandler{
		transport:        deps.Transport,
		recipientsConfig: deps.Config.TradeshiftAPI.Recipients,
	}
}

func (i *ImportOfferHandler) ImportOffers(offers []offerReader.RawOffer) {

	log.Printf("Import offers to Tradeshift has been started")
	for _, offer := range offers {
		if err := validateOffer(offer); err != nil {
			log.Printf("failed to import offer \"%v\". Reason:  %v", offer, err)
			break
		}
		_, err := i.importOffer(
			offer.Offer,
			offer.ReceiverName,
			offer.ValidFrom,
			offer.ExpiresAt,
			offer.Countries)
		if err != nil {
			log.Printf("Offer Import occured the error: %v", err)
		}
	}
}

func validateOffer(offer offerReader.RawOffer) error {
	if offer.Offer == "" {
		return fmt.Errorf("offer name should be defined")
	}
	if offer.ReceiverName == "" {
		return fmt.Errorf("offer receiver should be defined")
	}
	return nil
}

func (i *ImportOfferHandler) importOffer(
	offerName string,
	recipientName string,
	startDate time.Time,
	endDate time.Time,
	countries []string) (Status, error) {
	recipientID := i.recipientsConfig.GetRecipientIDByName(recipientName)
	if recipientID == "" {
		return Failed, fmt.Errorf("failed to find buyer %v in config", recipientName)
	}
	isFound, err := i.isRecipientExists(recipientID)
	if err != nil {
		return Failed, err
	}
	if !isFound {
		log.Printf("Offer '%v' can't be created for unknown buyer '%v'", offerName, recipientID)
		return BuyerNotFound, nil
	}

	offer, err := i.findOfferByNameAndBuyer(offerName, recipientID)
	if err != nil {
		return Failed, err
	}
	if offer != nil {
		err := i.updateOffer(offer, startDate, endDate, countries)
		if err != nil {
			return Failed, fmt.Errorf("failed to update offer %v: %v", offerName, err)
		}
		return OfferFound, nil
	}
	_, err = i.createOffer(offerName, recipientID, startDate, endDate, countries)
	if err != nil {
		return Failed, err
	}

	log.Printf("New offer with name %v has been created", offerName)
	return OfferCreated, nil
}

func (i *ImportOfferHandler) isRecipientExists(recipientID string) (bool, error) {
	res, err := i.transport.GetBuyer(recipientID)
	if err != nil {
		return false, fmt.Errorf("buyer \"%v\" not found: %v", recipientID, err)
	}
	if fmt.Sprintf("%v", res["Connected"]) != "true" {
		return false, nil
	}
	return true, nil
}
