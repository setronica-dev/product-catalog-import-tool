package offerReader

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"strings"
	"time"
	"ts/adapters"
	"ts/utils"
)

const DateLayout = "2006-01-02"

type OfferReader struct {
	reader adapters.HandlerInterface
}

type RawOffer struct {
	Offer        string
	ReceiverName string
	Contract     string
	ValidFrom    time.Time
	ExpiresAt    time.Time
	Countries    []string
}

type Deps struct {
	dig.In
	Reader      adapters.HandlerInterface
	FileManager *adapters.FileManager
}

func NewOfferReader(deps Deps) *OfferReader {
	return &OfferReader{
		reader: deps.Reader,
	}
}

func (o *OfferReader) UploadOffers(path string) []RawOffer {
	ext := adapters.GetFileType(path)
	o.reader.Init(ext)
	parsedRaws := o.reader.Parse(path)
	actualHeader := o.reader.GetHeader()
	header, err := processHeader(actualHeader)
	if err != nil {
		log.Printf("failed to upload offers: %v", err)
		return nil
	}
	or := processOffers(parsedRaws, header)
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
	if utils.IsEmptyMap(row) {
		return nil
	}
	if row[header.Offer] == nil || row[header.Receiver] == nil {
		log.Printf("row does not contain values in required columns (Offer, Receiver). Actual value: %v", row)
		return nil
	}

	offer := RawOffer{
		Offer:        fmt.Sprintf("%v", row[header.Offer]),
		ReceiverName: fmt.Sprintf("%v", row[header.Receiver]),
	}
	if header.ContractID != "" && row[header.ContractID] != "" {
		offer.Contract = fmt.Sprintf("%v", row[header.ContractID])
	}
	if header.ValidFrom != "" && row[header.ValidFrom] != "" {
		date, err := time.Parse(DateLayout, fmt.Sprintf("%v", row[header.ValidFrom]))
		if err == nil {
			offer.ValidFrom = date
		} else {
			log.Printf("invalid format of \"valid_from\" field: should be YYYY-MM-DD: %v", err)
		}
	}
	if header.ExpiresAt != "" && row[header.ExpiresAt] != "" {
		date, err := time.Parse(DateLayout, fmt.Sprintf("%v", row[header.ExpiresAt]))
		if err == nil {
			offer.ExpiresAt = date
		} else {
			log.Printf("invalid format of \"expies_at\" field: should be YYYY-MM-DD: %v", err)
		}
	}
	if header.Countries != "" && row[header.Countries] != "" {
		offer.Countries = getCountries(row[header.Countries])
	}
	return &offer
}

func getCountries(input interface{}) []string {
	return strings.SplitN(strings.ToLower(fmt.Sprintf("%v", input)), ",", -1)
}

func processHeader(parsedHeader []string) (*RawHeader, error) {
	resHeader := NewHeader(parsedHeader)
	if err := resHeader.ValidateHeader(); err != nil {
		return nil, err
	}
	return resHeader, nil
}
