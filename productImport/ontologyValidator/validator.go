package ontologyValidator

import (
	"fmt"
	"log"
	"strings"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/ontologyRead/rawOntology"
	"ts/productImport/reports"
	"ts/utils"
)

const (
	categoryKey  = "Category" // TS min required column
	productIdKey = "ID"       // TS min required
	nameKey      = "Name"
)

type Validator struct {
}

func NewValidator() ValidatorInterface {
	return &Validator{}
}

func (v *Validator) Validate(data struct {
	Mapping map[string]string
	Rules   *models.OntologyConfig
	Data    []map[string]interface{}
	Report  []*reports.Report
}) ([]reports.Report, bool) {

	if data.Report != nil && len(data.Report) > 0 {
		report, isErr := v.validateReport(data)
		return report, isErr
	}

	report, isErr := v.validateSource(data)
	return report, isErr
}

func (v *Validator) validateSource(data struct {
	Mapping map[string]string
	Rules   *models.OntologyConfig
	Data    []map[string]interface{}
	Report  []*reports.Report
}) ([]reports.Report, bool) {
	feed := make([]reports.Report, 0)
	var columnMapIndex map[string]string
	if data.Mapping != nil && len(data.Mapping) > 0 {
		columnMapIndex = utils.RevertMapKeyValue(data.Mapping)
	}
	catKey := categoryKey
	idKey := productIdKey
	productNameKey := nameKey
	if data.Mapping != nil {
		if val, ok := data.Mapping[catKey]; ok {
			catKey = val
		}
		if val, ok := data.Mapping[idKey]; ok {
			idKey = val
		}
		if val, ok := data.Mapping[productNameKey]; ok {
			productNameKey = val
		}
	}
	isError := false
	for _, product := range data.Data {
		var id string
		var category string
		if val, ok := product[catKey]; ok {
			category = fmt.Sprintf("%v", val)
		} else {
			log.Fatalf("category is not specified")
		}

		if val, ok := product[idKey]; ok {
			id = fmt.Sprintf("%v", val)
		} else {
			log.Fatalf("id is not specified")
		}
		name := ""
		if prodName, ok := product[productNameKey]; ok {
			name = fmt.Sprintf("%v", prodName)
		}

		prodToMapped := make(map[string]string, len(product))
		for k, v := range product {
			i := utils.GetMapOrDefault(k, columnMapIndex)
			prodToMapped[i] = fmt.Sprintf("%v", v)
		}
		if category == "" {
			feed = append(feed, reports.Report{
				ProductId: id,
				Name:      name,
				Category:  category,
				Errors:    []string{"The product category is not specified. The product can not be validated."},
			})
		} else {
			if ruleCategory, ok := data.Rules.Categories[category]; ok {
				for _, attr := range ruleCategory.Attributes {
					val := ""
					message := make([]string, 0)
					if attrVal, ok := prodToMapped[attr.Name]; ok && attrVal != "" {
						val = fmt.Sprintf("%v", attrVal)

						//attrVal check type
						if attr.DataType == rawOntology.FloatType || attr.DataType == rawOntology.NumberType {
							_, err := utils.GetFloat(attrVal)
							if err != nil {
								message = append(message, "The attribute value should be a "+
									strings.ToLower(fmt.Sprintf("%v.", attr.DataType)))
								isError = true
							}
						} else if attr.DataType == rawOntology.CodedType {
							values := strings.Split(attr.CodedValue, ",")
							if exists, _ := utils.InArray(val, values); !exists {
								message = append(
									message,
									"The provided attribute value doesn't match with any "+
										"from the list of predefined values. Look at 'Coded Value' column.")
								isError = true
							}
						}

						// max length
						if attr.MaxCharacterLength > 0 && len(val) > attr.MaxCharacterLength {
							message = append(
								message,
								"The attribute has a too long value (length: %v, max length: %v ).",
								fmt.Sprintf("%v", len(val)),
								fmt.Sprintf("%v", attr.MaxCharacterLength))
							isError = true
						}

						if len(message) == 0 {
							message = append(message, "It is ok!")
						}

					} else {
						text := ""
						if attr.IsMandatory {
							text = "The attribute is mandatory. A value should be provided."
							isError = true
						} else {
							text = "This attribute is optional."
						}
						message = append(message, text)
					}

					d := reports.Report{
						ProductId:    id,
						Name:         name,
						Category:     category,
						CategoryName: ruleCategory.Name,
						AttrName:     attr.Name,
						AttrValue:    val,
						UoM:          attr.MeasurementUoM,
						Errors:       message,
						DataType:     fmt.Sprintf("%v", attr.DataType),
						Description:  attr.Definition,
						IsMandatory:  fmt.Sprintf("%v", attr.IsMandatory),
						CodedVal:     attr.CodedValue,
					}
					feed = append(feed, d)
				}
			} else {
				feed = append(feed, reports.Report{
					ProductId: id,
					Name:      name,
					Category:  category,
					Errors:    []string{"The product category did not match any UNSPSC category from the ontology. The product can not be validated."},
				})
			}
		}
	}
	return feed, isError
}

