package ontologyValidator

import (
	"fmt"
	"strings"
	"ts/productImport/attribute"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/ontologyRead/rawOntology"
	"ts/productImport/product"
	"ts/productImport/reports"
	"ts/utils"
)

func (v *Validator) validateReport(
	rulesData *models.OntologyConfig,
	productData *product.Products,
	attributeData []*attribute.Attribute) ([]reports.Report, bool) {
	feed := make([]reports.Report, 0)
	isError := false

	for _, reportAttribute := range attributeData {
		category := getProductCategory(reportAttribute.ProductId, productData, reportAttribute)
		if category == "" {
			feed = append(feed, reports.Report{
				ProductId:   reportAttribute.ProductId,
				Name:        fmt.Sprintf("%v", reportAttribute.Name),
				Category:    category,
				AttrName:    reportAttribute.AttrName,
				AttrValue:   reportAttribute.AttrValue,
				UoM:         reportAttribute.UoM,
				DataType:    fmt.Sprintf("%v", reportAttribute.DataType),
				Description: reportAttribute.Description,
				IsMandatory: reportAttribute.IsMandatory,
				CodedVal:    reportAttribute.CodedVal,
				Errors:      []string{"The attribute category is not specified. The attribute can not be validated."},
			})
		} else {
			if ruleCategory, ok := rulesData.Categories[category]; ok {
				for _, attr := range ruleCategory.Attributes {
					message := make([]string, 0)
					if reportAttribute.AttrName == "" || (reportAttribute.AttrName != "" && reportAttribute.AttrName == attr.Name) {
						if reportAttribute.AttrValue != "" {
							//attrVal check type
							if attr.DataType == rawOntology.FloatType || attr.DataType == rawOntology.NumberType {
								_, err := utils.GetFloat(reportAttribute.AttrValue)
								if err != nil {
									message = append(message, "The attribute value should be a "+
										strings.ToLower(fmt.Sprintf("%v.", attr.DataType)))
									isError = true
								}
							} else if attr.DataType == rawOntology.CodedType {
								values := strings.Split(attr.CodedValue, ",")
								if exists, _ := utils.InArray(reportAttribute.AttrValue, values); !exists {
									message = append(
										message,
										"The provided attribute value doesn't match with any "+
											"from the list of predefined values. Look at 'Coded Value' column.")
									isError = true
								}
							}

							if attr.MaxCharacterLength > 0 && len(reportAttribute.AttrValue) > attr.MaxCharacterLength {
								message = append(
									message,
									"The attribute has a too long value (length: %v, max length: %v ).",
									fmt.Sprintf("%v", len(reportAttribute.AttrValue)),
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
							ProductId:    reportAttribute.ProductId,
							Name:         reportAttribute.Name,
							Category:     ruleCategory.UNSPSC,
							CategoryName: ruleCategory.Name,
							AttrName:     attr.Name,
							AttrValue:    reportAttribute.AttrValue,
							UoM:          attr.MeasurementUoM,
							Errors:       message,
							DataType:     fmt.Sprintf("%v", attr.DataType),
							Description:  reportAttribute.Description,
							IsMandatory:  reportAttribute.IsMandatory,
							CodedVal:     reportAttribute.CodedVal,
						}
						feed = append(feed, d)
						if reportAttribute.AttrName != "" {
							break
						}
					}
				}
			} else {
				feed = append(feed, reports.Report{
					ProductId:   reportAttribute.ProductId,
					Name:        fmt.Sprintf("%v", reportAttribute.Name),
					Category:    category,
					AttrName:    reportAttribute.AttrName,
					AttrValue:   reportAttribute.AttrValue,
					UoM:         reportAttribute.UoM,
					DataType:    fmt.Sprintf("%v", reportAttribute.DataType),
					Description: reportAttribute.Description,
					IsMandatory: reportAttribute.IsMandatory,
					CodedVal:    reportAttribute.CodedVal,
					Errors:      []string{"The attribute category did not match any UNSPSC category from the ontology. The attribute can not be validated."},
				})
			}
		}
	}
	return feed, isError
}

func getProductCategory(productID string, products *product.Products, reportAttribute *attribute.Attribute) string {
	targetProduct := products.FindProductByID(productID)
	if targetProduct != nil {
		return targetProduct.Category
	}
	return reportAttribute.Category
}
