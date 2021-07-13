package offerReader

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"ts/adapters"
)

type OfferReader struct {
	fileManager *adapters.FileManager
	reader      adapters.HandlerInterface
}

type RawOffer struct {
	Offer     string
	Receiver  string
	Contract  string
	ValidFrom string
	ExpiresAt string
	Countries string
}

type Deps struct {
	dig.In
	Reader      adapters.HandlerInterface
	FileManager *adapters.FileManager
}

func NewOfferReader(deps Deps) *OfferReader {
	return &OfferReader{
		reader:      deps.Reader,
		fileManager: deps.FileManager,
	}
}

func (o *OfferReader) UploadOffers(path string) []RawOffer {
	ext := o.fileManager.GetFileType(path)
	o.reader.Init(ext)
	parsedRaws := o.reader.Parse(path)
	actualHeader := o.reader.GetHeader()
	header, err := processHeader(actualHeader)
	if err != nil {
		log.Printf("failed to upload offers: %v", err)
		return nil
	}
	or := processOffers(parsedRaws, header)
	log.Printf("Offers upload finished.")
	return or

}

func processOffers(raws []map[string]interface{}, header *RawHeader) []RawOffer {
	res := make([]RawOffer, len(raws))
	for i, item := range raws {
		offer := processOffer(header, item)
		if offer != nil {
			res[i] = *offer
		}
	}
	return res
}

func processOffer(header *RawHeader, row map[string]interface{}) *RawOffer {
	if isEmptyRow(row) {
		return nil
	}
	if row[header.Offer] == nil || row[header.Receiver] == nil {
		log.Printf("row does not contain values in required columns (Offer, Receiver). Actual value: %v", row)
		return nil
	}

	offer := RawOffer{
		Offer:    fmt.Sprintf("%v", row[header.Offer]),
		Receiver: fmt.Sprintf("%v", row[header.Receiver]),
	}
	if header.ContractID != "" {
		offer.Contract = fmt.Sprintf("%v", row[header.ContractID])
	}
	if header.ValidFrom != "" {
		offer.ValidFrom = fmt.Sprintf("%v", row[header.ValidFrom])
	}
	if header.ExpiresAt != "" {
		offer.ExpiresAt = fmt.Sprintf("%v", row[header.ExpiresAt])
	}
	if header.Countries != "" {
		offer.Countries = fmt.Sprintf("%v", row[header.Countries])
	}
	return &offer
}

func isEmptyRow(row map[string]interface{}) bool {
	for _, value := range row {
		if value != nil && value != "" {
			return false
		}
	}
	return true
}

func processHeader(parsedHeader []string) (*RawHeader, error) {
	resHeader := NewHeader(parsedHeader)
	if err := resHeader.ValidateHeader(); err != nil {
		return nil, err
	}
	return resHeader, nil
}
