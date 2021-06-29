package rawOntology

import (
	"fmt"
	"log"
	"regexp"
)

const (
	Mandatory = "Mandatory"
	Optional  = "Optional"

	// Data Types
	CodedType  = "Coded"
	FloatType  = "Float"
	NumberType = "Number"
	StringType = "String"
	TextType   = "Text"
)

func ValidateRaw(raw map[string]interface{}, header *RawHeader) []error {
	var errors []error
	if len(raw) < RequiredColumnsCount {
		return []error{fmt.Errorf("not all required colunmns are specified")}
	}

	if header.category != "" {
		err := validateUNSPSC(raw[header.category].(string))
		if err != nil {
			errors = append(errors, err)
		}
	} else {
		log.Fatalf("column Category should be specified")
	}

	if header.categoryName != "" {
		err := validateUNSPSCName(raw[header.categoryName].(string))
		if err != nil {
			errors = append(errors, err)
		}
	}

	if header.attributeName != "" {
		err := validateAttributeName(raw[header.attributeName].(string))
		if err != nil {
			errors = append(errors, err)
		}
	} else {
		log.Fatalf("column Attribute Name should be specified")
	}

	if header.definition != "" {
		err := validateAttributeDefinition(raw[header.definition].(string))
		if err != nil {
			errors = append(errors, err)
		}
	}

	if header.dataType != "" {
		err := validateDataType(raw[header.dataType].(string))
		if err != nil {
			errors = append(errors, err)
		}
	} else {
		log.Fatalf("column Data Type should be specified")

	}

	if header.maxCharacterLength != "" {
		err := validateMaxCharacterLength(raw[header.maxCharacterLength].(string))
		if err != nil {
			errors = append(errors, err)
		}
	}

	if header.measurementUoM != "" {
		err := validateMeasurementUoM(raw[header.measurementUoM].(string))
		if err != nil {
			errors = append(errors, err)
		}
	}

	if header.isMandatory != "" {
		err := validateIsMandatory(raw[header.isMandatory].(string))
		if err != nil {
			errors = append(errors, err)
		}
	} else {
		log.Fatalf("column Is Mandatory should be specified")
	}

	if header.codedValue != "" {
		err := validateCodedValue(raw[header.codedValue].(string))
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func validateUNSPSC(input string) error {
	if input == "" {
		return fmt.Errorf("'UNSPSC' value is required")
	}
	if !isNumber(input) {
		return fmt.Errorf("'UNSPSC' contains not only digits: %v", input)
	}
	return nil
}

func validateUNSPSCName(input string) error {
	return nil
}

func validateAttributeID(input string) error {
	if input != "" && !isNumber(input) {
		return fmt.Errorf("'RawAttribute ID' contains not only digits: %v", input)
	}
	return nil
}

func validateAttributeName(input string) error {
	if input == "" {
		return fmt.Errorf("'RawAttribute UNSPSCName' value is required")
	}
	return nil
}

func validateAttributeDefinition(input string) error {
	return nil
}

func validateDataType(input string) error {
	if input == "" {
		return fmt.Errorf("'Data Type' value is required")
	}
	if input != string(CodedType) &&
		input != string(FloatType) &&
		input != string(NumberType) &&
		input != string(StringType) &&
		input != string(TextType) {
		return fmt.Errorf("'Data Type' field has invalid value %v", input)
	}
	return nil
}

func validateMaxCharacterLength(input string) error {
	if len(input) > 0 && !isNumber(input) {
		return fmt.Errorf("'MaxCharacterLength' contains not only digits: %v", input)
	}
	return nil
}

func validateIsRepeatable(input string) error {
	return nil
}

func validateMeasurementUoM(input string) error {
	//todo when mapping will be supported
	return nil
}

func validateIsMandatory(input string) error {
	if input == "" {
		return fmt.Errorf("'Is Mandatory' field is required")
	}
	if input != string(Mandatory) && input != string(Optional) {
		return fmt.Errorf("'Is Mandatory' field has invalid value: %v, expected: %v", input, string(Mandatory))
	}
	return nil
}

func validateCodedValue(input string) error {
	return nil
}

func isNumber(input string) bool {
	re := regexp.MustCompile(`[0-9]+`)
	return re.MatchString(input)
}
