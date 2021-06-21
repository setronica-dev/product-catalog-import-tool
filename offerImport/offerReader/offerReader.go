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
		log.Fatalf("failed to upload rules: %v", err)
	}
	or := processOffers(parsedRaws, header)
	log.Printf("Offers upload finished.")
	return or

}

func processOffers(raws []map[string]interface{}, header *RawHeader) []RawOffer {
	res := make([]RawOffer, len(raws))
	for i, item := range raws {
		offer := RawOffer{
			Offer:    fmt.Sprintf("%v", item[header.Offer]),
			Receiver: fmt.Sprintf("%v", item[header.Receiver]),
		}
		if header.ContractID != "" {
			offer.Contract = fmt.Sprintf("%v", item[header.ContractID])
		}
		if header.ValidFrom != "" {
			offer.ValidFrom = fmt.Sprintf("%v", item[header.ValidFrom])
		}
		if header.ExpiresAt != "" {
			offer.ExpiresAt = fmt.Sprintf("%v", item[header.ExpiresAt])
		}
		if header.Countries != "" {
			offer.Countries = fmt.Sprintf("%v", item[header.Countries])
		}
		res[i] = offer
	}

	return res
}

func processHeader(parsedHeader []string) (*RawHeader, error) {
	resHeader := NewHeader(parsedHeader)
	if err := resHeader.ValidateHeader(); err != nil {
		return nil, err
	}
	return resHeader, nil
}
