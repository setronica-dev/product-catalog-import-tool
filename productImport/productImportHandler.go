package productImport

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"os"
	"regexp"
	"ts/adapters"
	"ts/config"
	"ts/productImport/attribute"
	"ts/productImport/mapping"
	"ts/productImport/ontologyRead"
	"ts/productImport/ontologyRead/models"
	"ts/productImport/ontologyValidator"
	"ts/productImport/product"
	"ts/productImport/reports"
	"ts/productImport/tradeshiftImportHandler"
	"ts/utils"
)

const stageName = "Product Validation Import stage"

type ProductImportHandler struct {
	config           *config.Config
	mapHandler       mapping.MappingHandlerInterface
	rulesHandler     *ontologyRead.RulesHandler
	productHandler   product.ProductHandlerInterface
	attributeHandler attribute.AttributeHandlerInterface
	handler          adapters.HandlerInterface
	validator        ontologyValidator.ValidatorInterface
	reports          *reports.ReportsHandler
	fileManager      *adapters.FileManager
	importHandler    *tradeshiftImportHandler.TradeshiftHandler
	columnMap        *ColumnMap
}

type ColumnMap struct {
	Category  string
	ProductID string
	Name      string
}

type Deps struct {
	dig.In
	Config           *config.Config
	MapHandler       mapping.MappingHandlerInterface
	RulesHandler     *ontologyRead.RulesHandler
	ProductHandler   product.ProductHandlerInterface
	AttributeHandler attribute.AttributeHandlerInterface
	Handler          adapters.HandlerInterface
	Validator        ontologyValidator.ValidatorInterface
	Reports          *reports.ReportsHandler
	FileManager      *adapters.FileManager
	ImportHandler    *tradeshiftImportHandler.TradeshiftHandler
}

func NewProductImportHandler(deps Deps) *ProductImportHandler {
	m := deps.MapHandler.Parse()
	return &ProductImportHandler{
		config:           deps.Config,
		mapHandler:       deps.MapHandler,
		rulesHandler:     deps.RulesHandler,
		attributeHandler: deps.AttributeHandler,
		productHandler:   deps.ProductHandler,
		handler:          deps.Handler,
		validator:        deps.Validator,
		reports:          deps.Reports,
		fileManager:      deps.FileManager,
		importHandler:    deps.ImportHandler,
		columnMap: &ColumnMap{
			ProductID: m.ProductID,
			Category:  m.Category,
			Name:      m.Name,
		},
	}
}

func (ph *ProductImportHandler) RunCSV() {
	//ontology
	var rulesConfig *models.OntologyConfig

	rulesConfig, err := ph.rulesHandler.InitRulesConfig()
	if err != nil {
		log.Fatalf("ontology was not specified: %v", err)
	}

	// mappings
	columnMap := ph.mapHandler.Get()

	// feed
	err = ph.processProducts(columnMap, rulesConfig)
	if err != nil {
		log.Println(err)
	}
}

func (ph *ProductImportHandler) processProducts(columnMap map[string]string, rulesConfig *models.OntologyConfig) error {
	// if something in progress

	var processedSource []string
	inProgress := adapters.GetFiles(ph.config.ProductCatalog.InProgressPath)
	sources := adapters.GetFiles(ph.config.ProductCatalog.SourcePath)
	// identify fitting report
	if len(inProgress) > 0 {
		for _, processingFile := range inProgress {
			reportFile := findReport(processingFile, utils.SliceDiff(sources, processedSource))
			if reportFile != "" {
				processedSource = append(processedSource, reportFile)
				ph.processFeed(
					ph.config.ProductCatalog.InProgressPath+"/"+processingFile, //feed
					ph.config.ProductCatalog.SourcePath+reportFile,             //report
					columnMap,
					rulesConfig,
					false,
				)
			} else {
				log.Printf("You have the failed feed in progress '%v'. "+
					"Please check the failure report in '%v', "+
					"fill it with the data and appload to the '%v' folder.",
					ph.config.ProductCatalog.InProgressPath+"/"+processingFile,
					ph.config.ProductCatalog.FailResultPath,
					ph.config.ProductCatalog.SourcePath)
			}
		}
	} else if len(sources) == 0 {
		return fmt.Errorf("SOURCE IS NOT FOUND")
	}

	for _, source := range sources {
		if inArr, _ := utils.InArray(source, processedSource); !inArr {
			ph.processFeed(
				ph.config.ProductCatalog.SourcePath+source,
				"",
				columnMap,
				rulesConfig,
				true,
			)
		}
	}
	return nil
}

