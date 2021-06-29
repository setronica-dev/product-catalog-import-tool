package models

//ontology the index
type AttributeConfig struct {
	ID                 string
	Name               string
	Definition         string
	DataType           string
	MaxCharacterLength int
	IsRepeatable       bool
	MeasurementUoM     string
	IsMandatory        bool
	CodedValue         string
}

type CategoryConfig struct {
	UNSPSC     string
	Name       string
	Attributes map[string]*AttributeConfig
}

type OntologyConfig struct {
	Categories map[string]*CategoryConfig
}
