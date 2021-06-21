package rawOntology

import (
	"fmt"
	"ts/utils"
)

const (
	RequiredColumnsCount = 4

	defaultCategory           = "UNSPSC"
	defaultCategoryName       = "UNSPSC Name"
	defaultAttributeID        = "Attribute ID"
	defaultAttributeName      = "Attribute Name"
	defaultDefinition         = "Attribute Definition"
	defaultDataType           = "Data Type"
	defaultMaxCharacterLength = "Max Character Length"
	defaultIsRepeatable       = "Is Repeatable"
	defaultMeasurementUoM     = "Measurement UoM"
	defaultIsMandatory        = "Is Mandatory"
	defaultCodedValue         = "Coded Value"
)

type RawHeader struct {
	category           string //required
	categoryName       string
	attributeID        string //skipped
	attributeName      string //required
	definition         string
	dataType           string //required
	maxCharacterLength string
	isRepeatable       string //skipped
	measurementUoM     string
	isMandatory        string //required
	codedValue         string
}

func NewHeader(input []string) *RawHeader {
	var newHeader RawHeader
	for _, columnLabel := range input {
		//required
		trimmedColumnLabel := utils.TrimAll(columnLabel)
		switch trimmedColumnLabel {
		case utils.TrimAll(defaultCategory):
			newHeader.category = columnLabel
		case utils.TrimAll(defaultAttributeName):
			newHeader.attributeName = columnLabel
		case utils.TrimAll(defaultDataType):
			newHeader.dataType = columnLabel
		case utils.TrimAll(defaultIsMandatory):
			newHeader.isMandatory = columnLabel
		//unrequired
		case utils.TrimAll(defaultCategoryName):
			newHeader.categoryName = columnLabel
		case utils.TrimAll(defaultDefinition):
			newHeader.definition = columnLabel
		case utils.TrimAll(defaultMaxCharacterLength):
			newHeader.maxCharacterLength = columnLabel
		case utils.TrimAll(defaultMeasurementUoM):
			newHeader.measurementUoM = columnLabel
		case utils.TrimAll(defaultCodedValue):
			newHeader.codedValue = columnLabel
		}
	}
	return &newHeader
}

func (rh *RawHeader) ValidateHeader() error {
	if !rh.hasRequiredColumns() {
		return fmt.Errorf("ontology does not contains all requiered fields: actual header is "+
			"[Category:%v, AttributeName: %v, DataType: %v, IsMandatory: %v",
			rh.category, rh.attributeName, rh.dataType, rh.isMandatory)
	}
	return nil
}

func (rh *RawHeader) hasRequiredColumns() bool {
	if rh.category == "" || rh.attributeName == "" || rh.dataType == "" || rh.isMandatory == "" {
		return false
	}
	return true
}
