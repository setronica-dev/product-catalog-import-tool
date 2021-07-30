package ontologyRead

import (
	"log"
	"os"
	"ts/adapters"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/ontologyRead/rawOntology"
)

type RulesHandler struct {
	sourcePath string
	reader     adapters.HandlerInterface
}

func NewRulesHandler(deps Deps) *RulesHandler {
	return &RulesHandler{
		sourcePath: deps.Config.ProductCatalog.OntologyPath,
		reader:     deps.Handler,
	}
}

func (h *RulesHandler) InitRulesConfig() (*models.OntologyConfig, error) {
	var rulesConfig *models.OntologyConfig
	var rules *rawOntology.RawOntology
	sourcePath := h.sourcePath
	if sourcePath != "" {
		if _, err := os.Stat(sourcePath); !os.IsNotExist(err) {
			rules = h.UploadRules(sourcePath)
			rulesConfig = rules.ToConfig()
		} else {
			log.Fatalf("ontology file does not exists. Please fill and add it to %v", sourcePath)
		}
	} else {
		log.Fatalf("ontology path is not specified")
	}
	return rulesConfig, nil
}

func (h *RulesHandler) UploadRules(path string) *rawOntology.RawOntology {

	ext := adapters.GetFileType(path)
	h.reader.Init(ext)
	parsedRaws := h.reader.Parse(path)
	actualHeader := h.reader.GetHeader()
	header, err := processHeader(actualHeader)
	if err != nil {
		log.Fatalf("failed to upload rules: %v", err)
	}
	o := processOntology(parsedRaws, header)
	log.Printf("Rules upload finished. Proceeded %v lines, uploaded %v categories", len(parsedRaws)+1, o.GetCategoriesCount())
	return o
}

func processHeader(parsedHeader []string) (*rawOntology.RawHeader, error) {
	resHeader := rawOntology.NewHeader(parsedHeader)
	if err := resHeader.ValidateHeader(); err != nil {
		return nil, err
	}
	return resHeader, nil
}

func processOntology(parsedRaws []map[string]interface{}, header *rawOntology.RawHeader) *rawOntology.RawOntology {
	o := rawOntology.NewRawOntology()

	for i, raw := range parsedRaws {
		var errors []error

		errors = rawOntology.ValidateRaw(raw, header)
		if len(errors) == 0 {
			rawAttribute := rawOntology.NewRawAttribute(raw, header)
			rawCategory := rawOntology.NewRawCategory(raw, header)
			err := o.AddCategoryAttribute(rawCategory, rawAttribute)
			if err != nil {
				//log.Printf("raw %v error: %v", i, err)
			}
		} else {
			log.Printf("raw %v: validation errors: %v", i, errors)
		}
	}
	return o
}