func (v *Validator) validateReport(data struct {
	Mapping map[string]string
	Rules   *models.OntologyConfig
	Data    []map[string]interface{}
	Report  []*reports.Report
}) ([]reports.Report, bool) {

	feed := make([]reports.Report, 0)
	isError := false
	for _, product := range data.Report {
		category := product.Category
		if category == "" {
			feed = append(feed, reports.Report{
				ProductId: product.ProductId,
				Name:      fmt.Sprintf("%v", product.Name),
				Category:  category,
				Errors:    []string{"The product category is not specified. The product can not be validated."},
			})
		} else {
			if ruleCategory, ok := data.Rules.Categories[category]; ok {
				for _, attr := range ruleCategory.Attributes {

					message := make([]string, 0)
					if product.AttrName == "" || (product.AttrName != "" && product.AttrName == attr.Name) {
						if product.AttrValue != "" {
							//attrVal check type
							if attr.DataType == rawOntology.FloatType || attr.DataType == rawOntology.NumberType {
								_, err := utils.GetFloat(product.AttrValue)
								if err != nil {
									message = append(message, "The attribute value should be a "+
										strings.ToLower(fmt.Sprintf("%v.", attr.DataType)))
									isError = true
								}
							} else if attr.DataType == rawOntology.CodedType {
								values := strings.Split(attr.CodedValue, ",")
								if exists, _ := utils.InArray(product.AttrValue, values); !exists {
									message = append(
										message,
										"The provided attribute value doesn't match with any "+
											"from the list of predefined values. Look at 'Coded Value' column.")
									isError = true
								}
							}

							if attr.MaxCharacterLength > 0 && len(product.AttrValue) > attr.MaxCharacterLength {
								message = append(
									message,
									"The attribute has a too long value (length: %v, max length: %v ).",
									fmt.Sprintf("%v", len(product.AttrValue)),
									fmt.Sprintf("%v", attr.MaxCharacterLength))
								isError = true
							}

							if len(message) == 0 {
								message = append(message, "It is ok!")
							}
						} else if attr.IsMandatory {
							message = append(message, "The attribute is mandatory. A value should be provided.")
							isError = true

						} else {
							message = append(message, "This attribute is optional.")
						}

						d := reports.Report{
							ProductId:    product.ProductId,
							Name:         product.Name,
							Category:     ruleCategory.UNSPSC,
							CategoryName: ruleCategory.Name,
							AttrName:     attr.Name,
							AttrValue:    product.AttrValue,
							UoM:          attr.MeasurementUoM,
							Errors:       message,
							DataType:     fmt.Sprintf("%v", attr.DataType),
							Description:  product.Description,
							IsMandatory:  product.IsMandatory,
							CodedVal:     product.CodedVal,
						}
						feed = append(feed, d)
						if product.AttrName != "" {
							break
						}
					}
				}
			} else {
				feed = append(feed, reports.Report{
					ProductId: product.ProductId,
					Name:      fmt.Sprintf("%v", product.Name),
					Category:  category,
					Errors:    []string{"The product category did not match any UNSPSC category from the ontology. The product can not be validated."},
				})
			}
		}
	}
	return feed, isError
}
