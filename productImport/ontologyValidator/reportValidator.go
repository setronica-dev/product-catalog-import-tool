package ontologyValidator

import (
	"fmt"
	"strings"
	"ts/productImport/attribute"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/ontologyRead/rawOntology"
	"ts/productImport/reports"
	"ts/utils"
)

func (v *Validator) validateReport(data struct {
	Mapping       map[string]string
	Rules         *models.OntologyConfig
	SourceData    []map[string]interface{}
	AttributeData []*attribute.Attribute
}) ([]reports.Report, bool) {
	feed := make([]reports.Report, 0)
	isError := false
	for _, product := range data.AttributeData {
		category := product.Category
		if category == "" {
			feed = append(feed, reports.Report{
				ProductId:   product.ProductId,
				Name:        fmt.Sprintf("%v", product.Name),
				Category:    category,
				AttrName:    product.AttrName,
				AttrValue:   product.AttrValue,
				UoM:         product.UoM,
				DataType:    fmt.Sprintf("%v", product.DataType),
				Description: product.Description,
				IsMandatory: product.IsMandatory,
				CodedVal:    product.CodedVal,
				Errors:      []string{"The product category is not specified. The product can not be validated."},
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
					ProductId:   product.ProductId,
					Name:        fmt.Sprintf("%v", product.Name),
					Category:    category,
					AttrName:    product.AttrName,
					AttrValue:   product.AttrValue,
					UoM:         product.UoM,
					DataType:    fmt.Sprintf("%v", product.DataType),
					Description: product.Description,
					IsMandatory: product.IsMandatory,
					CodedVal:    product.CodedVal,
					Errors:      []string{"The product category did not match any UNSPSC category from the ontology. The product can not be validated."},
				})
			}
		}
	}
	return feed, isError
}
