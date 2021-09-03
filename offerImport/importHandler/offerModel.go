package importHandler

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Offer struct {
	ID        string
	Name      string
	Receiver  string
	ValidFrom time.Time
	ExpiresAt time.Time
	Countries []string
}

func newOffer(item map[string]interface{}) *Offer {
	offer := Offer{
		ID:       fmt.Sprintf("%v", item["offerKey"]),
		Name:     fmt.Sprintf("%v", item["name"]),
		Receiver: fmt.Sprintf("%v", item["buyerId"]),
	}

	if item["startDate"] != nil {
		offer.ValidFrom = convertMillisecondsToTime(item["startDate"])
	}

	if item["endDate"] != nil {
		offer.ExpiresAt = convertMillisecondsToTime(item["endDate"])
	}

	if item["countries"] != nil {
		offer.Countries = getCountriesAsArray(fmt.Sprintf("%s", item["countries"]))
	}

	return &offer
}

func convertMillisecondsToTime(input interface{}) time.Time {
	t := int64(input.(float64))
	return time.Unix(0, t*int64(time.Millisecond))
}

func getCountriesAsArray(input string) []string {
	input = strings.Replace(input, "[", "", -1)
	input = strings.Replace(input, "]", "", -1)
	re := regexp.MustCompile("\\b\\w{2}\\b")
	res := re.FindAllString(input, -1)
	return res
}