func (ph *ProductImportHandler) processFeed(
	sourceFeedPath string,
	validationReportPath string,
	columnMap map[string]string,
	ruleConfig *models.OntologyConfig,
	isInitial bool,
) {
	log.Println("_________________________________")
	log.Printf("PROCESSING SOURCE: %v", sourceFeedPath)
	var er error
	if validationReportPath != "" {
		log.Printf("EDITED REPORT: %v", validationReportPath)
		if validationReportPath, er = adapters.MoveToPath(validationReportPath, ph.config.ProductCatalog.InProgressPath); er != nil {
			log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", validationReportPath, ph.config.ProductCatalog.InProgressPath)
		}
	}

	if isInitial {
		if sourceFeedPath, er = adapters.MoveToPath(sourceFeedPath, ph.config.ProductCatalog.InProgressPath); er != nil {
			log.Printf("ERROR COPYING THE '%v' FILE to the  '%v' folder", sourceFeedPath, ph.config.ProductCatalog.InProgressPath)
		}
	}

	// fixed attributes
	attributeReportData, err := ph.attributeHandler.InitAttributeData(validationReportPath)
	if err != nil {
		log.Printf("failed to upload attributes report %v: %v", validationReportPath, err)
	}

	// source
	parsedData, err := ph.productHandler.InitSourceData(sourceFeedPath)
	if err != nil {
		log.Printf("failed to upload source data %v: %v", sourceFeedPath, err)
	}

	// validation feed
	feed, hasErrors := ph.validator.Validate(struct {
		Mapping       map[string]string
		Rules         *models.OntologyConfig
		SourceData    []map[string]interface{}
		AttributeData []*attribute.Attribute
	}{
		Mapping:       columnMap,
		Rules:         ruleConfig,
		SourceData:    parsedData,
		AttributeData: attributeReportData,
	})

	if !hasErrors {
		log.Printf("SUCCESS: FILE IS VALID. Please check the '%s' folder", ph.config.ProductCatalog.SentPath)
		if _, er = adapters.MoveToPath(sourceFeedPath, ph.config.ProductCatalog.SentPath); er != nil {
			log.Printf("ERROR COPYING THE SOURCE FILE %v to the '%v' folder", sourceFeedPath, ph.config.ProductCatalog.SentPath)
		}

		if validationReportPath != "" {
			if _, er = adapters.MoveToPath(validationReportPath, ph.config.ProductCatalog.SentPath); er != nil {
				log.Printf("ERROR COPYING THE REPORT FILE %v to the '%v' folder", validationReportPath, ph.config.ProductCatalog.SentPath)
			}
		}
	} else {
		log.Printf("FAILURE: check the failure report in '%v', fill it with the data and upload to the '%v' folder.",
			ph.config.ProductCatalog.ReportPath,
			ph.config.ProductCatalog.SourcePath)
		if validationReportPath != "" {
			e := os.Remove(validationReportPath)

			if e != nil {
				log.Println(e)
			}
		}
	}

	cleanUpFailures(sourceFeedPath, ph.config.ProductCatalog.FailResultPath)
	validationReportPath = ph.reports.WriteReport(sourceFeedPath, hasErrors, feed, parsedData)
	if !hasErrors {
		log.Println("IMPORT FEED TO TRADESHIFT WAS STARTED")
		er := ph.importHandler.ImportFeedToTradeshift(validationReportPath)
		if er != nil {
			log.Printf("FAILED TO IMPORT VALID FEED TO TRADESHIFT. Reason: %v", er)
		}
	}
}

func findReport(inProgressFile string, sources []string) string {
	report := ""
	pattern := adapters.GetFileName(inProgressFile)

	for _, source := range sources {
		regexp, _ := regexp.Compile(`(-failures)`)
		match := regexp.FindStringIndex(source)
		if len(match) == 2 {
			name := string(source[0:match[0]])
			if name == pattern {
				return source
			}
		}
	}
	return report
}

func cleanUpFailures(sourceFile string, folder string) {
	reportsList := adapters.GetFiles(folder)
	for _, source := range reportsList {
		del := findReport(sourceFile, []string{source})
		if del != "" {
			e := os.Remove(folder + "/" + del)
			if e != nil {
				log.Println(e)
			}
		}
	}
}
