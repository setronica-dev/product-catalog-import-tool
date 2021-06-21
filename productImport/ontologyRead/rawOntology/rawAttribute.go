package rawOntology

import (
	"strconv"
	"ts/productImport/ontologyRead/models"
	"ts/utils"
)

type RawAttribute struct {
	id                 string
	name               string
	definition         string
	dataType           string
	maxCharacterLength int
	isRepeatable       bool
	measurementUoM     string
	isMandatory        bool
	codedValue         string
}

func NewRawAttribute(raw map[string]interface{}, header *RawHeader) *RawAttribute {
	attribute := RawAttribute{
		name:        raw[header.attributeName].(string),
		dataType:    raw[header.dataType].(string),
		isMandatory: convertIsMandatory(raw[header.isMandatory].(string)),
	}
	if header.definition != "" {
		attribute.definition = raw[header.definition].(string)
	}
	if header.maxCharacterLength != "" {
		attribute.maxCharacterLength = convertMaxCharacterLength(raw[header.maxCharacterLength].(string))
	}
	if header.measurementUoM != "" {
		attribute.measurementUoM = raw[header.measurementUoM].(string)
	}
	if header.codedValue != "" {
		attribute.codedValue = raw[header.codedValue].(string)
	}
	return &attribute
}

func convertMaxCharacterLength(input string) int {
	res, _ := strconv.Atoi(input)
	return res
}

func convertIsRepeatable(input string) bool {
	if input == "No" {
		return false
	}
	return true
}

func convertIsMandatory(input string) bool {
	if utils.TrimAll(input) == utils.TrimAll(Mandatory) {
		return true
	}
	return false
}

func (a *RawAttribute) ToConfig() *models.AttributeConfig {
	return &models.AttributeConfig{
		ID:                 a.id,
		Name:               a.name,
		Definition:         a.definition,
		DataType:           a.dataType,
		MaxCharacterLength: a.maxCharacterLength,
		IsRepeatable:       a.isRepeatable,
		MeasurementUoM:     a.measurementUoM,
		IsMandatory:        a.isMandatory,
		CodedValue:         a.codedValue,
	}
}
